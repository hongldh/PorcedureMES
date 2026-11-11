package manufacturing

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	BomApi
	WorkOrderApi
	IssueApi
	ReceiptApi

	// 阶段 5a
	WorkstationApi
	TeamApi
	RoutingApi

	// 阶段 5b
	DispatchApi
	WorkReportApi
}

var (
	bomService       = service.ServiceGroupApp.ManufacturingServiceGroup.BomService
	workOrderService = service.ServiceGroupApp.ManufacturingServiceGroup.WorkOrderService
	issueService     = service.ServiceGroupApp.ManufacturingServiceGroup.IssueService
	receiptService   = service.ServiceGroupApp.ManufacturingServiceGroup.ReceiptService

	workstationService = service.ServiceGroupApp.ManufacturingServiceGroup.WorkstationService
	teamService        = service.ServiceGroupApp.ManufacturingServiceGroup.TeamService
	routingService     = service.ServiceGroupApp.ManufacturingServiceGroup.RoutingService

	dispatchService   = service.ServiceGroupApp.ManufacturingServiceGroup.DispatchService
	workReportService = service.ServiceGroupApp.ManufacturingServiceGroup.WorkReportService
)
