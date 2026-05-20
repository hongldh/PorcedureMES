package inventory

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/inventory"
)

// ============ Supplier ============

type SupplierService struct{}

func (s *SupplierService) Create(e *inventory.InvSupplier) error {
	return global.GVA_DB.Create(e).Error
}
func (s *SupplierService) Update(e *inventory.InvSupplier) error {
	return global.GVA_DB.Save(e).Error
}
func (s *SupplierService) Delete(id uint) error {
	return global.GVA_DB.Delete(&inventory.InvSupplier{}, id).Error
}
func (s *SupplierService) DeleteBatch(ids []uint) error {
	return global.GVA_DB.Delete(&[]inventory.InvSupplier{}, ids).Error
}
func (s *SupplierService) GetByID(id uint) (e inventory.InvSupplier, err error) {
	err = global.GVA_DB.First(&e, id).Error
	return
}
func (s *SupplierService) GetList(info request.PageInfo) (list []inventory.InvSupplier, total int64, err error) {
	db := global.GVA_DB.Model(&inventory.InvSupplier{})
	if info.Keyword != "" {
		db = db.Where("code LIKE ? OR name LIKE ? OR contact LIKE ?",
			"%"+info.Keyword+"%", "%"+info.Keyword+"%", "%"+info.Keyword+"%")
	}
	if err = db.Count(&total).Error; err != nil {
		return
	}
	err = db.Scopes(info.Paginate()).Order("id desc").Find(&list).Error
	return
}
func (s *SupplierService) GetAll() (list []inventory.InvSupplier, err error) {
	err = global.GVA_DB.Where("status = ?", 1).Find(&list).Error
	return
}

// ============ Customer ============

type CustomerService struct{}

func (s *CustomerService) Create(e *inventory.InvCustomer) error {
	return global.GVA_DB.Create(e).Error
}
func (s *CustomerService) Update(e *inventory.InvCustomer) error {
	return global.GVA_DB.Save(e).Error
}
func (s *CustomerService) Delete(id uint) error {
	return global.GVA_DB.Delete(&inventory.InvCustomer{}, id).Error
}
func (s *CustomerService) DeleteBatch(ids []uint) error {
	return global.GVA_DB.Delete(&[]inventory.InvCustomer{}, ids).Error
}
func (s *CustomerService) GetByID(id uint) (e inventory.InvCustomer, err error) {
	err = global.GVA_DB.First(&e, id).Error
	return
}
func (s *CustomerService) GetList(info request.PageInfo) (list []inventory.InvCustomer, total int64, err error) {
	db := global.GVA_DB.Model(&inventory.InvCustomer{})
	if info.Keyword != "" {
		db = db.Where("code LIKE ? OR name LIKE ? OR contact LIKE ?",
			"%"+info.Keyword+"%", "%"+info.Keyword+"%", "%"+info.Keyword+"%")
	}
	if err = db.Count(&total).Error; err != nil {
		return
	}
	err = db.Scopes(info.Paginate()).Order("id desc").Find(&list).Error
	return
}
func (s *CustomerService) GetAll() (list []inventory.InvCustomer, err error) {
	err = global.GVA_DB.Where("status = ?", 1).Find(&list).Error
	return
}
