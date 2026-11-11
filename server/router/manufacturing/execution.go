package manufacturing

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

// MfgExecutionRouter 阶段 5b：工序 / 派工 / 报工
type MfgExecutionRouter struct{}

func (r *MfgExecutionRouter) Init(Router *gin.RouterGroup) {
	g := Router.Group("mfg")
	gWrite := g.Group("").Use(middleware.OperationRecord())
	gRead := g.Group("")

	// 工单工序 + 派工
	gRead.GET("workOrder/ops", dispatchApi.GetOps)
	gWrite.POST("op/dispatch", dispatchApi.Dispatch)

	// 报工
	gWrite.POST("report", workReportApi.Create)
	gWrite.POST("report/void", workReportApi.Void)
	gRead.GET("report/list", workReportApi.GetList)
	gRead.GET("report/myOps", workReportApi.MyOps)
}
