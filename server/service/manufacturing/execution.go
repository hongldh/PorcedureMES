// 阶段 5b：派工 + 报工 service
//
// === 学习要点 ===
// 1. 工序顺序约束：第 N 道累计(合格+不良) ≤ 第 N-1 道累计合格 —— 防止"后道先做"
// 2. 报工即时生效：无草稿态，Create 直接写事实 + 回写工序累计；错了用 Void 反向冲回
// 3. 报工人校验：派给个人必须本人报；派给班组必须组员报；未派工不可报
// 4. 行锁：对工序行 SELECT ... FOR UPDATE，防止并发报工超计划
package manufacturing

import (
	"errors"
	"fmt"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/manufacturing"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const qtyEps = 0.0001

// ============================================================
// 派工
// ============================================================
type DispatchService struct{}

// DispatchParam 派工参数：员工 / 班组 至少填一个
type DispatchParam struct {
	OpID          uint `json:"opId"`
	AssignedUserID uint `json:"assignedUserId"`
	TeamID        uint `json:"teamId"`
	WorkstationID uint `json:"workstationId"` // 可覆盖默认工位
}

// Dispatch 派工 / 改派（已完工的工序不可改派）
func (s *DispatchService) Dispatch(p DispatchParam, opUserID uint) error {
	if p.AssignedUserID == 0 && p.TeamID == 0 {
		return errors.New("必须指定员工或班组")
	}
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var op manufacturing.MfgWorkOrderOp
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&op, p.OpID).Error; err != nil {
			return errors.New("工序不存在")
		}
		if op.Status == manufacturing.OpStatusFinished {
			return errors.New("工序已完工，不可改派")
		}
		var wo manufacturing.MfgWorkOrder
		if err := tx.First(&wo, op.WorkOrderID).Error; err != nil {
			return err
		}
		if wo.Status != manufacturing.WoStatusReleased && wo.Status != manufacturing.WoStatusInProgress {
			return errors.New("工单状态不允许派工（需已下达或生产中）")
		}
		// 班组校验：若同时指定员工与班组，员工必须属于该班组
		if p.TeamID > 0 {
			var team manufacturing.MfgTeam
			if err := tx.First(&team, p.TeamID).Error; err != nil {
				return errors.New("班组不存在")
			}
			if team.Status != manufacturing.BasicStatusEnabled {
				return errors.New("班组已停用")
			}
			if p.AssignedUserID > 0 {
				var cnt int64
				tx.Model(&manufacturing.MfgTeamMember{}).
					Where("team_id=? AND user_id=?", p.TeamID, p.AssignedUserID).Count(&cnt)
				if cnt == 0 {
					return errors.New("指定员工不属于该班组")
				}
			}
		}
		updates := map[string]interface{}{
			"assigned_user_id": p.AssignedUserID,
			"team_id":          p.TeamID,
		}
		if p.WorkstationID > 0 {
			updates["workstation_id"] = p.WorkstationID
		}
		// 未派工 → 已派工；进行中保持不变（改派不回退进度）
		if op.Status == manufacturing.OpStatusPending {
			updates["status"] = manufacturing.OpStatusDispatched
		}
		return tx.Model(&op).Updates(updates).Error
	})
}

// ============================================================
// 报工
// ============================================================
type WorkReportService struct{}

