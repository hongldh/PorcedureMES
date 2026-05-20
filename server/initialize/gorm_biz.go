package initialize

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/inventory"
)

func bizModel() error {
	db := global.GVA_DB
	err := db.AutoMigrate(
		// ===== 进销存：主数据 =====
		inventory.InvCategory{},
		inventory.InvUnit{},
		inventory.InvProduct{},
		inventory.InvWarehouse{},
		inventory.InvLocation{},
		inventory.InvSupplier{},
		inventory.InvCustomer{},

		// ===== 进销存：库存核心 =====
		inventory.InvStock{},
		inventory.InvStockFlow{},
		inventory.InvDocSequence{},

		// ===== 进销存：业务单据 =====
		inventory.InvPurchaseOrder{},
		inventory.InvPurchaseOrderItem{},
		inventory.InvPurchaseReturn{},
		inventory.InvPurchaseReturnItem{},
		inventory.InvSalesOrder{},
		inventory.InvSalesOrderItem{},
		inventory.InvSalesReturn{},
		inventory.InvSalesReturnItem{},
		inventory.InvTransfer{},
		inventory.InvTransferItem{},
		inventory.InvStocktake{},
		inventory.InvStocktakeItem{},
	)
	if err != nil {
		return err
	}
	return nil
}
