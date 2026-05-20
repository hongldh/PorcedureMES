package inventory

import "github.com/flipped-aurora/gin-vue-admin/server/global"

// InvProduct 商品（SKU）
//
// 关系字段说明：
//   CategoryID/UnitID 是外键 ID；Category/Unit 是 GORM 关联对象
//   GORM 通过 foreignKey 标签明确告诉它"用哪个字段做关联键"
//   查询时使用 Preload("Category") 加载关联数据
type InvProduct struct {
	global.GVA_MODEL
	Code          string  `json:"code" form:"code" gorm:"size:64;uniqueIndex;not null;comment:商品编码/SKU"`
	Name          string  `json:"name" form:"name" gorm:"size:128;not null;index;comment:商品名称"`
	Spec          string  `json:"spec" form:"spec" gorm:"size:128;comment:规格型号"`
	Barcode       string  `json:"barcode" form:"barcode" gorm:"size:64;index;comment:条形码"`
	CategoryID    uint    `json:"categoryId" form:"categoryId" gorm:"index;comment:分类ID"`
	UnitID        uint    `json:"unitId" form:"unitId" gorm:"comment:单位ID"`
	PurchasePrice float64 `json:"purchasePrice" form:"purchasePrice" gorm:"type:decimal(18,4);default:0;comment:参考采购价"`
	SalePrice     float64 `json:"salePrice" form:"salePrice" gorm:"type:decimal(18,4);default:0;comment:参考销售价"`
	SafetyStock   float64 `json:"safetyStock" form:"safetyStock" gorm:"type:decimal(18,4);default:0;comment:安全库存"`
	Status        int8    `json:"status" form:"status" gorm:"default:1;comment:状态 1启用 2停用"`
	Remark        string  `json:"remark" form:"remark" gorm:"size:255;comment:备注"`

	// 关联对象（不存数据库列，仅用于 Preload 加载）
	Category *InvCategory `json:"category,omitempty" gorm:"foreignKey:CategoryID;references:ID"`
	Unit     *InvUnit     `json:"unit,omitempty" gorm:"foreignKey:UnitID;references:ID"`
}

func (InvProduct) TableName() string { return "inv_product" }
