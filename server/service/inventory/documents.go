package inventory

import (
	"errors"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/inventory"
	"gorm.io/gorm"
)

// ============================================================
// 单据通用：DocSearch（按状态、日期、关键字、单号搜）
// ============================================================
type DocSearch struct {
	request.PageInfo
	Status    int8   `json:"status" form:"status"`
	DateStart string `json:"dateStart" form:"dateStart"`
	DateEnd   string `json:"dateEnd" form:"dateEnd"`
	DocNo     string `json:"docNo" form:"docNo"`
}

func applyDocFilter(db *gorm.DB, s DocSearch) *gorm.DB {
	if s.Status > 0 {
		db = db.Where("status = ?", s.Status)
	}
	if s.DateStart != "" {
		db = db.Where("biz_date >= ?", s.DateStart)
	}
	if s.DateEnd != "" {
		db = db.Where("biz_date <= ?", s.DateEnd)
	}
	if s.DocNo != "" {
		db = db.Where("doc_no LIKE ?", "%"+s.DocNo+"%")
	}
	return db
}

// 通用：单据状态校验
func mustDraft(status int8) error {
	if status != inventory.DocStatusDraft {
		return errors.New("仅草稿状态单据可执行此操作")
	}
	return nil
}
func mustDone(status int8) error {
	if status != inventory.DocStatusDone {
		return errors.New("仅已执行单据可作废")
	}
	return nil
}

// 包内共享引擎/序列服务
var (
	stockEngine = &StockEngine{}
	seqService  = &SequenceService{}
)

// ============================================================
// 1. 采购入库
// ============================================================
type PurchaseOrderService struct{}

func (s *PurchaseOrderService) Create(e *inventory.InvPurchaseOrder) error {
	e.Status = inventory.DocStatusDraft
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		no, err := seqService.NextDocNo(tx, inventory.BizPurchaseIn)
		if err != nil {
			return err
		}
		e.DocNo = no
		var total float64
		for i := range e.Items {
			it := &e.Items[i]
			it.Amount = it.Quantity * it.UnitPrice
			total += it.Amount
		}
		e.TotalAmount = total
		return tx.Create(e).Error
	})
}

func (s *PurchaseOrderService) Update(e *inventory.InvPurchaseOrder) error {
	var old inventory.InvPurchaseOrder
	if err := global.GVA_DB.First(&old, e.ID).Error; err != nil {
		return err
	}
	if err := mustDraft(old.Status); err != nil {
		return err
	}
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("order_id = ?", e.ID).Delete(&inventory.InvPurchaseOrderItem{}).Error; err != nil {
			return err
		}
		var total float64
		for i := range e.Items {
			it := &e.Items[i]
			it.ID = 0
			it.OrderID = e.ID
			it.Amount = it.Quantity * it.UnitPrice
			total += it.Amount
		}
		e.TotalAmount = total
		e.DocNo = old.DocNo
		e.Status = old.Status
		e.CreatedAt = old.CreatedAt
		return tx.Session(&gorm.Session{FullSaveAssociations: true}).Save(e).Error
	})
}

func (s *PurchaseOrderService) Delete(id uint) error {
	var d inventory.InvPurchaseOrder
	if err := global.GVA_DB.First(&d, id).Error; err != nil {
		return err
	}
	if err := mustDraft(d.Status); err != nil {
		return err
	}
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("order_id=?", id).Delete(&inventory.InvPurchaseOrderItem{}).Error; err != nil {
			return err
		}
		return tx.Delete(&inventory.InvPurchaseOrder{}, id).Error
	})
}

// Confirm 草稿 → 已执行：写库存流水 + 更新余额
func (s *PurchaseOrderService) Confirm(id uint, opUserID uint) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var d inventory.InvPurchaseOrder
		if err := tx.Preload("Items").First(&d, id).Error; err != nil {
			return err
		}
		if err := mustDraft(d.Status); err != nil {
			return err
		}
		for _, it := range d.Items {
			if err := stockEngine.Apply(tx, StockMoveParam{
				BizType:     inventory.BizPurchaseIn,
				BizNo:       d.DocNo,
				BizItemID:   it.ID,
				Direction:   inventory.DirIn,
				ProductID:   it.ProductID,
				WarehouseID: d.WarehouseID,
				LocationID:  it.LocationID,
				Quantity:    it.Quantity,
				UnitCost:    it.UnitPrice,
				OpUserID:    opUserID,
			}); err != nil {
				return err
			}
		}
		return tx.Model(&d).Updates(map[string]interface{}{
			"status":     inventory.DocStatusDone,
			"op_user_id": opUserID,
		}).Error
	})
}

