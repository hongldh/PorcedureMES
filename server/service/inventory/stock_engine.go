package inventory

import (
	"errors"
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/model/inventory"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// StockEngine 库存核心引擎
//
// 所有库存变动都走 Apply：
//   1. 行锁锁住 (product, warehouse, location) 余额行
//   2. 计算变动前/后数量与移动加权平均
//   3. 写一条流水（InvStockFlow）
//   4. upsert 余额表（InvStock）
//
// 加权平均规则：
//   入库：new_avg = (old_qty*old_avg + in_qty*in_cost) / (old_qty + in_qty)，old_qty<=0 时直接取 in_cost
//   出库：均价不变，按 old_avg 出
type StockEngine struct{}

// StockMoveParam 单次库存变动参数
type StockMoveParam struct {
	BizType     string  // 业务类型（常量，BizPurchaseIn 等）
	BizNo       string  // 业务单号（已生成的 DocNo）
	BizItemID   uint    // 业务明细ID
	Direction   int8    // DirIn / DirOut
	ProductID   uint
	WarehouseID uint
	LocationID  uint
	Quantity    float64 // 正数
	UnitCost    float64 // 入库时为采购单价；出库时可不传，引擎用现有均价
	Remark      string
	OpUserID    uint
}

// Apply 在已有事务 tx 中执行一次库存变动
func (e *StockEngine) Apply(tx *gorm.DB, p StockMoveParam) error {
	if p.Quantity <= 0 {
		return errors.New("数量必须>0")
	}
	if p.Direction != inventory.DirIn && p.Direction != inventory.DirOut {
		return errors.New("方向不合法")
	}

	// 1. 行锁查/建余额行
	var stock inventory.InvStock
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("product_id=? AND warehouse_id=? AND location_id=?",
			p.ProductID, p.WarehouseID, p.LocationID).
		First(&stock).Error
	if err == gorm.ErrRecordNotFound {
		stock = inventory.InvStock{
			ProductID:   p.ProductID,
			WarehouseID: p.WarehouseID,
			LocationID:  p.LocationID,
			Quantity:    0,
			AvgCost:     0,
		}
		if err = tx.Create(&stock).Error; err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	beforeQty := stock.Quantity
	beforeCost := stock.AvgCost
	var afterQty, afterCost, totalAmount, unitCost float64

	if p.Direction == inventory.DirIn {
		unitCost = p.UnitCost
		// 若上游不传成本（如销售退货回库），按当前均价回库，避免污染库存成本
		if unitCost <= 0 {
			unitCost = beforeCost
		}
		afterQty = beforeQty + p.Quantity
		if beforeQty <= 0 {
			afterCost = unitCost
		} else {
			afterCost = (beforeQty*beforeCost + p.Quantity*unitCost) / afterQty
		}
		totalAmount = p.Quantity * unitCost
	} else { // 出库
		if beforeQty < p.Quantity {
			return fmt.Errorf("库存不足: 商品%d 仓%d 位%d 现存%.4f 需出%.4f",
				p.ProductID, p.WarehouseID, p.LocationID, beforeQty, p.Quantity)
		}
		unitCost = beforeCost // 出库按当前均价
		afterQty = beforeQty - p.Quantity
		afterCost = beforeCost
		totalAmount = p.Quantity * unitCost
	}

	// 2. 写流水
	flow := inventory.InvStockFlow{
		BizType:     p.BizType,
		BizNo:       p.BizNo,
		BizItemID:   p.BizItemID,
		Direction:   p.Direction,
		ProductID:   p.ProductID,
		WarehouseID: p.WarehouseID,
		LocationID:  p.LocationID,
		Quantity:    p.Quantity,
		UnitCost:    unitCost,
		TotalAmount: totalAmount,
		BeforeQty:   beforeQty,
		AfterQty:    afterQty,
		BeforeCost:  beforeCost,
		AfterCost:   afterCost,
		Remark:      p.Remark,
		OpUserID:    p.OpUserID,
	}
	if err = tx.Create(&flow).Error; err != nil {
		return err
	}

	// 3. 更新余额
	stock.Quantity = afterQty
	stock.AvgCost = afterCost
	return tx.Save(&stock).Error
}
