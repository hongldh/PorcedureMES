package manufacturing

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

// MfgRouter 阶段 4：MES 基础（BOM/工单/领料/完工入库）
type MfgRouter struct{}

func (r *MfgRouter) Init(Router *gin.RouterGroup) {
	g := Router.Group("mfg")
	gWrite := g.Group("").Use(middleware.OperationRecord())
	gRead := g.Group("")

	// === BOM ===
	gWrite.POST("bom", bomApi.Create)
	gWrite.PUT("bom", bomApi.Update)
	gWrite.DELETE("bom", bomApi.Delete)
	gWrite.POST("bom/status", bomApi.SetStatus)
	gRead.GET("bom", bomApi.FindOne)
	gRead.GET("bom/list", bomApi.GetList)
	gRead.GET("bom/activeByProduct", bomApi.GetActiveByProduct)

	// === 工单 ===
	gWrite.POST("workOrder", workOrderApi.Create)
	gWrite.PUT("workOrder", workOrderApi.Update)
	gWrite.DELETE("workOrder", workOrderApi.Delete)
	gWrite.POST("workOrder/release", workOrderApi.Release)
	gWrite.POST("workOrder/close", workOrderApi.Close)
	gWrite.POST("workOrder/void", workOrderApi.Void)
	gRead.GET("workOrder", workOrderApi.FindOne)
	gRead.GET("workOrder/list", workOrderApi.GetList)

	// === 领料 ===
	gWrite.POST("issue", issueApi.Create)
	gWrite.PUT("issue", issueApi.Update)
	gWrite.DELETE("issue", issueApi.Delete)
	gWrite.POST("issue/confirm", issueApi.Confirm)
	gWrite.POST("issue/void", issueApi.Void)
	gRead.GET("issue", issueApi.FindOne)
	gRead.GET("issue/list", issueApi.GetList)

	// === 完工入库 ===
	gWrite.POST("receipt", receiptApi.Create)
	gWrite.PUT("receipt", receiptApi.Update)
	gWrite.DELETE("receipt", receiptApi.Delete)
	gWrite.POST("receipt/confirm", receiptApi.Confirm)
	gWrite.POST("receipt/void", receiptApi.Void)
	gRead.GET("receipt", receiptApi.FindOne)
	gRead.GET("receipt/list", receiptApi.GetList)
}
