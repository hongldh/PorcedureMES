package manufacturing

import (
	"errors"
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/inventory"
	"github.com/flipped-aurora/gin-vue-admin/server/model/manufacturing"
	invsvc "github.com/flipped-aurora/gin-vue-admin/server/service/inventory"
	"gorm.io/gorm"
)

// ============================================================
// 领料单服务
// 关键校验：累计领料 ≤ 物料需求量；Confirm 时调 StockEngine 出库
// ============================================================
type IssueService struct{}

func (s *IssueService) Create(e *manufacturing.MfgIssue) error {
	if e.WorkOrderID == 0 || e.WarehouseID == 0 {
		return errors.New("工单与仓库必填")
	}
	// 校验工单状态
	var wo manufacturing.MfgWorkOrder
	if err := global.GVA_DB.First(&wo, e.WorkOrderID).Error; err != nil {
		return fmt.Errorf("工单不存在")
	}
	if wo.Status != manufacturing.WoStatusReleased && wo.Status != manufacturing.WoStatusInProgress {
		return errors.New("工单状态不允许领料（需已下达或生产中）")
	}
	e.Status = manufacturing.DocStatusDraft
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		no, err := seqService.NextDocNo(tx, inventory.BizMfgIssue)
		if err != nil {
			return err
		}
		e.DocNo = no
		// 草稿态金额由 Confirm 计算（取实际出库均价），先按 0 存
		for i := range e.Items {
			e.Items[i].ID = 0
			e.Items[i].Amount = 0
		}
		return tx.Create(e).Error
	})
}

func (s *IssueService) Update(e *manufacturing.MfgIssue) error {
	var old manufacturing.MfgIssue
	if err := global.GVA_DB.First(&old, e.ID).Error; err != nil {
		return err
	}
	if err := mustDocDraft(old.Status); err != nil {
		return err
	}
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("issue_id=?", e.ID).Delete(&manufacturing.MfgIssueItem{}).Error; err != nil {
			return err
		}
		for i := range e.Items {
			e.Items[i].ID = 0
			e.Items[i].IssueID = e.ID
			e.Items[i].Amount = 0
		}
		e.DocNo = old.DocNo
		e.Status = old.Status
		e.WorkOrderID = old.WorkOrderID
		e.CreatedAt = old.CreatedAt
		return tx.Session(&gorm.Session{FullSaveAssociations: true}).Save(e).Error
	})
}

func (s *IssueService) Delete(id uint) error {
	var d manufacturing.MfgIssue
	if err := global.GVA_DB.First(&d, id).Error; err != nil {
		return err
	}
	if err := mustDocDraft(d.Status); err != nil {
		return err
	}
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("issue_id=?", id).Delete(&manufacturing.MfgIssueItem{}).Error; err != nil {
			return err
		}
		return tx.Delete(&manufacturing.MfgIssue{}, id).Error
	})
}

// Confirm 草稿→已执行
//
// 流程：
//  1. 校验工单状态、超领（累计领料 ≤ 需求量）
//  2. 对每行调 StockEngine.Apply(BizMfgIssue, DirOut)，UnitCost 由引擎按当前均价填入
//  3. 回写 mfg_work_order_item.IssuedQty、mfg_work_order.IssuedAmount
//  4. 工单状态：已下达 → 生产中
func (s *IssueService) Confirm(id, opUserID uint) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var d manufacturing.MfgIssue
		if err := tx.Preload("Items").First(&d, id).Error; err != nil {
			return err
		}
		if err := mustDocDraft(d.Status); err != nil {
			return err
		}
		var wo manufacturing.MfgWorkOrder
		if err := tx.Preload("Items").First(&wo, d.WorkOrderID).Error; err != nil {
			return err
		}
		if wo.Status != manufacturing.WoStatusReleased && wo.Status != manufacturing.WoStatusInProgress {
			return errors.New("工单状态不允许领料")
		}

		// 物料需求 map
		woItemByPid := make(map[uint]*manufacturing.MfgWorkOrderItem)
		for i := range wo.Items {
			woItemByPid[wo.Items[i].ProductID] = &wo.Items[i]
		}

		var totalAmount float64
		for i := range d.Items {
			it := &d.Items[i]
			if it.Quantity <= 0 {
				return errors.New("领料数量必须>0")
			}
			woItem, ok := woItemByPid[it.ProductID]
			if !ok {
				return fmt.Errorf("物料 %d 不在工单需求中", it.ProductID)
			}
			// 超领校验
			if woItem.IssuedQty+it.Quantity > woItem.RequiredQty+0.0001 {
				return fmt.Errorf("物料 %d 超领：需求%.4f 已领%.4f 本次%.4f",
					it.ProductID, woItem.RequiredQty, woItem.IssuedQty, it.Quantity)
			}
			// 调 StockEngine 出库；UnitCost 不传 → 引擎按当前均价
			if err := stockEngine.Apply(tx, invsvc.StockMoveParam{
				BizType:     inventory.BizMfgIssue,
				BizNo:       d.DocNo,
				BizItemID:   it.ID,
				Direction:   inventory.DirOut,
				ProductID:   it.ProductID,
				WarehouseID: d.WarehouseID,
				LocationID:  it.LocationID,
				Quantity:    it.Quantity,
				OpUserID:    opUserID,
			}); err != nil {
				return err
			}
			// 取出库后的实际单价（引擎已按当前均价出库）：从最新流水反查
			var lastFlow inventory.InvStockFlow
			if err := tx.Where("biz_type=? AND biz_no=? AND biz_item_id=?",
				inventory.BizMfgIssue, d.DocNo, it.ID).
				Order("id DESC").First(&lastFlow).Error; err != nil {
				return err
			}
			it.UnitCost = lastFlow.UnitCost
			it.Amount = lastFlow.TotalAmount
			if err := tx.Save(it).Error; err != nil {
				return err
			}
			totalAmount += it.Amount

			// 累计领料数
			woItem.IssuedQty += it.Quantity
			if err := tx.Save(woItem).Error; err != nil {
				return err
			}
		}

		// 更新单据头
		if err := tx.Model(&d).Updates(map[string]interface{}{
			"status":       manufacturing.DocStatusDone,
			"total_amount": totalAmount,
			"op_user_id":   opUserID,
		}).Error; err != nil {
			return err
		}

		// 回写工单累计金额 + 状态
		updates := map[string]interface{}{
			"issued_amount": wo.IssuedAmount + totalAmount,
		}
		if wo.Status == manufacturing.WoStatusReleased {
			updates["status"] = manufacturing.WoStatusInProgress
		}
		return tx.Model(&wo).Updates(updates).Error
	})
}

