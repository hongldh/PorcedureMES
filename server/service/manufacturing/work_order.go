package manufacturing

import (
	"errors"
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/manufacturing"
	"gorm.io/gorm"
)

// ============================================================
// 生产工单服务
// ============================================================
type WorkOrderService struct{}

// Create 创建工单（草稿）
func (s *WorkOrderService) Create(e *manufacturing.MfgWorkOrder) error {
	if e.ProductID == 0 || e.PlanQty <= 0 {
		return errors.New("成品与计划数量必填")
	}
	e.Status = manufacturing.WoStatusDraft
	e.ActualQty = 0
	e.IssuedAmount = 0
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		no, err := seqService.NextDocNo(tx, "MFG_WORK_ORDER")
		if err != nil {
			return err
		}
		e.DocNo = no
		// 草稿态不展开物料需求
		e.Items = nil
		return tx.Create(e).Error
	})
}

// Update 草稿态可修改成品/数量/BOM；其他状态只能改备注与计划日期
func (s *WorkOrderService) Update(e *manufacturing.MfgWorkOrder) error {
	var old manufacturing.MfgWorkOrder
	if err := global.GVA_DB.First(&old, e.ID).Error; err != nil {
		return err
	}
	if old.Status != manufacturing.WoStatusDraft {
		// 已下达：仅允许改备注、计划日期、入库仓
		return global.GVA_DB.Model(&old).Updates(map[string]interface{}{
			"remark":       e.Remark,
			"plan_start":   e.PlanStart,
			"plan_end":     e.PlanEnd,
			"warehouse_id": e.WarehouseID,
		}).Error
	}
	// 草稿态：全量替换头字段，明细在下达时才展开，这里不动 Items
	e.DocNo = old.DocNo
	e.Status = old.Status
	e.CreatedAt = old.CreatedAt
	e.ActualQty = old.ActualQty
	e.IssuedAmount = old.IssuedAmount
	e.Items = nil
	return global.GVA_DB.Save(e).Error
}

// Delete 删除工单（仅草稿）
func (s *WorkOrderService) Delete(id uint) error {
	var d manufacturing.MfgWorkOrder
	if err := global.GVA_DB.First(&d, id).Error; err != nil {
		return err
	}
	if err := mustWoDraft(d.Status); err != nil {
		return err
	}
	return global.GVA_DB.Delete(&manufacturing.MfgWorkOrder{}, id).Error
}

// Release 下达工单：按 BOM × PlanQty × (1+lossRate) 展开物料需求
func (s *WorkOrderService) Release(id, opUserID uint) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var d manufacturing.MfgWorkOrder
		if err := tx.First(&d, id).Error; err != nil {
			return err
		}
		if err := mustWoDraft(d.Status); err != nil {
			return err
		}
		if d.BomID == 0 {
			return errors.New("下达前必须指定 BOM")
		}
		// 读 BOM
		var bom manufacturing.MfgBom
		if err := tx.Preload("Items").First(&bom, d.BomID).Error; err != nil {
			return fmt.Errorf("BOM 不存在: %w", err)
		}
		if bom.Status != manufacturing.BomStatusActive {
			return errors.New("BOM 已停用，不可下达")
		}
		if bom.ProductID != d.ProductID {
			return errors.New("BOM 成品与工单成品不一致")
		}
		if len(bom.Items) == 0 {
			return errors.New("BOM 无子项，不可下达")
		}
		// === 5b：工艺路线强制 ===
		if d.RoutingID == 0 {
			return errors.New("下达前必须指定工艺路线")
		}
		var routing manufacturing.MfgRouting
		if err := tx.Preload("Items", func(db *gorm.DB) *gorm.DB { return db.Order("seq") }).
			First(&routing, d.RoutingID).Error; err != nil {
			return fmt.Errorf("工艺路线不存在: %w", err)
		}
		if routing.Status != manufacturing.RoutingStatusActive {
			return errors.New("工艺路线已停用，不可下达")
		}
		if routing.ProductID != d.ProductID {
			return errors.New("工艺路线成品与工单成品不一致")
		}
		if len(routing.Items) == 0 {
			return errors.New("工艺路线无工序，不可下达")
		}
		// 清理旧需求并展开
		if err := tx.Where("work_order_id=?", id).Delete(&manufacturing.MfgWorkOrderItem{}).Error; err != nil {
			return err
		}
		for _, bi := range bom.Items {
			req := bi.Quantity * d.PlanQty * (1 + bi.LossRate)
			item := manufacturing.MfgWorkOrderItem{
				WorkOrderID: id,
				ProductID:   bi.ProductID,
				RequiredQty: req,
				Remark:      bi.Remark,
			}
			if err := tx.Create(&item).Error; err != nil {
				return err
			}
		}
		// === 5b：展开工序快照 ===
		if err := tx.Where("work_order_id=?", id).Delete(&manufacturing.MfgWorkOrderOp{}).Error; err != nil {
			return err
		}
		for _, ri := range routing.Items {
			op := manufacturing.MfgWorkOrderOp{
				WorkOrderID:   id,
				Seq:           ri.Seq,
				Code:          ri.Code,
				Name:          ri.Name,
				WorkstationID: ri.WorkstationID,
				StdHours:      ri.StdHours,
				IsKey:         ri.IsKey,
				PlanQty:       d.PlanQty,
				Status:        manufacturing.OpStatusPending,
			}
			if err := tx.Create(&op).Error; err != nil {
				return err
			}
		}
		return tx.Model(&d).Updates(map[string]interface{}{
			"status":     manufacturing.WoStatusReleased,
			"op_user_id": opUserID,
		}).Error
	})
}