// Void 已执行 → 作废：反向冲销库存
func (s *PurchaseOrderService) Void(id uint, opUserID uint) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var d inventory.InvPurchaseOrder
		if err := tx.Preload("Items").First(&d, id).Error; err != nil {
			return err
		}
		if err := mustDone(d.Status); err != nil {
			return err
		}
		for _, it := range d.Items {
			if err := stockEngine.Apply(tx, StockMoveParam{
				BizType:     inventory.BizPurchaseIn,
				BizNo:       d.DocNo,
				BizItemID:   it.ID,
				Direction:   inventory.DirOut, // 反向
				ProductID:   it.ProductID,
				WarehouseID: d.WarehouseID,
				LocationID:  it.LocationID,
				Quantity:    it.Quantity,
				Remark:      "作废冲销",
				OpUserID:    opUserID,
			}); err != nil {
				return err
			}
		}
		return tx.Model(&d).Update("status", inventory.DocStatusVoid).Error
	})
}

func (s *PurchaseOrderService) GetByID(id uint) (e inventory.InvPurchaseOrder, err error) {
	err = global.GVA_DB.Preload("Supplier").Preload("Warehouse").
		Preload("Items.Product").Preload("Items.Location").First(&e, id).Error
	return
}

func (s *PurchaseOrderService) GetList(info DocSearch) (list []inventory.InvPurchaseOrder, total int64, err error) {
	db := applyDocFilter(global.GVA_DB.Model(&inventory.InvPurchaseOrder{}), info)
	if info.Keyword != "" {
		db = db.Where("doc_no LIKE ?", "%"+info.Keyword+"%")
	}
	if err = db.Count(&total).Error; err != nil {
		return
	}
	err = db.Preload("Supplier").Preload("Warehouse").
		Scopes(info.PageInfo.Paginate()).Order("id desc").Find(&list).Error
	return
}

// ============================================================
// 2. 采购退货 (出库)
// ============================================================
type PurchaseReturnService struct{}

func (s *PurchaseReturnService) Create(e *inventory.InvPurchaseReturn) error {
	e.Status = inventory.DocStatusDraft
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		no, err := seqService.NextDocNo(tx, inventory.BizPurchaseReturn)
		if err != nil {
			return err
		}
		e.DocNo = no
		var total float64
		for i := range e.Items {
			it := &e.Items[i]
			it.Amount = it.Quantity * it.UnitPrice
			total += it.Amount
		}
		e.TotalAmount = total
		return tx.Create(e).Error
	})
}

func (s *PurchaseReturnService) Update(e *inventory.InvPurchaseReturn) error {
	var old inventory.InvPurchaseReturn
	if err := global.GVA_DB.First(&old, e.ID).Error; err != nil {
		return err
	}
	if err := mustDraft(old.Status); err != nil {
		return err
	}
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("order_id = ?", e.ID).Delete(&inventory.InvPurchaseReturnItem{}).Error; err != nil {
			return err
		}
		var total float64
		for i := range e.Items {
			it := &e.Items[i]
			it.ID = 0
			it.OrderID = e.ID
			it.Amount = it.Quantity * it.UnitPrice
			total += it.Amount
		}
		e.TotalAmount = total
		e.DocNo = old.DocNo
		e.Status = old.Status
		e.CreatedAt = old.CreatedAt
		return tx.Session(&gorm.Session{FullSaveAssociations: true}).Save(e).Error
	})
}

func (s *PurchaseReturnService) Delete(id uint) error {
	var d inventory.InvPurchaseReturn
	if err := global.GVA_DB.First(&d, id).Error; err != nil {
		return err
	}
	if err := mustDraft(d.Status); err != nil {
		return err
	}
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("order_id=?", id).Delete(&inventory.InvPurchaseReturnItem{}).Error; err != nil {
			return err
		}
		return tx.Delete(&inventory.InvPurchaseReturn{}, id).Error
	})
}

