package inventory

import (
	"strconv"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/inventory"
	invSvc "github.com/flipped-aurora/gin-vue-admin/server/service/inventory"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ============ Warehouse Api ============

type WarehouseApi struct{}

func (a *WarehouseApi) Create(c *gin.Context) {
	var e inventory.InvWarehouse
	if err := c.ShouldBindJSON(&e); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := warehouseService.Create(&e); err != nil {
		global.GVA_LOG.Error("创建仓库失败", zap.Error(err))
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}
func (a *WarehouseApi) Update(c *gin.Context) {
	var e inventory.InvWarehouse
	if err := c.ShouldBindJSON(&e); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if e.ID == 0 {
		response.FailWithMessage("ID不能为空", c)
		return
	}
	if err := warehouseService.Update(&e); err != nil {
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}
func (a *WarehouseApi) Delete(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := warehouseService.Delete(req.Uint()); err != nil {
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}
func (a *WarehouseApi) DeleteBatch(c *gin.Context) {
	var req request.IdsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	ids := make([]uint, 0, len(req.Ids))
	for _, v := range req.Ids {
		ids = append(ids, uint(v))
	}
	if err := warehouseService.DeleteBatch(ids); err != nil {
		response.FailWithMessage("批量删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}
func (a *WarehouseApi) FindOne(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := warehouseService.GetByID(req.Uint())
	if err != nil {
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(gin.H{"record": data}, "查询成功", c)
}
func (a *WarehouseApi) GetList(c *gin.Context) {
	var info request.PageInfo
	if err := c.ShouldBindQuery(&info); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := warehouseService.GetList(info)
	if err != nil {
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List: list, Total: total, Page: info.Page, PageSize: info.PageSize,
	}, "查询成功", c)
}
func (a *WarehouseApi) GetAll(c *gin.Context) {
	list, err := warehouseService.GetAll()
	if err != nil {
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(gin.H{"list": list}, "查询成功", c)
}

// ============ Location Api ============

type LocationApi struct{}

func (a *LocationApi) Create(c *gin.Context) {
	var e inventory.InvLocation
	if err := c.ShouldBindJSON(&e); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := locationService.Create(&e); err != nil {
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}
func (a *LocationApi) Update(c *gin.Context) {
	var e inventory.InvLocation
	if err := c.ShouldBindJSON(&e); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if e.ID == 0 {
		response.FailWithMessage("ID不能为空", c)
		return
	}
	if err := locationService.Update(&e); err != nil {
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}
func (a *LocationApi) Delete(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := locationService.Delete(req.Uint()); err != nil {
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}
func (a *LocationApi) DeleteBatch(c *gin.Context) {
	var req request.IdsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	ids := make([]uint, 0, len(req.Ids))
	for _, v := range req.Ids {
		ids = append(ids, uint(v))
	}
	if err := locationService.DeleteBatch(ids); err != nil {
		response.FailWithMessage("批量删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}
func (a *LocationApi) FindOne(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := locationService.GetByID(req.Uint())
	if err != nil {
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(gin.H{"record": data}, "查询成功", c)
}
func (a *LocationApi) GetList(c *gin.Context) {
	var info invSvc.LocationSearch
	if err := c.ShouldBindQuery(&info); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := locationService.GetList(info)
	if err != nil {
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List: list, Total: total, Page: info.Page, PageSize: info.PageSize,
	}, "查询成功", c)
}
func (a *LocationApi) GetByWarehouse(c *gin.Context) {
	wid, _ := strconv.ParseUint(c.Query("warehouseId"), 10, 64)
	if wid == 0 {
		response.FailWithMessage("warehouseId 必填", c)
		return
	}
	list, err := locationService.GetByWarehouse(uint(wid))
	if err != nil {
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(gin.H{"list": list}, "查询成功", c)
}
