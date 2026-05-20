package inventory

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/inventory"
)

type UnitService struct{}

func (s *UnitService) Create(e *inventory.InvUnit) error { return global.GVA_DB.Create(e).Error }
func (s *UnitService) Update(e *inventory.InvUnit) error { return global.GVA_DB.Save(e).Error }
func (s *UnitService) Delete(id uint) error {
	return global.GVA_DB.Delete(&inventory.InvUnit{}, id).Error
}
func (s *UnitService) DeleteBatch(ids []uint) error {
	return global.GVA_DB.Delete(&[]inventory.InvUnit{}, ids).Error
}
func (s *UnitService) GetByID(id uint) (e inventory.InvUnit, err error) {
	err = global.GVA_DB.First(&e, id).Error
	return
}
func (s *UnitService) GetList(info request.PageInfo) (list []inventory.InvUnit, total int64, err error) {
	db := global.GVA_DB.Model(&inventory.InvUnit{})
	if info.Keyword != "" {
		db = db.Where("code LIKE ? OR name LIKE ?", "%"+info.Keyword+"%", "%"+info.Keyword+"%")
	}
	if err = db.Count(&total).Error; err != nil {
		return
	}
	err = db.Scopes(info.Paginate()).Order("id desc").Find(&list).Error
	return
}
func (s *UnitService) GetAll() (list []inventory.InvUnit, err error) {
	err = global.GVA_DB.Where("status = ?", 1).Order("id desc").Find(&list).Error
	return
}