func (s *PurchaseReturnService) Confirm(id, opUserID uint) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var d inventory.InvPurchaseReturn
		if err := tx.Preload("Items").First(&d, id).Error; err != nil {
			return err
		}
		if err := mustDraft(d.Status); err != nil {
			return err
		}
		for _, it := range d.Items {
			if err := stockEngine.Apply(tx, StockMoveParam{
				BizType: inventory.BizPurchaseReturn, BizNo: d.DocNo, BizItemID: it.ID,
				Direction: inventory.DirOut, ProductID: it.ProductID,
				WarehouseID: d.WarehouseID, LocationID: it.LocationID,
				Quantity: it.Quantity, OpUserID: opUserID,
			}); err != nil {
				return err
			}
		}
		return tx.Model(&d).Updates(map[string]interface{}{"status": inventory.DocStatusDone, "op_user_id": opUserID}).Error
	})
}

func (s *PurchaseReturnService) Void(id, opUserID uint) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var d inventory.InvPurchaseReturn
		if err := tx.Preload("Items").First(&d, id).Error; err != nil {
			return err
		}
		if err := mustDone(d.Status); err != nil {
			return err
		}
		for _, it := range d.Items {
			if err := stockEngine.Apply(tx, StockMoveParam{
				BizType: inventory.BizPurchaseReturn, BizNo: d.DocNo, BizItemID: it.ID,
				Direction: inventory.DirIn, ProductID: it.ProductID,
				WarehouseID: d.WarehouseID, LocationID: it.LocationID,
				Quantity: it.Quantity, UnitCost: 0, // 冲销按当前均价回库
				Remark: "作废冲销", OpUserID: opUserID,
			}); err != nil {
				return err
			}
		}
		return tx.Model(&d).Update("status", inventory.DocStatusVoid).Error
	})
}

func (s *PurchaseReturnService) GetByID(id uint) (e inventory.InvPurchaseReturn, err error) {
	err = global.GVA_DB.Preload("Supplier").Preload("Warehouse").
		Preload("Items.Product").Preload("Items.Location").First(&e, id).Error
	return
}

func (s *PurchaseReturnService) GetList(info DocSearch) (list []inventory.InvPurchaseReturn, total int64, err error) {
	db := applyDocFilter(global.GVA_DB.Model(&inventory.InvPurchaseReturn{}), info)
	if err = db.Count(&total).Error; err != nil {
		return
	}
	err = db.Preload("Supplier").Preload("Warehouse").
		Scopes(info.PageInfo.Paginate()).Order("id desc").Find(&list).Error
	return
}

// ============================================================
// 3. 销售出库
// ============================================================
type SalesOrderService struct{}

func (s *SalesOrderService) Create(e *inventory.InvSalesOrder) error {
	e.Status = inventory.DocStatusDraft
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		no, err := seqService.NextDocNo(tx, inventory.BizSalesOut)
		if err != nil {
			return err
		}
		e.DocNo = no
		var total float64
		for i := range e.Items {
			it := &e.Items[i]
			it.Amount = it.Quantity * it.UnitPrice
			total += it.Amount
		}
		e.TotalAmount = total
		return tx.Create(e).Error
	})
}

func (s *SalesOrderService) Update(e *inventory.InvSalesOrder) error {
	var old inventory.InvSalesOrder
	if err := global.GVA_DB.First(&old, e.ID).Error; err != nil {
		return err
	}
	if err := mustDraft(old.Status); err != nil {
		return err
	}
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("order_id = ?", e.ID).Delete(&inventory.InvSalesOrderItem{}).Error; err != nil {
			return err
		}
		var total float64
		for i := range e.Items {
			it := &e.Items[i]
			it.ID = 0
			it.OrderID = e.ID
			it.Amount = it.Quantity * it.UnitPrice
			total += it.Amount
		}
		e.TotalAmount = total
		e.DocNo = old.DocNo
		e.Status = old.Status
		e.CreatedAt = old.CreatedAt
		return tx.Session(&gorm.Session{FullSaveAssociations: true}).Save(e).Error
	})
}