// Close 关闭工单（已完工或人工提前关闭）
func (s *WorkOrderService) Close(id, opUserID uint) error {
	var d manufacturing.MfgWorkOrder
	if err := global.GVA_DB.First(&d, id).Error; err != nil {
		return err
	}
	if d.Status == manufacturing.WoStatusClosed || d.Status == manufacturing.WoStatusVoid {
		return errors.New("状态不允许")
	}
	return global.GVA_DB.Model(&d).Updates(map[string]interface{}{
		"status": manufacturing.WoStatusClosed, "op_user_id": opUserID,
	}).Error
}

// Void 作废工单（仅草稿/已下达，且未领料未完工）
func (s *WorkOrderService) Void(id, opUserID uint) error {
	var d manufacturing.MfgWorkOrder
	if err := global.GVA_DB.First(&d, id).Error; err != nil {
		return err
	}
	if d.Status == manufacturing.WoStatusVoid {
		return errors.New("已作废")
	}
	if d.ActualQty > 0 {
		return errors.New("已有完工记录，不可作废，请先冲销")
	}
	// 检查是否有已执行的领料单
	var cnt int64
	global.GVA_DB.Model(&manufacturing.MfgIssue{}).
		Where("work_order_id=? AND status=?", id, manufacturing.DocStatusDone).Count(&cnt)
	if cnt > 0 {
		return errors.New("已有领料记录，不可作废，请先冲销")
	}
	// 5b：有生效报工则禁止作废
	var rptCnt int64
	global.GVA_DB.Model(&manufacturing.MfgWorkReport{}).
		Where("work_order_id=? AND status=?", id, manufacturing.ReportStatusDone).Count(&rptCnt)
	if rptCnt > 0 {
		return errors.New("已有报工记录，不可作废，请先作废报工")
	}
	return global.GVA_DB.Model(&d).Updates(map[string]interface{}{
		"status": manufacturing.WoStatusVoid, "op_user_id": opUserID,
	}).Error
}

// GetByID 详细
func (s *WorkOrderService) GetByID(id uint) (e manufacturing.MfgWorkOrder, err error) {
	err = global.GVA_DB.Preload("Product").Preload("Warehouse").
		Preload("Bom").Preload("Routing").Preload("Items.Product").First(&e, id).Error
	return
}

// GetOps 工单工序列表（含派工信息）
func (s *WorkOrderService) GetOps(workOrderID uint) (list []manufacturing.MfgWorkOrderOp, err error) {
	err = global.GVA_DB.Preload("Workstation").Preload("AssignedUser").Preload("Team").
		Where("work_order_id=?", workOrderID).Order("seq").Find(&list).Error
	return
}

// GetList 分页
func (s *WorkOrderService) GetList(info WoSearch) (list []manufacturing.MfgWorkOrder, total int64, err error) {
	if info.Page <= 0 {
		info.Page = 1
	}
	if info.PageSize <= 0 {
		info.PageSize = 10
	}
	db := global.GVA_DB.Model(&manufacturing.MfgWorkOrder{})
	if info.Status > 0 {
		db = db.Where("status=?", info.Status)
	}
	if info.DocNo != "" {
		db = db.Where("doc_no LIKE ?", "%"+info.DocNo+"%")
	}
	if info.ProductID > 0 {
		db = db.Where("product_id=?", info.ProductID)
	}
	if err = db.Count(&total).Error; err != nil {
		return
	}
	err = db.Preload("Product").Preload("Warehouse").
		Limit(info.PageSize).Offset((info.Page - 1) * info.PageSize).
		Order("id DESC").Find(&list).Error
	return
}
