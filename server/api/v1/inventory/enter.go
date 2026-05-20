package inventory

import "github.com/flipped-aurora/gin-vue-admin/server/service"

// ApiGroup 进销存模块所有 Api 聚合
type ApiGroup struct {
	CategoryApi
	UnitApi
	ProductApi
	WarehouseApi
	LocationApi
	SupplierApi
	CustomerApi

	// 阶段 2
	PurchaseOrderApi
	PurchaseReturnApi
	SalesOrderApi
	SalesReturnApi
	TransferApi
	StocktakeApi
	StockQueryApi
}

// 包级单例：service 引用，供本包内所有 Api 使用
var (
	categoryService  = service.ServiceGroupApp.InventoryServiceGroup.CategoryService
	unitService      = service.ServiceGroupApp.InventoryServiceGroup.UnitService
	productService   = service.ServiceGroupApp.InventoryServiceGroup.ProductService
	warehouseService = service.ServiceGroupApp.InventoryServiceGroup.WarehouseService
	locationService  = service.ServiceGroupApp.InventoryServiceGroup.LocationService
	supplierService  = service.ServiceGroupApp.InventoryServiceGroup.SupplierService
	customerService  = service.ServiceGroupApp.InventoryServiceGroup.CustomerService

	purchaseOrderService  = service.ServiceGroupApp.InventoryServiceGroup.PurchaseOrderService
	purchaseReturnService = service.ServiceGroupApp.InventoryServiceGroup.PurchaseReturnService
	salesOrderService     = service.ServiceGroupApp.InventoryServiceGroup.SalesOrderService
	salesReturnService    = service.ServiceGroupApp.InventoryServiceGroup.SalesReturnService
	transferService       = service.ServiceGroupApp.InventoryServiceGroup.TransferService
	stocktakeService      = service.ServiceGroupApp.InventoryServiceGroup.StocktakeService
	stockQueryService     = service.ServiceGroupApp.InventoryServiceGroup.StockQueryService
)