func (s *SalesOrderService) Delete(id uint) error {
	var d inventory.InvSalesOrder
	if err := global.GVA_DB.First(&d, id).Error; err != nil {
		return err
	}
	if err := mustDraft(d.Status); err != nil {
		return err
	}
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("order_id=?", id).Delete(&inventory.InvSalesOrderItem{}).Error; err != nil {
			return err
		}
		return tx.Delete(&inventory.InvSalesOrder{}, id).Error
	})
}

func (s *SalesOrderService) Confirm(id, opUserID uint) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var d inventory.InvSalesOrder
		if err := tx.Preload("Items").First(&d, id).Error; err != nil {
			return err
		}
		if err := mustDraft(d.Status); err != nil {
			return err
		}
		for _, it := range d.Items {
			if err := stockEngine.Apply(tx, StockMoveParam{
				BizType: inventory.BizSalesOut, BizNo: d.DocNo, BizItemID: it.ID,
				Direction: inventory.DirOut, ProductID: it.ProductID,
				WarehouseID: d.WarehouseID, LocationID: it.LocationID,
				Quantity: it.Quantity, OpUserID: opUserID,
			}); err != nil {
				return err
			}
		}
		return tx.Model(&d).Updates(map[string]interface{}{"status": inventory.DocStatusDone, "op_user_id": opUserID}).Error
	})
}

func (s *SalesOrderService) Void(id, opUserID uint) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var d inventory.InvSalesOrder
		if err := tx.Preload("Items").First(&d, id).Error; err != nil {
			return err
		}
		if err := mustDone(d.Status); err != nil {
			return err
		}
		for _, it := range d.Items {
			if err := stockEngine.Apply(tx, StockMoveParam{
				BizType: inventory.BizSalesOut, BizNo: d.DocNo, BizItemID: it.ID,
				Direction: inventory.DirIn, ProductID: it.ProductID,
				WarehouseID: d.WarehouseID, LocationID: it.LocationID,
				Quantity: it.Quantity, UnitCost: 0, // 冲销按当前均价回库
				Remark: "作废冲销", OpUserID: opUserID,
			}); err != nil {
				return err
			}
		}
		return tx.Model(&d).Update("status", inventory.DocStatusVoid).Error
	})
}

func (s *SalesOrderService) GetByID(id uint) (e inventory.InvSalesOrder, err error) {
	err = global.GVA_DB.Preload("Customer").Preload("Warehouse").
		Preload("Items.Product").Preload("Items.Location").First(&e, id).Error
	return
}

func (s *SalesOrderService) GetList(info DocSearch) (list []inventory.InvSalesOrder, total int64, err error) {
	db := applyDocFilter(global.GVA_DB.Model(&inventory.InvSalesOrder{}), info)
	if err = db.Count(&total).Error; err != nil {
		return
	}
	err = db.Preload("Customer").Preload("Warehouse").
		Scopes(info.PageInfo.Paginate()).Order("id desc").Find(&list).Error
	return
}

// ============================================================
// 4. 销售退货 (入库)
// ============================================================
type SalesReturnService struct{}

func (s *SalesReturnService) Create(e *inventory.InvSalesReturn) error {
	e.Status = inventory.DocStatusDraft
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		no, err := seqService.NextDocNo(tx, inventory.BizSalesReturn)
		if err != nil {
			return err
		}
		e.DocNo = no
		var total float64
		for i := range e.Items {
			it := &e.Items[i]
			it.Amount = it.Quantity * it.UnitPrice
			total += it.Amount
		}
		e.TotalAmount = total
		return tx.Create(e).Error
	})
}

func (s *SalesReturnService) Update(e *inventory.InvSalesReturn) error {
	var old inventory.InvSalesReturn
	if err := global.GVA_DB.First(&old, e.ID).Error; err != nil {
		return err
	}
	if err := mustDraft(old.Status); err != nil {
		return err
	}
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("order_id = ?", e.ID).Delete(&inventory.InvSalesReturnItem{}).Error; err != nil {
			return err
		}
		var total float64
		for i := range e.Items {
			it := &e.Items[i]
			it.ID = 0
			it.OrderID = e.ID
			it.Amount = it.Quantity * it.UnitPrice
			total += it.Amount
		}
		e.TotalAmount = total
		e.DocNo = old.DocNo
		e.Status = old.Status
		e.CreatedAt = old.CreatedAt
		return tx.Session(&gorm.Session{FullSaveAssociations: true}).Save(e).Error
	})
}

