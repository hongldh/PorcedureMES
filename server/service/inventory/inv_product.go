package inventory

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/inventory"
)

type ProductService struct{}

func (s *ProductService) Create(e *inventory.InvProduct) error {
	return global.GVA_DB.Create(e).Error
}
func (s *ProductService) Update(e *inventory.InvProduct) error {
	return global.GVA_DB.Save(e).Error
}
func (s *ProductService) Delete(id uint) error {
	return global.GVA_DB.Delete(&inventory.InvProduct{}, id).Error
}
func (s *ProductService) DeleteBatch(ids []uint) error {
	return global.GVA_DB.Delete(&[]inventory.InvProduct{}, ids).Error
}
func (s *ProductService) GetByID(id uint) (e inventory.InvProduct, err error) {
	err = global.GVA_DB.Preload("Category").Preload("Unit").First(&e, id).Error
	return
}

// ProductSearch 商品搜索条件（请求层定义在此处简化，复杂模块再拆 request 包）
type ProductSearch struct {
	request.PageInfo
	CategoryID uint   `json:"categoryId" form:"categoryId"`
	Status     int8   `json:"status" form:"status"`
	Code       string `json:"code" form:"code"`
}

func (s *ProductService) GetList(info ProductSearch) (list []inventory.InvProduct, total int64, err error) {
	db := global.GVA_DB.Model(&inventory.InvProduct{})
	if info.Keyword != "" {
		db = db.Where("code LIKE ? OR name LIKE ? OR barcode LIKE ?",
			"%"+info.Keyword+"%", "%"+info.Keyword+"%", "%"+info.Keyword+"%")
	}
	if info.CategoryID > 0 {
		db = db.Where("category_id = ?", info.CategoryID)
	}
	if info.Status > 0 {
		db = db.Where("status = ?", info.Status)
	}
	if err = db.Count(&total).Error; err != nil {
		return
	}
	err = db.Preload("Category").Preload("Unit").
		Scopes(info.PageInfo.Paginate()).Order("id desc").Find(&list).Error
	return
}

func (s *ProductService) GetAll() (list []inventory.InvProduct, err error) {
	err = global.GVA_DB.Where("status = ?", 1).Preload("Unit").Find(&list).Error
	return
}
