package inventory

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/inventory"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ============ Supplier ============

type SupplierApi struct{}

func (a *SupplierApi) Create(c *gin.Context) {
	var e inventory.InvSupplier
	if err := c.ShouldBindJSON(&e); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := supplierService.Create(&e); err != nil {
		global.GVA_LOG.Error("创建供应商失败", zap.Error(err))
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}
func (a *SupplierApi) Update(c *gin.Context) {
	var e inventory.InvSupplier
	if err := c.ShouldBindJSON(&e); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if e.ID == 0 {
		response.FailWithMessage("ID不能为空", c)
		return
	}
	if err := supplierService.Update(&e); err != nil {
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}
func (a *SupplierApi) Delete(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := supplierService.Delete(req.Uint()); err != nil {
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}
func (a *SupplierApi) DeleteBatch(c *gin.Context) {
	var req request.IdsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	ids := make([]uint, 0, len(req.Ids))
	for _, v := range req.Ids {
		ids = append(ids, uint(v))
	}
	if err := supplierService.DeleteBatch(ids); err != nil {
		response.FailWithMessage("批量删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}
func (a *SupplierApi) FindOne(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := supplierService.GetByID(req.Uint())
	if err != nil {
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(gin.H{"record": data}, "查询成功", c)
}
func (a *SupplierApi) GetList(c *gin.Context) {
	var info request.PageInfo
	if err := c.ShouldBindQuery(&info); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := supplierService.GetList(info)
	if err != nil {
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List: list, Total: total, Page: info.Page, PageSize: info.PageSize,
	}, "查询成功", c)
}
func (a *SupplierApi) GetAll(c *gin.Context) {
	list, err := supplierService.GetAll()
	if err != nil {
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(gin.H{"list": list}, "查询成功", c)
}

// ============ Customer ============

type CustomerApi struct{}

func (a *CustomerApi) Create(c *gin.Context) {
	var e inventory.InvCustomer
	if err := c.ShouldBindJSON(&e); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := customerService.Create(&e); err != nil {
		global.GVA_LOG.Error("创建客户失败", zap.Error(err))
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}
func (a *CustomerApi) Update(c *gin.Context) {
	var e inventory.InvCustomer
	if err := c.ShouldBindJSON(&e); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if e.ID == 0 {
		response.FailWithMessage("ID不能为空", c)
		return
	}
	if err := customerService.Update(&e); err != nil {
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}
func (a *CustomerApi) Delete(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := customerService.Delete(req.Uint()); err != nil {
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}
func (a *CustomerApi) DeleteBatch(c *gin.Context) {
	var req request.IdsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	ids := make([]uint, 0, len(req.Ids))
	for _, v := range req.Ids {
		ids = append(ids, uint(v))
	}
	if err := customerService.DeleteBatch(ids); err != nil {
		response.FailWithMessage("批量删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}
func (a *CustomerApi) FindOne(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := customerService.GetByID(req.Uint())
	if err != nil {
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(gin.H{"record": data}, "查询成功", c)
}
func (a *CustomerApi) GetList(c *gin.Context) {
	var info request.PageInfo
	if err := c.ShouldBindQuery(&info); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := customerService.GetList(info)
	if err != nil {
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List: list, Total: total, Page: info.Page, PageSize: info.PageSize,
	}, "查询成功", c)
}
func (a *CustomerApi) GetAll(c *gin.Context) {
	list, err := customerService.GetAll()
	if err != nil {
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(gin.H{"list": list}, "查询成功", c)
}
