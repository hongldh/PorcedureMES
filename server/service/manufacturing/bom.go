package manufacturing

import (
	"errors"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/manufacturing"
	"gorm.io/gorm"
)

// ============================================================
// BOM 物料清单服务
// ============================================================
type BomService struct{}

// Create 创建 BOM（草稿态可继续编辑；启用后建议新版本）
func (s *BomService) Create(e *manufacturing.MfgBom) error {
	if e.ProductID == 0 {
		return errors.New("必须指定成品")
	}
	if e.Version == "" {
		return errors.New("必须指定版本号")
	}
	// 校验：同一成品+版本号唯一
	var cnt int64
	global.GVA_DB.Model(&manufacturing.MfgBom{}).
		Where("product_id=? AND version=?", e.ProductID, e.Version).Count(&cnt)
	if cnt > 0 {
		return errors.New("该成品已存在相同版本号的 BOM")
	}
	return global.GVA_DB.Create(e).Error
}

// Update 更新 BOM（含子项，全量替换）
func (s *BomService) Update(e *manufacturing.MfgBom) error {
	var old manufacturing.MfgBom
	if err := global.GVA_DB.First(&old, e.ID).Error; err != nil {
		return err
	}
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// 删旧子项
		if err := tx.Where("bom_id = ?", e.ID).Delete(&manufacturing.MfgBomItem{}).Error; err != nil {
			return err
		}
		for i := range e.Items {
			e.Items[i].ID = 0
			e.Items[i].BomID = e.ID
		}
		e.CreatedAt = old.CreatedAt
		return tx.Session(&gorm.Session{FullSaveAssociations: true}).Save(e).Error
	})
}

// Delete 删除 BOM
func (s *BomService) Delete(id uint) error {
	// 校验：被工单引用过的 BOM 不允许删除
	var cnt int64
	global.GVA_DB.Model(&manufacturing.MfgWorkOrder{}).Where("bom_id=?", id).Count(&cnt)
	if cnt > 0 {
		return errors.New("该 BOM 已被工单引用，不可删除（请改为停用）")
	}
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("bom_id=?", id).Delete(&manufacturing.MfgBomItem{}).Error; err != nil {
			return err
		}
		return tx.Delete(&manufacturing.MfgBom{}, id).Error
	})
}

// SetStatus 启用/停用
func (s *BomService) SetStatus(id uint, status int8) error {
	return global.GVA_DB.Model(&manufacturing.MfgBom{}).Where("id=?", id).
		Update("status", status).Error
}

// GetByID 含子项与商品
func (s *BomService) GetByID(id uint) (e manufacturing.MfgBom, err error) {
	err = global.GVA_DB.Preload("Product").
		Preload("Items.Product").First(&e, id).Error
	return
}

// BomSearch BOM 搜索
type BomSearch struct {
	Page      int  `json:"page" form:"page"`
	PageSize  int  `json:"pageSize" form:"pageSize"`
	Status    int8 `json:"status" form:"status"`
	ProductID uint `json:"productId" form:"productId"`
}

// GetList 分页
func (s *BomService) GetList(info BomSearch) (list []manufacturing.MfgBom, total int64, err error) {
	if info.Page <= 0 {
		info.Page = 1
	}
	if info.PageSize <= 0 {
		info.PageSize = 10
	}
	db := global.GVA_DB.Model(&manufacturing.MfgBom{})
	if info.Status > 0 {
		db = db.Where("status=?", info.Status)
	}
	if info.ProductID > 0 {
		db = db.Where("product_id=?", info.ProductID)
	}
	if err = db.Count(&total).Error; err != nil {
		return
	}
	err = db.Preload("Product").
		Limit(info.PageSize).Offset((info.Page - 1) * info.PageSize).
		Order("id DESC").Find(&list).Error
	return
}

// GetActiveByProduct 取某成品的所有启用 BOM（前端工单页选择）
func (s *BomService) GetActiveByProduct(productID uint) (list []manufacturing.MfgBom, err error) {
	err = global.GVA_DB.Where("product_id=? AND status=?", productID, manufacturing.BomStatusActive).
		Preload("Items.Product").Order("id DESC").Find(&list).Error
	return
}
