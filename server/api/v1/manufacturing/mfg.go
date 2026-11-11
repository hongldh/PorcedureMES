package manufacturing

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/manufacturing"
	mfgsvc "github.com/flipped-aurora/gin-vue-admin/server/service/manufacturing"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ============================================================
// 1. BOM
// ============================================================
type BomApi struct{}

func (a *BomApi) Create(c *gin.Context) {
	var e manufacturing.MfgBom
	if err := c.ShouldBindJSON(&e); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := bomService.Create(&e); err != nil {
		global.GVA_LOG.Error("创建BOM失败", zap.Error(err))
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(gin.H{"id": e.ID}, "创建成功", c)
}
func (a *BomApi) Update(c *gin.Context) {
	var e manufacturing.MfgBom
	if err := c.ShouldBindJSON(&e); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := bomService.Update(&e); err != nil {
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}
func (a *BomApi) Delete(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := bomService.Delete(req.Uint()); err != nil {
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}
func (a *BomApi) SetStatus(c *gin.Context) {
	var req struct {
		ID     uint `json:"id"`
		Status int8 `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := bomService.SetStatus(req.ID, req.Status); err != nil {
		response.FailWithMessage("操作失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("操作成功", c)
}
func (a *BomApi) FindOne(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := bomService.GetByID(req.Uint())
	if err != nil {
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(gin.H{"record": data}, "查询成功", c)
}
func (a *BomApi) GetList(c *gin.Context) {
	var info mfgsvc.BomSearch
	if err := c.ShouldBindQuery(&info); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := bomService.GetList(info)
	if err != nil {
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total,
		Page: info.Page, PageSize: info.PageSize}, "查询成功", c)
}
func (a *BomApi) GetActiveByProduct(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, err := bomService.GetActiveByProduct(req.Uint())
	if err != nil {
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(gin.H{"list": list}, "查询成功", c)
}

// ============================================================
// 2. 工单
// ============================================================
type WorkOrderApi struct{}

func (a *WorkOrderApi) Create(c *gin.Context) {
	var e manufacturing.MfgWorkOrder
	if err := c.ShouldBindJSON(&e); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := workOrderService.Create(&e); err != nil {
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(gin.H{"id": e.ID, "docNo": e.DocNo}, "创建成功", c)
}
func (a *WorkOrderApi) Update(c *gin.Context) {
	var e manufacturing.MfgWorkOrder
	if err := c.ShouldBindJSON(&e); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := workOrderService.Update(&e); err != nil {
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}
func (a *WorkOrderApi) Delete(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := workOrderService.Delete(req.Uint()); err != nil {
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}
func (a *WorkOrderApi) Release(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := workOrderService.Release(req.Uint(), utils.GetUserID(c)); err != nil {
		response.FailWithMessage("下达失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("下达成功", c)
}
func (a *WorkOrderApi) Close(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := workOrderService.Close(req.Uint(), utils.GetUserID(c)); err != nil {
		response.FailWithMessage("关闭失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("已关闭", c)
}
func (a *WorkOrderApi) Void(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := workOrderService.Void(req.Uint(), utils.GetUserID(c)); err != nil {
		response.FailWithMessage("作废失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("已作废", c)
}
func (a *WorkOrderApi) FindOne(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := workOrderService.GetByID(req.Uint())
	if err != nil {
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(gin.H{"record": data}, "查询成功", c)
}
func (a *WorkOrderApi) GetList(c *gin.Context) {
	var info mfgsvc.WoSearch
	if err := c.ShouldBindQuery(&info); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := workOrderService.GetList(info)
	if err != nil {
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total,
		Page: info.Page, PageSize: info.PageSize}, "查询成功", c)
}

// ============================================================
// 3. 领料单
// ============================================================
type IssueApi struct{}

func (a *IssueApi) Create(c *gin.Context) {
	var e manufacturing.MfgIssue
	if err := c.ShouldBindJSON(&e); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := issueService.Create(&e); err != nil {
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(gin.H{"id": e.ID, "docNo": e.DocNo}, "创建成功", c)
}
func (a *IssueApi) Update(c *gin.Context) {
	var e manufacturing.MfgIssue
	if err := c.ShouldBindJSON(&e); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := issueService.Update(&e); err != nil {
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}
func (a *IssueApi) Delete(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := issueService.Delete(req.Uint()); err != nil {
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}
func (a *IssueApi) Confirm(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := issueService.Confirm(req.Uint(), utils.GetUserID(c)); err != nil {
		response.FailWithMessage("确认失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("已确认", c)
}
func (a *IssueApi) Void(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := issueService.Void(req.Uint(), utils.GetUserID(c)); err != nil {
		response.FailWithMessage("作废失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("已作废", c)
}
func (a *IssueApi) FindOne(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := issueService.GetByID(req.Uint())
	if err != nil {
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(gin.H{"record": data}, "查询成功", c)
}
func (a *IssueApi) GetList(c *gin.Context) {
	var info mfgsvc.DocSearch
	if err := c.ShouldBindQuery(&info); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := issueService.GetList(info)
	if err != nil {
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total,
		Page: info.Page, PageSize: info.PageSize}, "查询成功", c)
}

// ============================================================
// 4. 完工入库单
// ============================================================
type ReceiptApi struct{}

func (a *ReceiptApi) Create(c *gin.Context) {
	var e manufacturing.MfgReceipt
	if err := c.ShouldBindJSON(&e); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := receiptService.Create(&e); err != nil {
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(gin.H{"id": e.ID, "docNo": e.DocNo}, "创建成功", c)
}
func (a *ReceiptApi) Update(c *gin.Context) {
	var e manufacturing.MfgReceipt
	if err := c.ShouldBindJSON(&e); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := receiptService.Update(&e); err != nil {
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}
func (a *ReceiptApi) Delete(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := receiptService.Delete(req.Uint()); err != nil {
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}
func (a *ReceiptApi) Confirm(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := receiptService.Confirm(req.Uint(), utils.GetUserID(c)); err != nil {
		response.FailWithMessage("确认失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("已确认", c)
}
func (a *ReceiptApi) Void(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := receiptService.Void(req.Uint(), utils.GetUserID(c)); err != nil {
		response.FailWithMessage("作废失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("已作废", c)
}
func (a *ReceiptApi) FindOne(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := receiptService.GetByID(req.Uint())
	if err != nil {
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(gin.H{"record": data}, "查询成功", c)
}
func (a *ReceiptApi) GetList(c *gin.Context) {
	var info mfgsvc.DocSearch
	if err := c.ShouldBindQuery(&info); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := receiptService.GetList(info)
	if err != nil {
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total,
		Page: info.Page, PageSize: info.PageSize}, "查询成功", c)
}
