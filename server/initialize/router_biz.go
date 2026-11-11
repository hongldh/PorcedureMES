package initialize

import (
	"github.com/flipped-aurora/gin-vue-admin/server/router"
	"github.com/gin-gonic/gin"
)

// 占位方法，保证文件可以正确加载，避免go空变量检测报错，请勿删除。
func holder(routers ...*gin.RouterGroup) {
	_ = routers
	_ = router.RouterGroupApp
}

func initBizRouter(routers ...*gin.RouterGroup) {
	privateGroup := routers[0]
	publicGroup := routers[1]

	// ===== 进销存模块路由 =====
	inventoryRouter := router.RouterGroupApp.Inventory
	inventoryRouter.MasterDataRouter.Init(privateGroup)
	inventoryRouter.DocumentsRouter.Init(privateGroup)

	// ===== MES 模块路由（阶段4 + 5a）=====
	mfgRouter := router.RouterGroupApp.Manufacturing
	mfgRouter.MfgRouter.Init(privateGroup)
		mfgRouter.MfgMasterRouter.Init(privateGroup)
		mfgRouter.MfgExecutionRouter.Init(privateGroup) // 5b 派工报工

	holder(publicGroup, privateGroup)
}