// Create 报工（即时生效）
func (s *WorkReportService) Create(e *manufacturing.MfgWorkReport, opUserID uint) error {
	if e.OpID == 0 {
		return errors.New("必须指定工序")
	}
	if e.GoodQty < 0 || e.RejectQty < 0 || e.Hours < 0 {
		return errors.New("数量/工时不能为负")
	}
	if e.GoodQty+e.RejectQty <= 0 {
		return errors.New("合格数与不良数不能同时为 0")
	}
	if e.UserID == 0 {
		e.UserID = opUserID // 默认当前登录人报工
	}
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// 1. 锁工序行
		var op manufacturing.MfgWorkOrderOp
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&op, e.OpID).Error; err != nil {
			return errors.New("工序不存在")
		}
		if op.Status == manufacturing.OpStatusPending {
			return errors.New("工序未派工，不可报工")
		}
		if op.Status == manufacturing.OpStatusFinished {
			return errors.New("工序已完工，不可再报工")
		}
		// 2. 工单状态
		var wo manufacturing.MfgWorkOrder
		if err := tx.First(&wo, op.WorkOrderID).Error; err != nil {
			return err
		}
		if wo.Status != manufacturing.WoStatusReleased && wo.Status != manufacturing.WoStatusInProgress {
			return errors.New("工单状态不允许报工")
		}
		// 3. 报工人校验
		if err := checkReporter(tx, &op, e.UserID); err != nil {
			return err
		}
		// 4. 数量约束
		reportTotal := e.GoodQty + e.RejectQty
		// 4a. 本道累计不超计划
		if op.DoneQty+op.RejectQty+reportTotal > op.PlanQty+qtyEps {
			return fmt.Errorf("超计划：工序[%s]计划%.4f 已报%.4f 本次%.4f",
				op.Code, op.PlanQty, op.DoneQty+op.RejectQty, reportTotal)
		}
		// 4b. 工序顺序：本道累计(合格+不良) ≤ 上道累计合格
		var prev manufacturing.MfgWorkOrderOp
		err := tx.Where("work_order_id=? AND seq<?", op.WorkOrderID, op.Seq).
			Order("seq DESC").First(&prev).Error
		if err == nil {
			if op.DoneQty+op.RejectQty+reportTotal > prev.DoneQty+qtyEps {
				return fmt.Errorf("上道工序[%s]合格数仅 %.4f，本道累计报工不能超过该值（已报%.4f 本次%.4f）",
					prev.Code, prev.DoneQty, op.DoneQty+op.RejectQty, reportTotal)
			}
		} else if err != gorm.ErrRecordNotFound {
			return err
		}
		// 5. 写报工记录
		no, err := seqService.NextDocNo(tx, "MFG_WORK_REPORT")
		if err != nil {
			return err
		}
		e.DocNo = no
		e.WorkOrderID = op.WorkOrderID
		e.Status = manufacturing.ReportStatusDone
		if e.WorkstationID == 0 {
			e.WorkstationID = op.WorkstationID
		}
		if e.ReportTime == "" {
			e.ReportTime = time.Now().Format("2006-01-02 15:04:05")
		}
		if err := tx.Create(e).Error; err != nil {
			return err
		}
		// 6. 回写工序累计 + 状态
		newDone := op.DoneQty + e.GoodQty
		newReject := op.RejectQty + e.RejectQty
		newHours := op.ActualHours + e.Hours
		newStatus := manufacturing.OpStatusInProgress
		if newDone >= op.PlanQty-qtyEps {
			newStatus = manufacturing.OpStatusFinished
		}
		if err := tx.Model(&op).Updates(map[string]interface{}{
			"done_qty": newDone, "reject_qty": newReject,
			"actual_hours": newHours, "status": newStatus,
		}).Error; err != nil {
			return err
		}
		// 7. 工单：已下达 → 生产中
		if wo.Status == manufacturing.WoStatusReleased {
			return tx.Model(&wo).Update("status", manufacturing.WoStatusInProgress).Error
		}
		return nil
	})
}

// Void 作废报工：反向冲回工序累计；校验下道工序未依赖本道合格数
func (s *WorkReportService) Void(id, opUserID uint) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var r manufacturing.MfgWorkReport
		if err := tx.First(&r, id).Error; err != nil {
			return err
		}
		if r.Status != manufacturing.ReportStatusDone {
			return errors.New("仅已生效的报工可作废")
		}
		var op manufacturing.MfgWorkOrderOp
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&op, r.OpID).Error; err != nil {
			return err
		}
		newDone := op.DoneQty - r.GoodQty
		newReject := op.RejectQty - r.RejectQty
		newHours := op.ActualHours - r.Hours
		if newDone < -qtyEps || newReject < -qtyEps {
			return errors.New("冲回后累计为负，数据异常")
		}
		if newDone < 0 {
			newDone = 0
		}
		if newReject < 0 {
			newReject = 0
		}
		if newHours < 0 {
			newHours = 0
		}
		// 下道工序依赖校验：下道累计(合格+不良) 不能超过本道冲回后的合格
		var next manufacturing.MfgWorkOrderOp
		err := tx.Where("work_order_id=? AND seq>?", op.WorkOrderID, op.Seq).
			Order("seq ASC").First(&next).Error
		if err == nil {
			if next.DoneQty+next.RejectQty > newDone+qtyEps {
				return fmt.Errorf("下道工序[%s]已报工 %.4f，作废后本道合格仅 %.4f，请先作废下道报工",
					next.Code, next.DoneQty+next.RejectQty, newDone)
			}
		} else if err != gorm.ErrRecordNotFound {
			return err
		}
		// 完工入库依赖校验：末道工序被冲回时，已入库数不能超过冲回后合格数
		var last manufacturing.MfgWorkOrderOp
		if err := tx.Where("work_order_id=?", op.WorkOrderID).Order("seq DESC").First(&last).Error; err == nil && last.ID == op.ID {
			var wo manufacturing.MfgWorkOrder
			if err := tx.First(&wo, op.WorkOrderID).Error; err != nil {
				return err
			}
			if wo.ActualQty > newDone+qtyEps {
				return fmt.Errorf("已完工入库 %.4f，作废后末道合格仅 %.4f，请先作废完工入库单", wo.ActualQty, newDone)
			}
		}
		// 状态回退
		newStatus := op.Status
		if newDone+newReject <= qtyEps {
			newStatus = manufacturing.OpStatusDispatched
		} else if newStatus == manufacturing.OpStatusFinished {
			newStatus = manufacturing.OpStatusInProgress
		}
		if err := tx.Model(&op).Updates(map[string]interface{}{
			"done_qty": newDone, "reject_qty": newReject,
			"actual_hours": newHours, "status": newStatus,
		}).Error; err != nil {
			return err
		}
		return tx.Model(&r).Update("status", manufacturing.ReportStatusVoid).Error
	})
}

