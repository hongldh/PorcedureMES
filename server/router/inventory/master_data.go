package inventory

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

// MasterDataRouter 主数据相关路由（分类/单位/商品/仓库/货位/供应商/客户）
type MasterDataRouter struct{}

// Init 注册路由到 PrivateGroup（需 JWT + Casbin 鉴权）
//
// === Go 学习要点：路由分组与中间件 ===
// 1. Router.Group("xx") 创建子分组
// 2. .Use(middleware.OperationRecord()) 给写操作链路加操作日志中间件
// 3. 写操作（POST/PUT/DELETE）走有日志的组；查询走没日志的组（避免日志爆炸）
func (r *MasterDataRouter) Init(Router *gin.RouterGroup) {
	g := Router.Group("inv")
	gWrite := g.Group("").Use(middleware.OperationRecord())
	gRead := g.Group("")

	// ====== 分类 ======
	{
		gWrite.POST("category", categoryApi.Create)
		gWrite.PUT("category", categoryApi.Update)
		gWrite.DELETE("category", categoryApi.Delete)
		gWrite.DELETE("category/batch", categoryApi.DeleteBatch)
		gRead.GET("category", categoryApi.FindOne)
		gRead.GET("category/list", categoryApi.GetList)
		gRead.GET("category/all", categoryApi.GetAll)
	}
	// ====== 单位 ======
	{
		gWrite.POST("unit", unitApi.Create)
		gWrite.PUT("unit", unitApi.Update)
		gWrite.DELETE("unit", unitApi.Delete)
		gWrite.DELETE("unit/batch", unitApi.DeleteBatch)
		gRead.GET("unit", unitApi.FindOne)
		gRead.GET("unit/list", unitApi.GetList)
		gRead.GET("unit/all", unitApi.GetAll)
	}
	// ====== 商品 ======
	{
		gWrite.POST("product", productApi.Create)
		gWrite.PUT("product", productApi.Update)
		gWrite.DELETE("product", productApi.Delete)
		gWrite.DELETE("product/batch", productApi.DeleteBatch)
		gRead.GET("product", productApi.FindOne)
		gRead.GET("product/list", productApi.GetList)
		gRead.GET("product/all", productApi.GetAll)
	}
	// ====== 仓库 ======
	{
		gWrite.POST("warehouse", warehouseApi.Create)
		gWrite.PUT("warehouse", warehouseApi.Update)
		gWrite.DELETE("warehouse", warehouseApi.Delete)
		gWrite.DELETE("warehouse/batch", warehouseApi.DeleteBatch)
		gRead.GET("warehouse", warehouseApi.FindOne)
		gRead.GET("warehouse/list", warehouseApi.GetList)
		gRead.GET("warehouse/all", warehouseApi.GetAll)
	}
	// ====== 货位 ======
	{
		gWrite.POST("location", locationApi.Create)
		gWrite.PUT("location", locationApi.Update)
		gWrite.DELETE("location", locationApi.Delete)
		gWrite.DELETE("location/batch", locationApi.DeleteBatch)
		gRead.GET("location", locationApi.FindOne)
		gRead.GET("location/list", locationApi.GetList)
		gRead.GET("location/byWarehouse", locationApi.GetByWarehouse)
	}
	// ====== 供应商 ======
	{
		gWrite.POST("supplier", supplierApi.Create)
		gWrite.PUT("supplier", supplierApi.Update)
		gWrite.DELETE("supplier", supplierApi.Delete)
		gWrite.DELETE("supplier/batch", supplierApi.DeleteBatch)
		gRead.GET("supplier", supplierApi.FindOne)
		gRead.GET("supplier/list", supplierApi.GetList)
		gRead.GET("supplier/all", supplierApi.GetAll)
	}
	// ====== 客户 ======
	{
		gWrite.POST("customer", customerApi.Create)
		gWrite.PUT("customer", customerApi.Update)
		gWrite.DELETE("customer", customerApi.Delete)
		gWrite.DELETE("customer/batch", customerApi.DeleteBatch)
		gRead.GET("customer", customerApi.FindOne)
		gRead.GET("customer/list", customerApi.GetList)
		gRead.GET("customer/all", customerApi.GetAll)
	}
}
