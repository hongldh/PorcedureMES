// 阶段 5a：工位 + 班组 + 工艺路线 service
package manufacturing

import (
	"errors"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/manufacturing"
	"gorm.io/gorm"
)

// ============================================================
// 1. 工位
// ============================================================
type WorkstationService struct{}

func (s *WorkstationService) Create(e *manufacturing.MfgWorkstation) error {
	if e.Code == "" || e.Name == "" {
		return errors.New("编码与名称必填")
	}
	return global.GVA_DB.Create(e).Error
}
func (s *WorkstationService) Update(e *manufacturing.MfgWorkstation) error {
	var old manufacturing.MfgWorkstation
	if err := global.GVA_DB.First(&old, e.ID).Error; err != nil {
		return err
	}
	e.CreatedAt = old.CreatedAt
	return global.GVA_DB.Save(e).Error
}
func (s *WorkstationService) Delete(id uint) error {
	// 被工艺路线引用的不可删
	var cnt int64
	global.GVA_DB.Model(&manufacturing.MfgRoutingItem{}).Where("workstation_id=?", id).Count(&cnt)
	if cnt > 0 {
		return errors.New("该工位已被工艺路线引用，不可删除（请改为停用）")
	}
	return global.GVA_DB.Delete(&manufacturing.MfgWorkstation{}, id).Error
}
func (s *WorkstationService) GetByID(id uint) (e manufacturing.MfgWorkstation, err error) {
	err = global.GVA_DB.Preload("Warehouse").First(&e, id).Error
	return
}

type WorkstationSearch struct {
	request.PageInfo
	Keyword string `json:"keyword" form:"keyword"`
	Type    int8   `json:"type" form:"type"`
	Status  int8   `json:"status" form:"status"`
}

func (s *WorkstationService) GetList(info WorkstationSearch) (list []manufacturing.MfgWorkstation, total int64, err error) {
	if info.Page <= 0 {
		info.Page = 1
	}
	if info.PageSize <= 0 {
		info.PageSize = 10
	}
	db := global.GVA_DB.Model(&manufacturing.MfgWorkstation{})
	if info.Keyword != "" {
		db = db.Where("code LIKE ? OR name LIKE ?", "%"+info.Keyword+"%", "%"+info.Keyword+"%")
	}
	if info.Type > 0 {
		db = db.Where("type=?", info.Type)
	}
	if info.Status > 0 {
		db = db.Where("status=?", info.Status)
	}
	if err = db.Count(&total).Error; err != nil {
		return
	}
	err = db.Preload("Warehouse").
		Limit(info.PageSize).Offset((info.Page - 1) * info.PageSize).
		Order("id DESC").Find(&list).Error
	return
}

func (s *WorkstationService) GetAll() (list []manufacturing.MfgWorkstation, err error) {
	err = global.GVA_DB.Where("status=?", manufacturing.BasicStatusEnabled).Order("code").Find(&list).Error
	return
}

// ============================================================
// 2. 班组
// ============================================================
type TeamService struct{}

