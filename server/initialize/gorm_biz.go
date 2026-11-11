package initialize

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/inventory"
	"github.com/flipped-aurora/gin-vue-admin/server/model/manufacturing"
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

		// ===== MES：制造执行（阶段4）=====
		manufacturing.MfgBom{},
		manufacturing.MfgBomItem{},
		manufacturing.MfgWorkOrder{},
		manufacturing.MfgWorkOrderItem{},
		manufacturing.MfgIssue{},
		manufacturing.MfgIssueItem{},
		manufacturing.MfgReceipt{},

		// ===== MES：工艺主数据（阶段5a）=====
		manufacturing.MfgWorkstation{},
		manufacturing.MfgTeam{},
		manufacturing.MfgTeamMember{},
		manufacturing.MfgRouting{},
		manufacturing.MfgRoutingItem{},

		// ===== MES：派工报工（阶段5b）=====
		manufacturing.MfgWorkOrderOp{},
		manufacturing.MfgWorkReport{},
	)
	if err != nil {
		return err
	}
	return nil
}
