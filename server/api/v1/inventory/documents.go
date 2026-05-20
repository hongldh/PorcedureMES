package inventory

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	invsvc "github.com/flipped-aurora/gin-vue-admin/server/service/inventory"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/inventory"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// === Go 学习要点 ===
// 同一类业务（单据）共享一组 handler 模板：
//   Create / Update / Delete / Confirm / Void / FindOne / GetList
// 这里把 6 张单据的 API 写在一个文件里，便于对照。

// ============================================================
// 1. 采购入库
// ============================================================
type PurchaseOrderApi struct{}

func (a *PurchaseOrderApi) Create(c *gin.Context) {
	var e inventory.InvPurchaseOrder
	if err := c.ShouldBindJSON(&e); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := purchaseOrderService.Create(&e); err != nil {
		global.GVA_LOG.Error("创建采购入库单失败", zap.Error(err))
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(gin.H{"id": e.ID, "docNo": e.DocNo}, "创建成功", c)
}
func (a *PurchaseOrderApi) Update(c *gin.Context) {
	var e inventory.InvPurchaseOrder
	if err := c.ShouldBindJSON(&e); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := purchaseOrderService.Update(&e); err != nil {
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}
func (a *PurchaseOrderApi) Delete(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := purchaseOrderService.Delete(req.Uint()); err != nil {
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}
func (a *PurchaseOrderApi) Confirm(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := purchaseOrderService.Confirm(req.Uint(), utils.GetUserID(c)); err != nil {
		response.FailWithMessage("执行失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("执行成功", c)
}
func (a *PurchaseOrderApi) Void(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := purchaseOrderService.Void(req.Uint(), utils.GetUserID(c)); err != nil {
		response.FailWithMessage("作废失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("作废成功", c)
}
func (a *PurchaseOrderApi) FindOne(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := purchaseOrderService.GetByID(req.Uint())
	if err != nil {
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(gin.H{"record": data}, "查询成功", c)
}
func (a *PurchaseOrderApi) GetList(c *gin.Context) {
	var info invsvc.DocSearch
	if err := c.ShouldBindQuery(&info); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := purchaseOrderService.GetList(info)
	if err != nil {
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total, Page: info.Page, PageSize: info.PageSize}, "查询成功", c)
}

// ============================================================
// 2. 采购退货
// ============================================================
type PurchaseReturnApi struct{}

func (a *PurchaseReturnApi) Create(c *gin.Context) {
	var e inventory.InvPurchaseReturn
	if err := c.ShouldBindJSON(&e); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := purchaseReturnService.Create(&e); err != nil {
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(gin.H{"id": e.ID, "docNo": e.DocNo}, "创建成功", c)
}
func (a *PurchaseReturnApi) Update(c *gin.Context) {
	var e inventory.InvPurchaseReturn
	if err := c.ShouldBindJSON(&e); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := purchaseReturnService.Update(&e); err != nil {
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}
func (a *PurchaseReturnApi) Delete(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := purchaseReturnService.Delete(req.Uint()); err != nil {
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}
func (a *PurchaseReturnApi) Confirm(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := purchaseReturnService.Confirm(req.Uint(), utils.GetUserID(c)); err != nil {
		response.FailWithMessage("执行失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("执行成功", c)
}
func (a *PurchaseReturnApi) Void(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := purchaseReturnService.Void(req.Uint(), utils.GetUserID(c)); err != nil {
		response.FailWithMessage("作废失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("作废成功", c)
}
func (a *PurchaseReturnApi) FindOne(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := purchaseReturnService.GetByID(req.Uint())
	if err != nil {
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(gin.H{"record": data}, "查询成功", c)
}
func (a *PurchaseReturnApi) GetList(c *gin.Context) {
	var info invsvc.DocSearch
	if err := c.ShouldBindQuery(&info); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := purchaseReturnService.GetList(info)
	if err != nil {
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total, Page: info.Page, PageSize: info.PageSize}, "查询成功", c)
}

// ============================================================
// 3. 销售出库
// ============================================================
type SalesOrderApi struct{}

func (a *SalesOrderApi) Create(c *gin.Context) {
	var e inventory.InvSalesOrder
	if err := c.ShouldBindJSON(&e); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := salesOrderService.Create(&e); err != nil {
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(gin.H{"id": e.ID, "docNo": e.DocNo}, "创建成功", c)
}
func (a *SalesOrderApi) Update(c *gin.Context) {
	var e inventory.InvSalesOrder
	if err := c.ShouldBindJSON(&e); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := salesOrderService.Update(&e); err != nil {
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}
func (a *SalesOrderApi) Delete(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := salesOrderService.Delete(req.Uint()); err != nil {
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}
func (a *SalesOrderApi) Confirm(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := salesOrderService.Confirm(req.Uint(), utils.GetUserID(c)); err != nil {
		response.FailWithMessage("执行失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("执行成功", c)
}
func (a *SalesOrderApi) Void(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := salesOrderService.Void(req.Uint(), utils.GetUserID(c)); err != nil {
		response.FailWithMessage("作废失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("作废成功", c)
}
func (a *SalesOrderApi) FindOne(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := salesOrderService.GetByID(req.Uint())
	if err != nil {
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(gin.H{"record": data}, "查询成功", c)
}
func (a *SalesOrderApi) GetList(c *gin.Context) {
	var info invsvc.DocSearch
	if err := c.ShouldBindQuery(&info); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := salesOrderService.GetList(info)
	if err != nil {
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total, Page: info.Page, PageSize: info.PageSize}, "查询成功", c)
}

// ============================================================
// 4. 销售退货
// ============================================================
type SalesReturnApi struct{}

func (a *SalesReturnApi) Create(c *gin.Context) {
	var e inventory.InvSalesReturn
	if err := c.ShouldBindJSON(&e); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := salesReturnService.Create(&e); err != nil {
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(gin.H{"id": e.ID, "docNo": e.DocNo}, "创建成功", c)
}
func (a *SalesReturnApi) Update(c *gin.Context) {
	var e inventory.InvSalesReturn
	if err := c.ShouldBindJSON(&e); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := salesReturnService.Update(&e); err != nil {
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}
func (a *SalesReturnApi) Delete(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := salesReturnService.Delete(req.Uint()); err != nil {
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}
func (a *SalesReturnApi) Confirm(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := salesReturnService.Confirm(req.Uint(), utils.GetUserID(c)); err != nil {
		response.FailWithMessage("执行失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("执行成功", c)
}
func (a *SalesReturnApi) Void(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := salesReturnService.Void(req.Uint(), utils.GetUserID(c)); err != nil {
		response.FailWithMessage("作废失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("作废成功", c)
}
func (a *SalesReturnApi) FindOne(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := salesReturnService.GetByID(req.Uint())
	if err != nil {
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(gin.H{"record": data}, "查询成功", c)
}
func (a *SalesReturnApi) GetList(c *gin.Context) {
	var info invsvc.DocSearch
	if err := c.ShouldBindQuery(&info); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := salesReturnService.GetList(info)
	if err != nil {
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total, Page: info.Page, PageSize: info.PageSize}, "查询成功", c)
}

// ============================================================
// 5. 调拨
// ============================================================
type TransferApi struct{}

func (a *TransferApi) Create(c *gin.Context) {
	var e inventory.InvTransfer
	if err := c.ShouldBindJSON(&e); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := transferService.Create(&e); err != nil {
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(gin.H{"id": e.ID, "docNo": e.DocNo}, "创建成功", c)
}
func (a *TransferApi) Update(c *gin.Context) {
	var e inventory.InvTransfer
	if err := c.ShouldBindJSON(&e); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := transferService.Update(&e); err != nil {
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}
func (a *TransferApi) Delete(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := transferService.Delete(req.Uint()); err != nil {
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}
func (a *TransferApi) Confirm(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := transferService.Confirm(req.Uint(), utils.GetUserID(c)); err != nil {
		response.FailWithMessage("执行失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("执行成功", c)
}
func (a *TransferApi) Void(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := transferService.Void(req.Uint(), utils.GetUserID(c)); err != nil {
		response.FailWithMessage("作废失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("作废成功", c)
}
func (a *TransferApi) FindOne(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := transferService.GetByID(req.Uint())
	if err != nil {
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(gin.H{"record": data}, "查询成功", c)
}
func (a *TransferApi) GetList(c *gin.Context) {
	var info invsvc.DocSearch
	if err := c.ShouldBindQuery(&info); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := transferService.GetList(info)
	if err != nil {
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total, Page: info.Page, PageSize: info.PageSize}, "查询成功", c)
}

// ============================================================
// 6. 盘点
// ============================================================
type StocktakeApi struct{}

func (a *StocktakeApi) Create(c *gin.Context) {
	var e inventory.InvStocktake
	if err := c.ShouldBindJSON(&e); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := stocktakeService.Create(&e); err != nil {
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(gin.H{"id": e.ID, "docNo": e.DocNo}, "创建成功", c)
}
func (a *StocktakeApi) Update(c *gin.Context) {
	var e inventory.InvStocktake
	if err := c.ShouldBindJSON(&e); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := stocktakeService.Update(&e); err != nil {
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}
func (a *StocktakeApi) Delete(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := stocktakeService.Delete(req.Uint()); err != nil {
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}
func (a *StocktakeApi) Confirm(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := stocktakeService.Confirm(req.Uint(), utils.GetUserID(c)); err != nil {
		response.FailWithMessage("执行失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("执行成功", c)
}
func (a *StocktakeApi) Void(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := stocktakeService.Void(req.Uint(), utils.GetUserID(c)); err != nil {
		response.FailWithMessage("作废失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("作废成功", c)
}
func (a *StocktakeApi) FindOne(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := stocktakeService.GetByID(req.Uint())
	if err != nil {
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(gin.H{"record": data}, "查询成功", c)
}
func (a *StocktakeApi) GetList(c *gin.Context) {
	var info invsvc.DocSearch
	if err := c.ShouldBindQuery(&info); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := stocktakeService.GetList(info)
	if err != nil {
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total, Page: info.Page, PageSize: info.PageSize}, "查询成功", c)
}

// ============================================================
// 7. 库存查询（余额 + 流水）
// ============================================================
type StockQueryApi struct{}

func (a *StockQueryApi) Balance(c *gin.Context) {
	var info invsvc.StockSearch
	if err := c.ShouldBindQuery(&info); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := stockQueryService.BalanceList(info)
	if err != nil {
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total, Page: info.Page, PageSize: info.PageSize}, "查询成功", c)
}

func (a *StockQueryApi) Flow(c *gin.Context) {
	var info invsvc.FlowSearch
	if err := c.ShouldBindQuery(&info); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := stockQueryService.FlowList(info)
	if err != nil {
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total, Page: info.Page, PageSize: info.PageSize}, "查询成功", c)
}