func (s *SalesReturnService) Delete(id uint) error {
	var d inventory.InvSalesReturn
	if err := global.GVA_DB.First(&d, id).Error; err != nil {
		return err
	}
	if err := mustDraft(d.Status); err != nil {
		return err
	}
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("order_id=?", id).Delete(&inventory.InvSalesReturnItem{}).Error; err != nil {
			return err
		}
		return tx.Delete(&inventory.InvSalesReturn{}, id).Error
	})
}

func (s *SalesReturnService) Confirm(id, opUserID uint) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var d inventory.InvSalesReturn
		if err := tx.Preload("Items").First(&d, id).Error; err != nil {
			return err
		}
		if err := mustDraft(d.Status); err != nil {
			return err
		}
		for _, it := range d.Items {
			if err := stockEngine.Apply(tx, StockMoveParam{
				BizType: inventory.BizSalesReturn, BizNo: d.DocNo, BizItemID: it.ID,
				Direction: inventory.DirIn, ProductID: it.ProductID,
				WarehouseID: d.WarehouseID, LocationID: it.LocationID,
				Quantity: it.Quantity, UnitCost: 0, // 销售退货按当前均价回库（StockEngine 自动取 beforeCost）
				OpUserID: opUserID,
			}); err != nil {
				return err
			}
		}
		return tx.Model(&d).Updates(map[string]interface{}{"status": inventory.DocStatusDone, "op_user_id": opUserID}).Error
	})
}

func (s *SalesReturnService) Void(id, opUserID uint) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var d inventory.InvSalesReturn
		if err := tx.Preload("Items").First(&d, id).Error; err != nil {
			return err
		}
		if err := mustDone(d.Status); err != nil {
			return err
		}
		for _, it := range d.Items {
			if err := stockEngine.Apply(tx, StockMoveParam{
				BizType: inventory.BizSalesReturn, BizNo: d.DocNo, BizItemID: it.ID,
				Direction: inventory.DirOut, ProductID: it.ProductID,
				WarehouseID: d.WarehouseID, LocationID: it.LocationID,
				Quantity: it.Quantity, Remark: "作废冲销", OpUserID: opUserID,
			}); err != nil {
				return err
			}
		}
		return tx.Model(&d).Update("status", inventory.DocStatusVoid).Error
	})
}

func (s *SalesReturnService) GetByID(id uint) (e inventory.InvSalesReturn, err error) {
	err = global.GVA_DB.Preload("Customer").Preload("Warehouse").
		Preload("Items.Product").Preload("Items.Location").First(&e, id).Error
	return
}

func (s *SalesReturnService) GetList(info DocSearch) (list []inventory.InvSalesReturn, total int64, err error) {
	db := applyDocFilter(global.GVA_DB.Model(&inventory.InvSalesReturn{}), info)
	if err = db.Count(&total).Error; err != nil {
		return
	}
	err = db.Preload("Customer").Preload("Warehouse").
		Scopes(info.PageInfo.Paginate()).Order("id desc").Find(&list).Error
	return
}

// ============================================================
// 5. 调拨（一单两动：源仓出 + 目标仓入，按源仓均价做入库单价）
// ============================================================
type TransferService struct{}

func (s *TransferService) Create(e *inventory.InvTransfer) error {
	e.Status = inventory.DocStatusDraft
	if e.FromWarehouse == e.ToWarehouse {
		return errors.New("源仓和目标仓不能相同")
	}
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		no, err := seqService.NextDocNo(tx, inventory.BizTransferOut)
		if err != nil {
			return err
		}
		e.DocNo = no
		return tx.Create(e).Error
	})
}

func (s *TransferService) Update(e *inventory.InvTransfer) error {
	var old inventory.InvTransfer
	if err := global.GVA_DB.First(&old, e.ID).Error; err != nil {
		return err
	}
	if err := mustDraft(old.Status); err != nil {
		return err
	}
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("order_id = ?", e.ID).Delete(&inventory.InvTransferItem{}).Error; err != nil {
			return err
		}
		for i := range e.Items {
			e.Items[i].ID = 0
			e.Items[i].OrderID = e.ID
		}
		e.DocNo = old.DocNo
		e.Status = old.Status
		e.CreatedAt = old.CreatedAt
		return tx.Session(&gorm.Session{FullSaveAssociations: true}).Save(e).Error
	})
}

