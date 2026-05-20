package inventory

import "github.com/flipped-aurora/gin-vue-admin/server/global"

// InvSupplier 供应商
type InvSupplier struct {
	global.GVA_MODEL
	Code    string `json:"code" form:"code" gorm:"size:64;uniqueIndex;not null;comment:供应商编码"`
	Name    string `json:"name" form:"name" gorm:"size:128;not null;index;comment:供应商名称"`
	Contact string `json:"contact" form:"contact" gorm:"size:64;comment:联系人"`
	Phone   string `json:"phone" form:"phone" gorm:"size:32;comment:电话"`
	Email   string `json:"email" form:"email" gorm:"size:128;comment:邮箱"`
	Address string `json:"address" form:"address" gorm:"size:255;comment:地址"`
	Status  int8   `json:"status" form:"status" gorm:"default:1;comment:状态 1启用 2停用"`
	Remark  string `json:"remark" form:"remark" gorm:"size:255;comment:备注"`
}

func (InvSupplier) TableName() string { return "inv_supplier" }

// InvCustomer 客户（独立于 GVA 自带的 ExaCustomer）
type InvCustomer struct {
	global.GVA_MODEL
	Code    string `json:"code" form:"code" gorm:"size:64;uniqueIndex;not null;comment:客户编码"`
	Name    string `json:"name" form:"name" gorm:"size:128;not null;index;comment:客户名称"`
	Contact string `json:"contact" form:"contact" gorm:"size:64;comment:联系人"`
	Phone   string `json:"phone" form:"phone" gorm:"size:32;comment:电话"`
	Email   string `json:"email" form:"email" gorm:"size:128;comment:邮箱"`
	Address string `json:"address" form:"address" gorm:"size:255;comment:地址"`
	Status  int8   `json:"status" form:"status" gorm:"default:1;comment:状态 1启用 2停用"`
	Remark  string `json:"remark" form:"remark" gorm:"size:255;comment:备注"`
}

func (InvCustomer) TableName() string { return "inv_customer" }
