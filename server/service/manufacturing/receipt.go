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
// 完工入库服务
// 一张完工单只入库工单上的成品；UnitCost 在 Confirm 时按"实际成本法"算：
//
//   unitCost = 工单累计领料金额 / (工单已完工 + 本次完工)
//
// 简化版：本次完工分摊 (本次/(已完工+本次)) × 累计领料金额
// 实现：直接用 unitCost = (累计领料金额 / 总完工数量) ，
//       回写时 wo.IssuedAmount 不变（领料金额还在那），通过 ActualQty 累加自然算出新均价
// ============================================================
type ReceiptService struct{}

func (s *ReceiptService) Create(e *manufacturing.MfgReceipt) error {
	if e.WorkOrderID == 0 || e.WarehouseID == 0 || e.Quantity <= 0 {
		return errors.New("工单、入库仓、数量必填")
	}
	var wo manufacturing.MfgWorkOrder
	if err := global.GVA_DB.Preload("Items").First(&wo, e.WorkOrderID).Error; err != nil {
		return errors.New("工单不存在")
	}
	if wo.Status != manufacturing.WoStatusReleased && wo.Status != manufacturing.WoStatusInProgress {
		return errors.New("工单状态不允许完工入库（需已下达或生产中）")
	}
	// 超完工校验
	if wo.ActualQty+e.Quantity > wo.PlanQty+0.0001 {
		return fmt.Errorf("超完工：计划%.4f 已完%.4f 本次%.4f",
			wo.PlanQty, wo.ActualQty, e.Quantity)
	}
	// 物料齐套校验（草稿阶段就拦截，UX 更好）
	if len(wo.Items) == 0 {
		return errors.New("工单未下达或无物料需求，不可完工入库")
	}
	theoreticalMax := wo.PlanQty
	for _, woItem := range wo.Items {
		if woItem.RequiredQty <= 0 {
			continue
		}
		perUnit := woItem.RequiredQty / wo.PlanQty
		maxByThis := woItem.IssuedQty / perUnit
		if maxByThis < theoreticalMax {
			theoreticalMax = maxByThis
		}
	}
	if wo.ActualQty+e.Quantity > theoreticalMax+0.0001 {
		return fmt.Errorf("物料不足：当前已领料最多支撑完工%.4f，已完%.4f 本次拟完%.4f；请先领料",
			theoreticalMax, wo.ActualQty, e.Quantity)
	}
	// 5b：工序联动校验
	if err := checkLastOpDone(global.GVA_DB, wo.ID, wo.ActualQty+e.Quantity); err != nil {
		return err
	}
	e.Status = manufacturing.DocStatusDraft
	e.UnitCost = 0
	e.TotalAmount = 0
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		no, err := seqService.NextDocNo(tx, inventory.BizMfgReceipt)
		if err != nil {
			return err
		}
		e.DocNo = no
		return tx.Create(e).Error
	})
}

func (s *ReceiptService) Update(e *manufacturing.MfgReceipt) error {
	var old manufacturing.MfgReceipt
	if err := global.GVA_DB.First(&old, e.ID).Error; err != nil {
		return err
	}
	if err := mustDocDraft(old.Status); err != nil {
		return err
	}
	e.DocNo = old.DocNo
	e.Status = old.Status
	e.WorkOrderID = old.WorkOrderID
	e.UnitCost = 0
	e.TotalAmount = 0
	e.CreatedAt = old.CreatedAt
	return global.GVA_DB.Save(e).Error
}

func (s *ReceiptService) Delete(id uint) error {
	var d manufacturing.MfgReceipt
	if err := global.GVA_DB.First(&d, id).Error; err != nil {
		return err
	}
	if err := mustDocDraft(d.Status); err != nil {
		return err
	}
	return global.GVA_DB.Delete(&manufacturing.MfgReceipt{}, id).Error
}

