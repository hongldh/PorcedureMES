package manufacturing

import api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"

type RouterGroup struct {
	MfgRouter
	MfgMasterRouter
	MfgExecutionRouter
}

var (
	bomApi       = api.ApiGroupApp.ManufacturingApiGroup.BomApi
	workOrderApi = api.ApiGroupApp.ManufacturingApiGroup.WorkOrderApi
	issueApi     = api.ApiGroupApp.ManufacturingApiGroup.IssueApi
	receiptApi   = api.ApiGroupApp.ManufacturingApiGroup.ReceiptApi

	workstationApi = api.ApiGroupApp.ManufacturingApiGroup.WorkstationApi
	teamApi        = api.ApiGroupApp.ManufacturingApiGroup.TeamApi
	routingApi     = api.ApiGroupApp.ManufacturingApiGroup.RoutingApi

	dispatchApi   = api.ApiGroupApp.ManufacturingApiGroup.DispatchApi
	workReportApi = api.ApiGroupApp.ManufacturingApiGroup.WorkReportApi
)