func (s *TransferService) Delete(id uint) error {
	var d inventory.InvTransfer
	if err := global.GVA_DB.First(&d, id).Error; err != nil {
		return err
	}
	if err := mustDraft(d.Status); err != nil {
		return err
	}
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("order_id=?", id).Delete(&inventory.InvTransferItem{}).Error; err != nil {
			return err
		}
		return tx.Delete(&inventory.InvTransfer{}, id).Error
	})
}

func (s *TransferService) Confirm(id, opUserID uint) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var d inventory.InvTransfer
		if err := tx.Preload("Items").First(&d, id).Error; err != nil {
			return err
		}
		if err := mustDraft(d.Status); err != nil {
			return err
		}
		for _, it := range d.Items {
			// 1) 先查源仓均价（出库前）
			var src inventory.InvStock
			if err := tx.Where("product_id=? AND warehouse_id=? AND location_id=?",
				it.ProductID, d.FromWarehouse, it.FromLocationID).First(&src).Error; err != nil {
				return errors.New("源仓库存不存在")
			}
			cost := src.AvgCost
			// 2) 出
			if err := stockEngine.Apply(tx, StockMoveParam{
				BizType: inventory.BizTransferOut, BizNo: d.DocNo, BizItemID: it.ID,
				Direction: inventory.DirOut, ProductID: it.ProductID,
				WarehouseID: d.FromWarehouse, LocationID: it.FromLocationID,
				Quantity: it.Quantity, OpUserID: opUserID,
			}); err != nil {
				return err
			}
			// 3) 入（用源仓均价）
			if err := stockEngine.Apply(tx, StockMoveParam{
				BizType: inventory.BizTransferIn, BizNo: d.DocNo, BizItemID: it.ID,
				Direction: inventory.DirIn, ProductID: it.ProductID,
				WarehouseID: d.ToWarehouse, LocationID: it.ToLocationID,
				Quantity: it.Quantity, UnitCost: cost, OpUserID: opUserID,
			}); err != nil {
				return err
			}
		}
		return tx.Model(&d).Updates(map[string]interface{}{"status": inventory.DocStatusDone, "op_user_id": opUserID}).Error
	})
}

func (s *TransferService) Void(id, opUserID uint) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var d inventory.InvTransfer
		if err := tx.Preload("Items").First(&d, id).Error; err != nil {
			return err
		}
		if err := mustDone(d.Status); err != nil {
			return err
		}
		for _, it := range d.Items {
			// 反向：目标仓出 → 源仓入
			var dst inventory.InvStock
			if err := tx.Where("product_id=? AND warehouse_id=? AND location_id=?",
				it.ProductID, d.ToWarehouse, it.ToLocationID).First(&dst).Error; err != nil {
				return errors.New("目标仓库存不存在,无法冲销")
			}
			cost := dst.AvgCost
			if err := stockEngine.Apply(tx, StockMoveParam{
				BizType: inventory.BizTransferIn, BizNo: d.DocNo, BizItemID: it.ID,
				Direction: inventory.DirOut, ProductID: it.ProductID,
				WarehouseID: d.ToWarehouse, LocationID: it.ToLocationID,
				Quantity: it.Quantity, Remark: "作废冲销", OpUserID: opUserID,
			}); err != nil {
				return err
			}
			if err := stockEngine.Apply(tx, StockMoveParam{
				BizType: inventory.BizTransferOut, BizNo: d.DocNo, BizItemID: it.ID,
				Direction: inventory.DirIn, ProductID: it.ProductID,
				WarehouseID: d.FromWarehouse, LocationID: it.FromLocationID,
				Quantity: it.Quantity, UnitCost: cost,
				Remark: "作废冲销", OpUserID: opUserID,
			}); err != nil {
				return err
			}
		}
		return tx.Model(&d).Update("status", inventory.DocStatusVoid).Error
	})
}

func (s *TransferService) GetByID(id uint) (e inventory.InvTransfer, err error) {
	err = global.GVA_DB.Preload("FromWh").Preload("ToWh").
		Preload("Items.Product").First(&e, id).Error
	return
}