// Void 作废已执行的领料单：反向入库 + 回退累计领料数
func (s *IssueService) Void(id, opUserID uint) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var d manufacturing.MfgIssue
		if err := tx.Preload("Items").First(&d, id).Error; err != nil {
			return err
		}
		if err := mustDocDone(d.Status); err != nil {
			return err
		}
		var wo manufacturing.MfgWorkOrder
		if err := tx.Preload("Items").First(&wo, d.WorkOrderID).Error; err != nil {
			return err
		}
		woItemByPid := make(map[uint]*manufacturing.MfgWorkOrderItem)
		for i := range wo.Items {
			woItemByPid[wo.Items[i].ProductID] = &wo.Items[i]
		}
		for _, it := range d.Items {
			// 反向入库（按当时出库成本回库）
			if err := stockEngine.Apply(tx, invsvc.StockMoveParam{
				BizType:     inventory.BizMfgIssue,
				BizNo:       d.DocNo,
				BizItemID:   it.ID,
				Direction:   inventory.DirIn,
				ProductID:   it.ProductID,
				WarehouseID: d.WarehouseID,
				LocationID:  it.LocationID,
				Quantity:    it.Quantity,
				UnitCost:    it.UnitCost,
				Remark:      "作废冲销",
				OpUserID:    opUserID,
			}); err != nil {
				return err
			}
			if woItem, ok := woItemByPid[it.ProductID]; ok {
				woItem.IssuedQty -= it.Quantity
				if woItem.IssuedQty < 0 {
					woItem.IssuedQty = 0
				}
				if err := tx.Save(woItem).Error; err != nil {
					return err
				}
			}
		}
		// 工单累计金额扣回
		newAmount := wo.IssuedAmount - d.TotalAmount
		if newAmount < 0 {
			newAmount = 0
		}
		if err := tx.Model(&wo).Update("issued_amount", newAmount).Error; err != nil {
			return err
		}
		return tx.Model(&d).Updates(map[string]interface{}{
			"status": manufacturing.DocStatusVoid, "op_user_id": opUserID,
		}).Error
	})
}

func (s *IssueService) GetByID(id uint) (e manufacturing.MfgIssue, err error) {
	err = global.GVA_DB.Preload("Warehouse").Preload("WorkOrder.Product").
		Preload("Items.Product").Preload("Items.Location").First(&e, id).Error
	return
}

func (s *IssueService) GetList(info DocSearch) (list []manufacturing.MfgIssue, total int64, err error) {
	if info.Page <= 0 {
		info.Page = 1
	}
	if info.PageSize <= 0 {
		info.PageSize = 10
	}
	db := applyDocFilter(global.GVA_DB.Model(&manufacturing.MfgIssue{}), info)
	if err = db.Count(&total).Error; err != nil {
		return
	}
	err = db.Preload("Warehouse").Preload("WorkOrder").
		Limit(info.PageSize).Offset((info.Page - 1) * info.PageSize).
		Order("id DESC").Find(&list).Error
	return
}
