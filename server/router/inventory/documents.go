package inventory

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

// DocumentsRouter 阶段 2：业务单据 + 库存查询路由
type DocumentsRouter struct{}

func (r *DocumentsRouter) Init(Router *gin.RouterGroup) {
	g := Router.Group("inv")
	gWrite := g.Group("").Use(middleware.OperationRecord())
	gRead := g.Group("")

	// 通用注册函数 —— 每张单据都是 7 个端点
	bind := func(name string, api interface {
		Create(c *gin.Context)
		Update(c *gin.Context)
		Delete(c *gin.Context)
		Confirm(c *gin.Context)
		Void(c *gin.Context)
		FindOne(c *gin.Context)
		GetList(c *gin.Context)
	}) {
		gWrite.POST(name, api.Create)
		gWrite.PUT(name, api.Update)
		gWrite.DELETE(name, api.Delete)
		gWrite.POST(name+"/confirm", api.Confirm)
		gWrite.POST(name+"/void", api.Void)
		gRead.GET(name, api.FindOne)
		gRead.GET(name+"/list", api.GetList)
	}

	bind("purchaseOrder", &purchaseOrderApi)
	bind("purchaseReturn", &purchaseReturnApi)
	bind("salesOrder", &salesOrderApi)
	bind("salesReturn", &salesReturnApi)
	bind("transfer", &transferApi)
	bind("stocktake", &stocktakeApi)

	// ====== 库存查询 ======
	gRead.GET("stock/balance", stockQueryApi.Balance)
	gRead.GET("stock/flow", stockQueryApi.Flow)
}