func (s *TransferService) GetList(info DocSearch) (list []inventory.InvTransfer, total int64, err error) {
	db := applyDocFilter(global.GVA_DB.Model(&inventory.InvTransfer{}), info)
	if err = db.Count(&total).Error; err != nil {
		return
	}
	err = db.Preload("FromWh").Preload("ToWh").
		Scopes(info.PageInfo.Paginate()).Order("id desc").Find(&list).Error
	return
}

// ============================================================
// 6. 盘点（差异 → 一条调整流水）
// ============================================================
type StocktakeService struct{}

func (s *StocktakeService) Create(e *inventory.InvStocktake) error {
	e.Status = inventory.DocStatusDraft
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		no, err := seqService.NextDocNo(tx, inventory.BizStocktakeAdjust)
		if err != nil {
			return err
		}
		e.DocNo = no
		// 创建时按当前账面数填 BookQty
		for i := range e.Items {
			it := &e.Items[i]
			var stock inventory.InvStock
			if err := tx.Where("product_id=? AND warehouse_id=? AND location_id=?",
				it.ProductID, e.WarehouseID, it.LocationID).First(&stock).Error; err == nil {
				it.BookQty = stock.Quantity
			}
			it.DiffQty = it.ActualQty - it.BookQty
		}
		return tx.Create(e).Error
	})
}

func (s *StocktakeService) Update(e *inventory.InvStocktake) error {
	var old inventory.InvStocktake
	if err := global.GVA_DB.First(&old, e.ID).Error; err != nil {
		return err
	}
	if err := mustDraft(old.Status); err != nil {
		return err
	}
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("order_id = ?", e.ID).Delete(&inventory.InvStocktakeItem{}).Error; err != nil {
			return err
		}
		for i := range e.Items {
			it := &e.Items[i]
			it.ID = 0
			it.OrderID = e.ID
			var stock inventory.InvStock
			if err := tx.Where("product_id=? AND warehouse_id=? AND location_id=?",
				it.ProductID, e.WarehouseID, it.LocationID).First(&stock).Error; err == nil {
				it.BookQty = stock.Quantity
			}
			it.DiffQty = it.ActualQty - it.BookQty
		}
		e.DocNo = old.DocNo
		e.Status = old.Status
		e.CreatedAt = old.CreatedAt
		return tx.Session(&gorm.Session{FullSaveAssociations: true}).Save(e).Error
	})
}

func (s *StocktakeService) Delete(id uint) error {
	var d inventory.InvStocktake
	if err := global.GVA_DB.First(&d, id).Error; err != nil {
		return err
	}
	if err := mustDraft(d.Status); err != nil {
		return err
	}
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("order_id=?", id).Delete(&inventory.InvStocktakeItem{}).Error; err != nil {
			return err
		}
		return tx.Delete(&inventory.InvStocktake{}, id).Error
	})
}

func (s *StocktakeService) Confirm(id, opUserID uint) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var d inventory.InvStocktake
		if err := tx.Preload("Items").First(&d, id).Error; err != nil {
			return err
		}
		if err := mustDraft(d.Status); err != nil {
			return err
		}
		for _, it := range d.Items {
			if it.DiffQty == 0 {
				continue
			}
			dir := inventory.DirIn
			qty := it.DiffQty
			cost := 0.0
			if it.DiffQty < 0 { // 盘亏 → 出
				dir = inventory.DirOut
				qty = -it.DiffQty
			} else { // 盘盈 → 入，单价用现有均价
				var stock inventory.InvStock
				if err := tx.Where("product_id=? AND warehouse_id=? AND location_id=?",
					it.ProductID, d.WarehouseID, it.LocationID).First(&stock).Error; err == nil {
					cost = stock.AvgCost
				}
			}
			if err := stockEngine.Apply(tx, StockMoveParam{
				BizType: inventory.BizStocktakeAdjust, BizNo: d.DocNo, BizItemID: it.ID,
				Direction: dir, ProductID: it.ProductID,
				WarehouseID: d.WarehouseID, LocationID: it.LocationID,
				Quantity: qty, UnitCost: cost, OpUserID: opUserID,
			}); err != nil {
				return err
			}
		}
		return tx.Model(&d).Updates(map[string]interface{}{"status": inventory.DocStatusDone, "op_user_id": opUserID}).Error
	})
}

