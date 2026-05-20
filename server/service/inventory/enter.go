package inventory

// ServiceGroup 进销存模块所有 service 的聚合
//
// === Go 学习要点 ===
// 这种"嵌入聚合"的写法是 GVA 的统一惯例：
//   - 把所有子 service 用结构体嵌入，外部通过 ServiceGroup.CategoryService 直接调用
//   - 嵌入的字段没有名字，类型名即字段名（CategoryService 嵌入后用 .CategoryService 访问）
type ServiceGroup struct {
	CategoryService
	UnitService
	ProductService
	WarehouseService
	LocationService
	SupplierService
	CustomerService

	// === 阶段 2：单据 + 库存引擎 ===
	SequenceService
	StockEngine
	PurchaseOrderService
	PurchaseReturnService
	SalesOrderService
	SalesReturnService
	TransferService
	StocktakeService
	StockQueryService
}
