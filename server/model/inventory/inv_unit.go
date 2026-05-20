package inventory

import "github.com/flipped-aurora/gin-vue-admin/server/global"

// InvUnit 计量单位（个/箱/kg/...）
type InvUnit struct {
	global.GVA_MODEL
	Code   string `json:"code" form:"code" gorm:"size:32;uniqueIndex;not null;comment:单位编码"`
	Name   string `json:"name" form:"name" gorm:"size:64;not null;comment:单位名称"`
	Status int8   `json:"status" form:"status" gorm:"default:1;comment:状态 1启用 2停用"`
	Remark string `json:"remark" form:"remark" gorm:"size:255;comment:备注"`
}

func (InvUnit) TableName() string { return "inv_unit" }
