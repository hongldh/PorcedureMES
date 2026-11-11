package manufacturing

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

// MfgMasterRouter 阶段 5a：工位 / 班组 / 工艺路线
type MfgMasterRouter struct{}

func (r *MfgMasterRouter) Init(Router *gin.RouterGroup) {
	g := Router.Group("mfg")
	gWrite := g.Group("").Use(middleware.OperationRecord())
	gRead := g.Group("")

	// 工位
	gWrite.POST("workstation", workstationApi.Create)
	gWrite.PUT("workstation", workstationApi.Update)
	gWrite.DELETE("workstation", workstationApi.Delete)
	gRead.GET("workstation", workstationApi.FindOne)
	gRead.GET("workstation/list", workstationApi.GetList)
	gRead.GET("workstation/all", workstationApi.GetAll)

	// 班组
	gWrite.POST("team", teamApi.Create)
	gWrite.PUT("team", teamApi.Update)
	gWrite.DELETE("team", teamApi.Delete)
	gRead.GET("team", teamApi.FindOne)
	gRead.GET("team/list", teamApi.GetList)

	// 工艺路线
	gWrite.POST("routing", routingApi.Create)
	gWrite.PUT("routing", routingApi.Update)
	gWrite.DELETE("routing", routingApi.Delete)
	gWrite.POST("routing/status", routingApi.SetStatus)
	gRead.GET("routing", routingApi.FindOne)
	gRead.GET("routing/list", routingApi.GetList)
	gRead.GET("routing/activeByProduct", routingApi.GetActiveByProduct)
}