// Confirm 完工入库
//
// 成本核算（实际成本法）：
//   单位成本 = 工单累计领料金额 / (已完工数量 + 本次完工数量)
//
// 注意：这是简化的"按比例分摊"，最终工单完全完工后总成本 ≈ 累计领料金额。
// 若工单尚未领完料就完工，单位成本会偏低；若超领料，单位成本偏高。
// 实际企业可在阶段5 加"标准成本/差异分析"。
func (s *ReceiptService) Confirm(id, opUserID uint) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var d manufacturing.MfgReceipt
		if err := tx.First(&d, id).Error; err != nil {
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
			return errors.New("工单状态不允许完工入库")
		}
		if wo.ActualQty+d.Quantity > wo.PlanQty+0.0001 {
			return fmt.Errorf("超完工：计划%.4f 已完%.4f 本次%.4f",
				wo.PlanQty, wo.ActualQty, d.Quantity)
		}

		// === 物料齐套校验：领料量必须能支撑本次完工数 ===
		// 理论最大可完工数 = min(每种物料的 IssuedQty / 单位用量) ，
		// 单位用量 = RequiredQty / PlanQty （含损耗率）
		if len(wo.Items) == 0 {
			return errors.New("工单未下达或无物料需求，不可完工入库")
		}
		theoreticalMax := wo.PlanQty // 上界
		for _, woItem := range wo.Items {
			if woItem.RequiredQty <= 0 {
				continue
			}
			perUnit := woItem.RequiredQty / wo.PlanQty
			maxByThis := woItem.IssuedQty / perUnit
			if maxByThis < theoreticalMax {
				theoreticalMax = maxByThis
			}
		}
		if wo.ActualQty+d.Quantity > theoreticalMax+0.0001 {
			return fmt.Errorf("物料不足：当前已领料最多支撑完工%.4f，已完%.4f 本次拟完%.4f；请先领料",
				theoreticalMax, wo.ActualQty, d.Quantity)
		}
		// 5b：工序联动校验
		if err := checkLastOpDone(tx, wo.ID, wo.ActualQty+d.Quantity); err != nil {
			return err
		}

		// 计算单位成本
		totalQtyAfter := wo.ActualQty + d.Quantity
		var unitCost float64
		if wo.IssuedAmount > 0 && totalQtyAfter > 0 {
			unitCost = wo.IssuedAmount / totalQtyAfter
		}
		totalAmount := unitCost * d.Quantity

		// 调 StockEngine 入库
		if err := stockEngine.Apply(tx, invsvc.StockMoveParam{
			BizType:     inventory.BizMfgReceipt,
			BizNo:       d.DocNo,
			BizItemID:   d.ID,
			Direction:   inventory.DirIn,
			ProductID:   wo.ProductID,
			WarehouseID: d.WarehouseID,
			LocationID:  d.LocationID,
			Quantity:    d.Quantity,
			UnitCost:    unitCost,
			OpUserID:    opUserID,
		}); err != nil {
			return err
		}

		// 更新单据头
		if err := tx.Model(&d).Updates(map[string]interface{}{
			"status":       manufacturing.DocStatusDone,
			"unit_cost":    unitCost,
			"total_amount": totalAmount,
			"op_user_id":   opUserID,
		}).Error; err != nil {
			return err
		}

		// 回写工单完工数 + 状态
		newActual := wo.ActualQty + d.Quantity
		updates := map[string]interface{}{"actual_qty": newActual}
		if newActual >= wo.PlanQty-0.0001 {
			updates["status"] = manufacturing.WoStatusFinished
		} else if wo.Status == manufacturing.WoStatusReleased {
			updates["status"] = manufacturing.WoStatusInProgress
		}
		return tx.Model(&wo).Updates(updates).Error
	})
}

// Void 作废已执行的完工单：反向出库 + 回退完工数
func (s *ReceiptService) Void(id, opUserID uint) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var d manufacturing.MfgReceipt
		if err := tx.First(&d, id).Error; err != nil {
			return err
		}
		if err := mustDocDone(d.Status); err != nil {
			return err
		}
		var wo manufacturing.MfgWorkOrder
		if err := tx.First(&wo, d.WorkOrderID).Error; err != nil {
			return err
		}
		// 反向：把成品出库
		if err := stockEngine.Apply(tx, invsvc.StockMoveParam{
			BizType:     inventory.BizMfgReceipt,
			BizNo:       d.DocNo,
			BizItemID:   d.ID,
			Direction:   inventory.DirOut,
			ProductID:   wo.ProductID,
			WarehouseID: d.WarehouseID,
			LocationID:  d.LocationID,
			Quantity:    d.Quantity,
			Remark:      "作废冲销",
			OpUserID:    opUserID,
		}); err != nil {
			return err
		}
		// 回退完工数
		newActual := wo.ActualQty - d.Quantity
		if newActual < 0 {
			newActual = 0
		}
		updates := map[string]interface{}{"actual_qty": newActual}
		// 状态回退：从已完工 → 生产中
		if wo.Status == manufacturing.WoStatusFinished {
			updates["status"] = manufacturing.WoStatusInProgress
		}
		if err := tx.Model(&wo).Updates(updates).Error; err != nil {
			return err
		}
		return tx.Model(&d).Updates(map[string]interface{}{
			"status": manufacturing.DocStatusVoid, "op_user_id": opUserID,
		}).Error
	})
}

func (s *ReceiptService) GetByID(id uint) (e manufacturing.MfgReceipt, err error) {
	err = global.GVA_DB.Preload("Warehouse").Preload("Location").
		Preload("WorkOrder.Product").First(&e, id).Error
	return
}

func (s *ReceiptService) GetList(info DocSearch) (list []manufacturing.MfgReceipt, total int64, err error) {
	if info.Page <= 0 {
		info.Page = 1
	}
	if info.PageSize <= 0 {
		info.PageSize = 10
	}
	db := applyDocFilter(global.GVA_DB.Model(&manufacturing.MfgReceipt{}), info)
	if err = db.Count(&total).Error; err != nil {
		return
	}
	err = db.Preload("Warehouse").Preload("WorkOrder.Product").
		Limit(info.PageSize).Offset((info.Page - 1) * info.PageSize).
		Order("id DESC").Find(&list).Error
	return
}
