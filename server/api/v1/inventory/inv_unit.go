package inventory

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/inventory"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UnitApi struct{}

func (a *UnitApi) Create(c *gin.Context) {
	var e inventory.InvUnit
	if err := c.ShouldBindJSON(&e); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := unitService.Create(&e); err != nil {
		global.GVA_LOG.Error("创建单位失败", zap.Error(err))
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

func (a *UnitApi) Update(c *gin.Context) {
	var e inventory.InvUnit
	if err := c.ShouldBindJSON(&e); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if e.ID == 0 {
		response.FailWithMessage("ID不能为空", c)
		return
	}
	if err := unitService.Update(&e); err != nil {
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

func (a *UnitApi) Delete(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := unitService.Delete(req.Uint()); err != nil {
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

func (a *UnitApi) DeleteBatch(c *gin.Context) {
	var req request.IdsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	ids := make([]uint, 0, len(req.Ids))
	for _, v := range req.Ids {
		ids = append(ids, uint(v))
	}
	if err := unitService.DeleteBatch(ids); err != nil {
		response.FailWithMessage("批量删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

func (a *UnitApi) FindOne(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := unitService.GetByID(req.Uint())
	if err != nil {
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(gin.H{"record": data}, "查询成功", c)
}

func (a *UnitApi) GetList(c *gin.Context) {
	var info request.PageInfo
	if err := c.ShouldBindQuery(&info); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := unitService.GetList(info)
	if err != nil {
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List: list, Total: total, Page: info.Page, PageSize: info.PageSize,
	}, "查询成功", c)
}

func (a *UnitApi) GetAll(c *gin.Context) {
	list, err := unitService.GetAll()
	if err != nil {
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(gin.H{"list": list}, "查询成功", c)
}