// checkLastOpDone 完工入库联动：累计入库数不得超过末道工序累计合格数
//
// 工单若无工序（5b 之前下达的旧工单）则跳过，保持向后兼容
func checkLastOpDone(tx *gorm.DB, workOrderID uint, totalAfter float64) error {
	var last manufacturing.MfgWorkOrderOp
	err := tx.Where("work_order_id=?", workOrderID).Order("seq DESC").First(&last).Error
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	if err != nil {
		return err
	}
	if totalAfter > last.DoneQty+qtyEps {
		return fmt.Errorf("末道工序[%s]累计合格仅 %.4f，完工入库累计不能超过该值（拟入库后累计 %.4f），请先报工",
			last.Code, last.DoneQty, totalAfter)
	}
	return nil
}

// checkReporter 报工人校验
func checkReporter(tx *gorm.DB, op *manufacturing.MfgWorkOrderOp, userID uint) error {
	if op.AssignedUserID > 0 {
		if op.AssignedUserID != userID {
			return errors.New("该工序已指派给其他员工，您无权报工")
		}
		return nil
	}
	if op.TeamID > 0 {
		var cnt int64
		tx.Model(&manufacturing.MfgTeamMember{}).
			Where("team_id=? AND user_id=?", op.TeamID, userID).Count(&cnt)
		if cnt == 0 {
			return errors.New("您不是该工序指派班组的成员，无权报工")
		}
	}
	return nil
}

// ============================================================
// 查询
// ============================================================
type ReportSearch struct {
	request.PageInfo
	WorkOrderID uint   `json:"workOrderId" form:"workOrderId"`
	OpID        uint   `json:"opId" form:"opId"`
	UserID      uint   `json:"userId" form:"userId"`
	Status      int8   `json:"status" form:"status"`
	DocNo       string `json:"docNo" form:"docNo"`
}

func (s *WorkReportService) GetList(info ReportSearch) (list []manufacturing.MfgWorkReport, total int64, err error) {
	if info.Page <= 0 {
		info.Page = 1
	}
	if info.PageSize <= 0 {
		info.PageSize = 10
	}
	db := global.GVA_DB.Model(&manufacturing.MfgWorkReport{})
	if info.WorkOrderID > 0 {
		db = db.Where("work_order_id=?", info.WorkOrderID)
	}
	if info.OpID > 0 {
		db = db.Where("op_id=?", info.OpID)
	}
	if info.UserID > 0 {
		db = db.Where("user_id=?", info.UserID)
	}
	if info.Status > 0 {
		db = db.Where("status=?", info.Status)
	}
	if info.DocNo != "" {
		db = db.Where("doc_no LIKE ?", "%"+info.DocNo+"%")
	}
	if err = db.Count(&total).Error; err != nil {
		return
	}
	err = db.Preload("WorkOrder.Product").Preload("Op").Preload("User").Preload("Workstation").
		Limit(info.PageSize).Offset((info.Page - 1) * info.PageSize).
		Order("id DESC").Find(&list).Error
	return
}

// MyOps 当前用户可报工的工序（派给本人 或 派给本人所在班组），且工单在产
func (s *WorkReportService) MyOps(userID uint) (list []manufacturing.MfgWorkOrderOp, err error) {
	var teamIDs []uint
	global.GVA_DB.Model(&manufacturing.MfgTeamMember{}).Where("user_id=?", userID).Pluck("team_id", &teamIDs)

	db := global.GVA_DB.Model(&manufacturing.MfgWorkOrderOp{}).
		Joins("JOIN mfg_work_order wo ON wo.id = mfg_work_order_op.work_order_id AND wo.deleted_at IS NULL").
		Where("wo.status IN ?", []int8{manufacturing.WoStatusReleased, manufacturing.WoStatusInProgress}).
		Where("mfg_work_order_op.status IN ?", []int8{manufacturing.OpStatusDispatched, manufacturing.OpStatusInProgress})
	if len(teamIDs) > 0 {
		db = db.Where("mfg_work_order_op.assigned_user_id = ? OR (mfg_work_order_op.assigned_user_id = 0 AND mfg_work_order_op.team_id IN ?)", userID, teamIDs)
	} else {
		db = db.Where("mfg_work_order_op.assigned_user_id = ?", userID)
	}
	err = db.Preload("Workstation").Order("mfg_work_order_op.work_order_id, mfg_work_order_op.seq").Find(&list).Error
	return
}
