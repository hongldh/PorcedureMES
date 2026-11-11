package manufacturing

// ServiceGroup MES 模块 service 聚合
type ServiceGroup struct {
	BomService
	WorkOrderService
	IssueService
	ReceiptService

	// 阶段5a 主数据
	WorkstationService
	TeamService
	RoutingService

	// 阶段5b 派工报工
	DispatchService
	WorkReportService
}
