package inventory

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/inventory"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CategoryApi 商品分类 Controller
//
// === Go 学习要点：Gin Controller 通用套路 ===
// 每个 handler 都是 func(c *gin.Context) 签名
// 步骤通常是：
//   1. 用 ShouldBindJSON / ShouldBindQuery 把请求绑定到结构体
//   2. 校验
//   3. 调 service
//   4. 用 response.OkXxx / FailXxx 统一返回
type CategoryApi struct{}

// Create 创建分类
func (a *CategoryApi) Create(c *gin.Context) {
	var e inventory.InvCategory
	if err := c.ShouldBindJSON(&e); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := categoryService.Create(&e); err != nil {
		global.GVA_LOG.Error("创建分类失败", zap.Error(err))
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// Update 更新
func (a *CategoryApi) Update(c *gin.Context) {
	var e inventory.InvCategory
	if err := c.ShouldBindJSON(&e); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if e.ID == 0 {
		response.FailWithMessage("ID不能为空", c)
		return
	}
	if err := categoryService.Update(&e); err != nil {
		global.GVA_LOG.Error("更新失败", zap.Error(err))
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// Delete 删除
func (a *CategoryApi) Delete(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := categoryService.Delete(req.Uint()); err != nil {
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteBatch 批量删除
func (a *CategoryApi) DeleteBatch(c *gin.Context) {
	var req request.IdsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	ids := make([]uint, 0, len(req.Ids))
	for _, v := range req.Ids {
		ids = append(ids, uint(v))
	}
	if err := categoryService.DeleteBatch(ids); err != nil {
		response.FailWithMessage("批量删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// FindOne 查询单条
func (a *CategoryApi) FindOne(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := categoryService.GetByID(req.Uint())
	if err != nil {
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(gin.H{"record": data}, "查询成功", c)
}

// GetList 分页查询
func (a *CategoryApi) GetList(c *gin.Context) {
	var info request.PageInfo
	if err := c.ShouldBindQuery(&info); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := categoryService.GetList(info)
	if err != nil {
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List: list, Total: total, Page: info.Page, PageSize: info.PageSize,
	}, "查询成功", c)
}

// GetAll 不分页（下拉用）
func (a *CategoryApi) GetAll(c *gin.Context) {
	list, err := categoryService.GetAll()
	if err != nil {
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(gin.H{"list": list}, "查询成功", c)
}