func (s *StocktakeService) Void(id, opUserID uint) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var d inventory.InvStocktake
		if err := tx.Preload("Items").First(&d, id).Error; err != nil {
			return err
		}
		if err := mustDone(d.Status); err != nil {
			return err
		}
		for _, it := range d.Items {
			if it.DiffQty == 0 {
				continue
			}
			// 反向
			dir := inventory.DirOut
			qty := it.DiffQty
			cost := 0.0
			if it.DiffQty < 0 {
				dir = inventory.DirIn
				qty = -it.DiffQty
				var stock inventory.InvStock
				if err := tx.Where("product_id=? AND warehouse_id=? AND location_id=?",
					it.ProductID, d.WarehouseID, it.LocationID).First(&stock).Error; err == nil {
					cost = stock.AvgCost
				}
			}
			if err := stockEngine.Apply(tx, StockMoveParam{
				BizType: inventory.BizStocktakeAdjust, BizNo: d.DocNo, BizItemID: it.ID,
				Direction: dir, ProductID: it.ProductID,
				WarehouseID: d.WarehouseID, LocationID: it.LocationID,
				Quantity: qty, UnitCost: cost, Remark: "作废冲销", OpUserID: opUserID,
			}); err != nil {
				return err
			}
		}
		return tx.Model(&d).Update("status", inventory.DocStatusVoid).Error
	})
}

func (s *StocktakeService) GetByID(id uint) (e inventory.InvStocktake, err error) {
	err = global.GVA_DB.Preload("Warehouse").
		Preload("Items.Product").Preload("Items.Location").First(&e, id).Error
	return
}

func (s *StocktakeService) GetList(info DocSearch) (list []inventory.InvStocktake, total int64, err error) {
	db := applyDocFilter(global.GVA_DB.Model(&inventory.InvStocktake{}), info)
	if err = db.Count(&total).Error; err != nil {
		return
	}
	err = db.Preload("Warehouse").
		Scopes(info.PageInfo.Paginate()).Order("id desc").Find(&list).Error
	return
}

// ============================================================
// 7. 库存查询（余额 + 流水）
// ============================================================
type StockQueryService struct{}

type StockSearch struct {
	request.PageInfo
	ProductID   uint `json:"productId" form:"productId"`
	WarehouseID uint `json:"warehouseId" form:"warehouseId"`
	LocationID  uint `json:"locationId" form:"locationId"`
}

func (s *StockQueryService) BalanceList(info StockSearch) (list []inventory.InvStock, total int64, err error) {
	db := global.GVA_DB.Model(&inventory.InvStock{})
	if info.ProductID > 0 {
		db = db.Where("product_id = ?", info.ProductID)
	}
	if info.WarehouseID > 0 {
		db = db.Where("warehouse_id = ?", info.WarehouseID)
	}
	if info.LocationID > 0 {
		db = db.Where("location_id = ?", info.LocationID)
	}
	if err = db.Count(&total).Error; err != nil {
		return
	}
	err = db.Preload("Product").Preload("Warehouse").Preload("Location").
		Scopes(info.PageInfo.Paginate()).Order("id desc").Find(&list).Error
	return
}

type FlowSearch struct {
	request.PageInfo
	ProductID   uint   `json:"productId" form:"productId"`
	WarehouseID uint   `json:"warehouseId" form:"warehouseId"`
	BizType     string `json:"bizType" form:"bizType"`
	BizNo       string `json:"bizNo" form:"bizNo"`
}

func (s *StockQueryService) FlowList(info FlowSearch) (list []inventory.InvStockFlow, total int64, err error) {
	db := global.GVA_DB.Model(&inventory.InvStockFlow{})
	if info.ProductID > 0 {
		db = db.Where("product_id = ?", info.ProductID)
	}
	if info.WarehouseID > 0 {
		db = db.Where("warehouse_id = ?", info.WarehouseID)
	}
	if info.BizType != "" {
		db = db.Where("biz_type = ?", info.BizType)
	}
	if info.BizNo != "" {
		db = db.Where("biz_no LIKE ?", "%"+info.BizNo+"%")
	}
	if err = db.Count(&total).Error; err != nil {
		return
	}
	err = db.Preload("Product").Preload("Warehouse").
		Scopes(info.PageInfo.Paginate()).Order("id desc").Find(&list).Error
	return
}
