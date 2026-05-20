package inventory

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/inventory"
)

// ============ Warehouse ============

type WarehouseService struct{}

func (s *WarehouseService) Create(e *inventory.InvWarehouse) error {
	// 创建仓库时自动建一条默认货位 DEFAULT，方便没启用货位管理的场景
	tx := global.GVA_DB.Begin()
	if err := tx.Create(e).Error; err != nil {
		tx.Rollback()
		return err
	}
	loc := inventory.InvLocation{
		WarehouseID: e.ID,
		Code:        "DEFAULT",
		Name:        "默认货位",
		Status:      1,
	}
	if err := tx.Create(&loc).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
func (s *WarehouseService) Update(e *inventory.InvWarehouse) error {
	return global.GVA_DB.Save(e).Error
}
func (s *WarehouseService) Delete(id uint) error {
	return global.GVA_DB.Delete(&inventory.InvWarehouse{}, id).Error
}
func (s *WarehouseService) DeleteBatch(ids []uint) error {
	return global.GVA_DB.Delete(&[]inventory.InvWarehouse{}, ids).Error
}
func (s *WarehouseService) GetByID(id uint) (e inventory.InvWarehouse, err error) {
	err = global.GVA_DB.First(&e, id).Error
	return
}
func (s *WarehouseService) GetList(info request.PageInfo) (list []inventory.InvWarehouse, total int64, err error) {
	db := global.GVA_DB.Model(&inventory.InvWarehouse{})
	if info.Keyword != "" {
		db = db.Where("code LIKE ? OR name LIKE ?", "%"+info.Keyword+"%", "%"+info.Keyword+"%")
	}
	if err = db.Count(&total).Error; err != nil {
		return
	}
	err = db.Scopes(info.Paginate()).Order("id desc").Find(&list).Error
	return
}
func (s *WarehouseService) GetAll() (list []inventory.InvWarehouse, err error) {
	err = global.GVA_DB.Where("status = ?", 1).Find(&list).Error
	return
}

// ============ Location ============

type LocationService struct{}

func (s *LocationService) Create(e *inventory.InvLocation) error {
	return global.GVA_DB.Create(e).Error
}
func (s *LocationService) Update(e *inventory.InvLocation) error {
	return global.GVA_DB.Save(e).Error
}
func (s *LocationService) Delete(id uint) error {
	return global.GVA_DB.Delete(&inventory.InvLocation{}, id).Error
}
func (s *LocationService) DeleteBatch(ids []uint) error {
	return global.GVA_DB.Delete(&[]inventory.InvLocation{}, ids).Error
}
func (s *LocationService) GetByID(id uint) (e inventory.InvLocation, err error) {
	err = global.GVA_DB.Preload("Warehouse").First(&e, id).Error
	return
}

type LocationSearch struct {
	request.PageInfo
	WarehouseID uint `json:"warehouseId" form:"warehouseId"`
}

func (s *LocationService) GetList(info LocationSearch) (list []inventory.InvLocation, total int64, err error) {
	db := global.GVA_DB.Model(&inventory.InvLocation{})
	if info.Keyword != "" {
		db = db.Where("code LIKE ? OR name LIKE ?", "%"+info.Keyword+"%", "%"+info.Keyword+"%")
	}
	if info.WarehouseID > 0 {
		db = db.Where("warehouse_id = ?", info.WarehouseID)
	}
	if err = db.Count(&total).Error; err != nil {
		return
	}
	err = db.Preload("Warehouse").Scopes(info.PageInfo.Paginate()).Order("sort asc, id desc").Find(&list).Error
	return
}

// GetByWarehouse 获取某仓库下所有货位（下拉用）
func (s *LocationService) GetByWarehouse(wid uint) (list []inventory.InvLocation, err error) {
	err = global.GVA_DB.Where("warehouse_id = ? AND status = ?", wid, 1).Order("sort asc").Find(&list).Error
	return
}