func (s *TeamService) Create(e *manufacturing.MfgTeam) error {
	if e.Code == "" || e.Name == "" {
		return errors.New("编码与名称必填")
	}
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// 提取成员，先存班组再补成员
		members := e.Members
		e.Members = nil
		if err := tx.Create(e).Error; err != nil {
			return err
		}
		for i := range members {
			members[i].ID = 0
			members[i].TeamID = e.ID
			if err := tx.Create(&members[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *TeamService) Update(e *manufacturing.MfgTeam) error {
	var old manufacturing.MfgTeam
	if err := global.GVA_DB.First(&old, e.ID).Error; err != nil {
		return err
	}
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// 全量替换成员
		if err := tx.Where("team_id=?", e.ID).Delete(&manufacturing.MfgTeamMember{}).Error; err != nil {
			return err
		}
		members := e.Members
		e.Members = nil
		e.CreatedAt = old.CreatedAt
		if err := tx.Save(e).Error; err != nil {
			return err
		}
		for i := range members {
			members[i].ID = 0
			members[i].TeamID = e.ID
			if err := tx.Create(&members[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *TeamService) Delete(id uint) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("team_id=?", id).Delete(&manufacturing.MfgTeamMember{}).Error; err != nil {
			return err
		}
		return tx.Delete(&manufacturing.MfgTeam{}, id).Error
	})
}

func (s *TeamService) GetByID(id uint) (e manufacturing.MfgTeam, err error) {
	err = global.GVA_DB.Preload("Leader").Preload("Members.User").First(&e, id).Error
	return
}

type TeamSearch struct {
	request.PageInfo
	Keyword string `json:"keyword" form:"keyword"`
	Status  int8   `json:"status" form:"status"`
}

func (s *TeamService) GetList(info TeamSearch) (list []manufacturing.MfgTeam, total int64, err error) {
	if info.Page <= 0 {
		info.Page = 1
	}
	if info.PageSize <= 0 {
		info.PageSize = 10
	}
	db := global.GVA_DB.Model(&manufacturing.MfgTeam{})
	if info.Keyword != "" {
		db = db.Where("code LIKE ? OR name LIKE ?", "%"+info.Keyword+"%", "%"+info.Keyword+"%")
	}
	if info.Status > 0 {
		db = db.Where("status=?", info.Status)
	}
	if err = db.Count(&total).Error; err != nil {
		return
	}
	err = db.Preload("Leader").
		Limit(info.PageSize).Offset((info.Page - 1) * info.PageSize).
		Order("id DESC").Find(&list).Error
	return
}

// ============================================================
// 3. 工艺路线
// ============================================================
type RoutingService struct{}

func (s *RoutingService) Create(e *manufacturing.MfgRouting) error {
	if e.ProductID == 0 || e.Version == "" {
		return errors.New("成品与版本号必填")
	}
	var cnt int64
	global.GVA_DB.Model(&manufacturing.MfgRouting{}).
		Where("product_id=? AND version=?", e.ProductID, e.Version).Count(&cnt)
	if cnt > 0 {
		return errors.New("该成品已存在相同版本号的工艺路线")
	}
	if len(e.Items) == 0 {
		return errors.New("至少添加一道工序")
	}
	return global.GVA_DB.Create(e).Error
}

func (s *RoutingService) Update(e *manufacturing.MfgRouting) error {
	var old manufacturing.MfgRouting
	if err := global.GVA_DB.First(&old, e.ID).Error; err != nil {
		return err
	}
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("routing_id=?", e.ID).Delete(&manufacturing.MfgRoutingItem{}).Error; err != nil {
			return err
		}
		for i := range e.Items {
			e.Items[i].ID = 0
			e.Items[i].RoutingID = e.ID
		}
		e.CreatedAt = old.CreatedAt
		return tx.Session(&gorm.Session{FullSaveAssociations: true}).Save(e).Error
	})
}

func (s *RoutingService) Delete(id uint) error {
	// 被工单引用过的不可删
	var cnt int64
	global.GVA_DB.Model(&manufacturing.MfgWorkOrder{}).Where("routing_id=?", id).Count(&cnt)
	if cnt > 0 {
		return errors.New("该工艺路线已被工单引用，不可删除（请改为停用）")
	}
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("routing_id=?", id).Delete(&manufacturing.MfgRoutingItem{}).Error; err != nil {
			return err
		}
		return tx.Delete(&manufacturing.MfgRouting{}, id).Error
	})
}

func (s *RoutingService) SetStatus(id uint, status int8) error {
	return global.GVA_DB.Model(&manufacturing.MfgRouting{}).Where("id=?", id).
		Update("status", status).Error
}

func (s *RoutingService) GetByID(id uint) (e manufacturing.MfgRouting, err error) {
	err = global.GVA_DB.Preload("Product").
		Preload("Items.Workstation").
		Preload("Items", func(db *gorm.DB) *gorm.DB { return db.Order("seq") }).
		First(&e, id).Error
	return
}

type RoutingSearch struct {
	request.PageInfo
	Status    int8 `json:"status" form:"status"`
	ProductID uint `json:"productId" form:"productId"`
}

func (s *RoutingService) GetList(info RoutingSearch) (list []manufacturing.MfgRouting, total int64, err error) {
	if info.Page <= 0 {
		info.Page = 1
	}
	if info.PageSize <= 0 {
		info.PageSize = 10
	}
	db := global.GVA_DB.Model(&manufacturing.MfgRouting{})
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

// GetActiveByProduct 取某成品所有启用工艺路线（前端工单页面用）
func (s *RoutingService) GetActiveByProduct(productID uint) (list []manufacturing.MfgRouting, err error) {
	err = global.GVA_DB.Where("product_id=? AND status=?", productID, manufacturing.RoutingStatusActive).
		Preload("Items", func(db *gorm.DB) *gorm.DB { return db.Order("seq") }).
		Order("id DESC").Find(&list).Error
	return
}
