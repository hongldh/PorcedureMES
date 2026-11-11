package inventory

import "github.com/flipped-aurora/gin-vue-admin/server/global"

// InvStock 库存余额（按 商品 × 仓库 × 货位 聚合）
//
// 唯一索引 (product_id, warehouse_id, location_id) 保证一个三元组只有一条余额记录
// 业务上由 StockEngine 通过 upsert 维护
type InvStock struct {
	global.GVA_MODEL
	ProductID   uint    `json:"productId" gorm:"uniqueIndex:uk_stock;not null;comment:商品ID"`
	WarehouseID uint    `json:"warehouseId" gorm:"uniqueIndex:uk_stock;not null;comment:仓库ID"`
	LocationID  uint    `json:"locationId" gorm:"uniqueIndex:uk_stock;not null;comment:货位ID"`
	Quantity    float64 `json:"quantity" gorm:"type:decimal(18,4);default:0;comment:库存数量"`
	AvgCost     float64 `json:"avgCost" gorm:"type:decimal(18,4);default:0;comment:移动加权平均成本"`

	Product   *InvProduct   `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	Warehouse *InvWarehouse `json:"warehouse,omitempty" gorm:"foreignKey:WarehouseID"`
	Location  *InvLocation  `json:"location,omitempty" gorm:"foreignKey:LocationID"`
}

func (InvStock) TableName() string { return "inv_stock" }

// 业务类型常量（用于流水的 biz_type 字段）
const (
	BizPurchaseIn      = "PURCHASE_IN"      // 采购入库
	BizPurchaseReturn  = "PURCHASE_RETURN"  // 采购退货（出）
	BizSalesOut        = "SALES_OUT"        // 销售出库
	BizSalesReturn     = "SALES_RETURN"     // 销售退货（入）
	BizTransferOut     = "TRANSFER_OUT"     // 调拨出
	BizTransferIn      = "TRANSFER_IN"      // 调拨入
	BizStocktakeAdjust = "STOCKTAKE_ADJUST" // 盘点调整
	BizMfgIssue        = "MFG_ISSUE"        // 生产领料（出）
	BizMfgReceipt      = "MFG_RECEIPT"      // 生产完工入库（入）
)

// 方向常量
const (
	DirIn  int8 = 1 // 入库
	DirOut int8 = 2 // 出库
)

// InvStockFlow 库存流水（核心审计表，永不删除/修改）
//
// 任何库存变动都先写一条流水，再 upsert inv_stock。
// 流水里保存"变动前/后"快照，便于对账和重算余额。
type InvStockFlow struct {
	global.GVA_MODEL
	BizType     string  `json:"bizType" gorm:"size:32;index;not null;comment:业务类型"`
	BizNo       string  `json:"bizNo" gorm:"size:64;index;not null;comment:业务单号"`
	BizItemID   uint    `json:"bizItemId" gorm:"comment:业务明细ID"`
	Direction   int8    `json:"direction" gorm:"not null;comment:方向 1入 2出"`
	ProductID   uint    `json:"productId" gorm:"index;not null;comment:商品ID"`
	WarehouseID uint    `json:"warehouseId" gorm:"index;not null;comment:仓库ID"`
	LocationID  uint    `json:"locationId" gorm:"not null;comment:货位ID"`
	Quantity    float64 `json:"quantity" gorm:"type:decimal(18,4);not null;comment:本次变动数量(正数)"`
	UnitCost    float64 `json:"unitCost" gorm:"type:decimal(18,4);default:0;comment:本次单位成本"`
	TotalAmount float64 `json:"totalAmount" gorm:"type:decimal(18,4);default:0;comment:本次金额"`
	BeforeQty   float64 `json:"beforeQty" gorm:"type:decimal(18,4);default:0;comment:变动前数量"`
	AfterQty    float64 `json:"afterQty" gorm:"type:decimal(18,4);default:0;comment:变动后数量"`
	BeforeCost  float64 `json:"beforeCost" gorm:"type:decimal(18,4);default:0;comment:变动前均价"`
	AfterCost   float64 `json:"afterCost" gorm:"type:decimal(18,4);default:0;comment:变动后均价"`
	Remark      string  `json:"remark" gorm:"size:255;comment:备注"`
	OpUserID    uint    `json:"opUserId" gorm:"comment:操作人ID"`

	Product   *InvProduct   `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	Warehouse *InvWarehouse `json:"warehouse,omitempty" gorm:"foreignKey:WarehouseID"`
}

func (InvStockFlow) TableName() string { return "inv_stock_flow" }

// InvDocSequence 单号序列表（每天每业务一条）
// 用 (biz_type, date_str) 做唯一键，每次 +1 取号
type InvDocSequence struct {
	ID      uint   `gorm:"primaryKey"`
	BizType string `gorm:"size:32;uniqueIndex:uk_seq;not null"`
	DateStr string `gorm:"size:8;uniqueIndex:uk_seq;not null;comment:YYYYMMDD"`
	Seq     int64  `gorm:"not null;default:0;comment:当日序号"`
}

func (InvDocSequence) TableName() string { return "inv_doc_sequence" }
