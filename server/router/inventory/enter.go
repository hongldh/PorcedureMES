package inventory

import api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"

// RouterGroup 进销存模块路由聚合
type RouterGroup struct {
	MasterDataRouter
	DocumentsRouter
}

// 包级单例：api 引用
var (
	categoryApi  = api.ApiGroupApp.InventoryApiGroup.CategoryApi
	unitApi      = api.ApiGroupApp.InventoryApiGroup.UnitApi
	productApi   = api.ApiGroupApp.InventoryApiGroup.ProductApi
	warehouseApi = api.ApiGroupApp.InventoryApiGroup.WarehouseApi
	locationApi  = api.ApiGroupApp.InventoryApiGroup.LocationApi
	supplierApi  = api.ApiGroupApp.InventoryApiGroup.SupplierApi
	customerApi  = api.ApiGroupApp.InventoryApiGroup.CustomerApi

	purchaseOrderApi  = api.ApiGroupApp.InventoryApiGroup.PurchaseOrderApi
	purchaseReturnApi = api.ApiGroupApp.InventoryApiGroup.PurchaseReturnApi
	salesOrderApi     = api.ApiGroupApp.InventoryApiGroup.SalesOrderApi
	salesReturnApi    = api.ApiGroupApp.InventoryApiGroup.SalesReturnApi
	transferApi       = api.ApiGroupApp.InventoryApiGroup.TransferApi
	stocktakeApi      = api.ApiGroupApp.InventoryApiGroup.StocktakeApi
	stockQueryApi     = api.ApiGroupApp.InventoryApiGroup.StockQueryApi
)
