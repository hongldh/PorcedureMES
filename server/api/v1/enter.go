package v1

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/example"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/inventory"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/manufacturing"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/system"
)

var ApiGroupApp = new(ApiGroup)

type ApiGroup struct {
	SystemApiGroup        system.ApiGroup
	ExampleApiGroup       example.ApiGroup
	InventoryApiGroup     inventory.ApiGroup
	ManufacturingApiGroup manufacturing.ApiGroup
}
