// 阶段 5b：工单工序 / 派工 / 报工 API
package manufacturing

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/manufacturing"
	mfgsvc "github.com/flipped-aurora/gin-vue-admin/server/service/manufacturing"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
)

// ============================================================
// 1. 工单工序 + 派工
// ============================================================
type DispatchApi struct{}

// GetOps 工单工序列表  GET /mfg/workOrder/ops?id=
func (a *DispatchApi) GetOps(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, err := workOrderService.GetOps(req.Uint())
	if err != nil {
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(gin.H{"list": list}, "查询成功", c)
}

// Dispatch 派工 / 改派  POST /mfg/op/dispatch
func (a *DispatchApi) Dispatch(c *gin.Context) {
	var p mfgsvc.DispatchParam
	if err := c.ShouldBindJSON(&p); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := dispatchService.Dispatch(p, utils.GetUserID(c)); err != nil {
		response.FailWithMessage("派工失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("派工成功", c)
}

// ============================================================
// 2. 报工
// ============================================================
type WorkReportApi struct{}

// Create 报工（即时生效）  POST /mfg/report
func (a *WorkReportApi) Create(c *gin.Context) {
	var e manufacturing.MfgWorkReport
	if err := c.ShouldBindJSON(&e); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := workReportService.Create(&e, utils.GetUserID(c)); err != nil {
		response.FailWithMessage("报工失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(gin.H{"id": e.ID, "docNo": e.DocNo}, "报工成功", c)
}

// Void 作废报工  POST /mfg/report/void
func (a *WorkReportApi) Void(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := workReportService.Void(req.Uint(), utils.GetUserID(c)); err != nil {
		response.FailWithMessage("作废失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("已作废", c)
}

// GetList 报工记录分页  GET /mfg/report/list
func (a *WorkReportApi) GetList(c *gin.Context) {
	var info mfgsvc.ReportSearch
	if err := c.ShouldBindQuery(&info); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := workReportService.GetList(info)
	if err != nil {
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total,
		Page: info.Page, PageSize: info.PageSize}, "查询成功", c)
}

// MyOps 当前登录人可报工的工序  GET /mfg/report/myOps
func (a *WorkReportApi) MyOps(c *gin.Context) {
	list, err := workReportService.MyOps(utils.GetUserID(c))
	if err != nil {
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(gin.H{"list": list}, "查询成功", c)
}
