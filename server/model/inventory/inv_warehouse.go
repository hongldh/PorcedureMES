package inventory

import "github.com/flipped-aurora/gin-vue-admin/server/global"

// InvWarehouse 仓库
type InvWarehouse struct {
	global.GVA_MODEL
	Code    string `json:"code" form:"code" gorm:"size:32;uniqueIndex;not null;comment:仓库编码"`
	Name    string `json:"name" form:"name" gorm:"size:128;not null;comment:仓库名称"`
	Address string `json:"address" form:"address" gorm:"size:255;comment:地址"`
	Manager string `json:"manager" form:"manager" gorm:"size:64;comment:负责人"`
	Phone   string `json:"phone" form:"phone" gorm:"size:32;comment:联系电话"`
	Status  int8   `json:"status" form:"status" gorm:"default:1;comment:状态 1启用 2停用"`
	Remark  string `json:"remark" form:"remark" gorm:"size:255;comment:备注"`
}

func (InvWarehouse) TableName() string { return "inv_warehouse" }

// InvLocation 货位（属于某个仓库）
type InvLocation struct {
	global.GVA_MODEL
	WarehouseID uint   `json:"warehouseId" form:"warehouseId" gorm:"index;not null;comment:仓库ID"`
	Code        string `json:"code" form:"code" gorm:"size:64;not null;comment:货位编码"`
	Name        string `json:"name" form:"name" gorm:"size:128;not null;comment:货位名称"`
	Sort        int    `json:"sort" form:"sort" gorm:"default:0;comment:排序"`
	Status      int8   `json:"status" form:"status" gorm:"default:1;comment:状态"`
	Remark      string `json:"remark" form:"remark" gorm:"size:255;comment:备注"`

	Warehouse *InvWarehouse `json:"warehouse,omitempty" gorm:"foreignKey:WarehouseID;references:ID"`
}

func (InvLocation) TableName() string { return "inv_location" }
