// Package manufacturing MES 制造执行模块 - 数据模型
//
// === 学习要点 ===
// 1. MES 不重写库存逻辑，所有库存动作复用 inventory.StockEngine
// 2. BOM "快照"模式：工单下达时把 BOM 行复制到 mfg_work_order_item，避免后续改 BOM 影响历史工单
// 3. 状态机用 int8 常量明确表达，service 层入口校验
package manufacturing

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/inventory"
)

// ==================== 状态常量 ====================

// 工单状态
const (
	WoStatusDraft      int8 = 1 // 草稿（可改可删）
	WoStatusReleased   int8 = 2 // 已下达（物料需求已生成，可领料）
	WoStatusInProgress int8 = 3 // 生产中（首次领料后自动进入）
	WoStatusFinished   int8 = 4 // 已完工（实际数量达到计划数）
	WoStatusClosed     int8 = 5 // 已关闭（人工结案，不可再操作）
	WoStatusVoid       int8 = 9 // 已作废
)

// BOM 状态
const (
	BomStatusActive   int8 = 1 // 启用
	BomStatusInactive int8 = 2 // 停用
)

// 领料/完工单据状态（与进销存一致）
const (
	DocStatusDraft int8 = 1
	DocStatusDone  int8 = 3
	DocStatusVoid  int8 = 4
)

// ==================== 1. BOM 物料清单 ====================

// MfgBom BOM 头表
//
// 同一成品支持多版本，工单引用 BOM 后会快照到工单行，避免后续改 BOM 影响历史
type MfgBom struct {
	global.GVA_MODEL
	ProductID     uint    `json:"productId" gorm:"index;not null;comment:成品ID"`
	Version       string  `json:"version" gorm:"size:32;not null;comment:版本号 v1.0"`
	EffectiveDate string  `json:"effectiveDate" gorm:"size:10;comment:生效日期 YYYY-MM-DD"`
	Status        int8    `json:"status" gorm:"default:1;index;comment:1启用 2停用"`
	Remark        string  `json:"remark" gorm:"size:255"`
	OpUserID      uint    `json:"opUserId"`

	Product *inventory.InvProduct `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	Items   []MfgBomItem          `json:"items,omitempty" gorm:"foreignKey:BomID"`
}

func (MfgBom) TableName() string { return "mfg_bom" }

// MfgBomItem BOM 行（子项物料）
type MfgBomItem struct {
	global.GVA_MODEL
	BomID     uint    `json:"bomId" gorm:"index;not null"`
	ProductID uint    `json:"productId" gorm:"not null;comment:子物料ID"`
	Quantity  float64 `json:"quantity" gorm:"type:decimal(18,4);not null;comment:单位用量"`
	LossRate  float64 `json:"lossRate" gorm:"type:decimal(8,4);default:0;comment:损耗率(0.05=5%)"`
	Remark    string  `json:"remark" gorm:"size:255"`

	Product *inventory.InvProduct `json:"product,omitempty" gorm:"foreignKey:ProductID"`
}

func (MfgBomItem) TableName() string { return "mfg_bom_item" }

// ==================== 2. 生产工单 ====================

// MfgWorkOrder 工单头
type MfgWorkOrder struct {
	global.GVA_MODEL
	DocNo        string  `json:"docNo" gorm:"size:64;uniqueIndex;not null"`
	ProductID    uint    `json:"productId" gorm:"index;not null;comment:成品ID"`
	BomID        uint    `json:"bomId" gorm:"comment:引用的BOM版本ID"`
	RoutingID    uint    `json:"routingId" gorm:"comment:引用的工艺路线版本ID（5a扩展，可空）"`
	PlanQty      float64 `json:"planQty" gorm:"type:decimal(18,4);not null;comment:计划数量"`
	ActualQty    float64 `json:"actualQty" gorm:"type:decimal(18,4);default:0;comment:已完工数量"`
	IssuedAmount float64 `json:"issuedAmount" gorm:"type:decimal(18,4);default:0;comment:累计领料金额"`
	PlanStart    string  `json:"planStart" gorm:"size:10"`
	PlanEnd      string  `json:"planEnd" gorm:"size:10"`
	WarehouseID  uint    `json:"warehouseId" gorm:"comment:默认入库仓"`
	Status       int8    `json:"status" gorm:"default:1;index;comment:1草稿 2已下达 3生产中 4已完工 5已关闭 9作废"`
	Remark       string  `json:"remark" gorm:"size:255"`
	OpUserID     uint    `json:"opUserId"`

	Product   *inventory.InvProduct   `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	Bom       *MfgBom                 `json:"bom,omitempty" gorm:"foreignKey:BomID"`
	Routing   *MfgRouting             `json:"routing,omitempty" gorm:"foreignKey:RoutingID"`
	Warehouse *inventory.InvWarehouse `json:"warehouse,omitempty" gorm:"foreignKey:WarehouseID"`
	Items     []MfgWorkOrderItem      `json:"items,omitempty" gorm:"foreignKey:WorkOrderID"`
}

func (MfgWorkOrder) TableName() string { return "mfg_work_order" }

// MfgWorkOrderItem 工单物料需求（下达时按 BOM × PlanQty × (1+损耗率) 展开）
type MfgWorkOrderItem struct {
	global.GVA_MODEL
	WorkOrderID uint    `json:"workOrderId" gorm:"index;not null"`
	ProductID   uint    `json:"productId" gorm:"not null"`
	RequiredQty float64 `json:"requiredQty" gorm:"type:decimal(18,4);not null;comment:计划用量"`
	IssuedQty   float64 `json:"issuedQty" gorm:"type:decimal(18,4);default:0;comment:累计领料数量"`
	Remark      string  `json:"remark" gorm:"size:255"`

	Product *inventory.InvProduct `json:"product,omitempty" gorm:"foreignKey:ProductID"`
}

func (MfgWorkOrderItem) TableName() string { return "mfg_work_order_item" }

// ==================== 3. 领料单 ====================

// MfgIssue 领料单头
type MfgIssue struct {
	global.GVA_MODEL
	DocNo       string  `json:"docNo" gorm:"size:64;uniqueIndex;not null"`
	WorkOrderID uint    `json:"workOrderId" gorm:"index;not null;comment:关联工单"`
	WarehouseID uint    `json:"warehouseId" gorm:"index;not null;comment:领料仓"`
	BizDate     string  `json:"bizDate" gorm:"size:10;not null"`
	TotalAmount float64 `json:"totalAmount" gorm:"type:decimal(18,4);default:0;comment:本次领料总金额"`
	Status      int8    `json:"status" gorm:"default:1;index;comment:1草稿 3已执行 4作废"`
	Remark      string  `json:"remark" gorm:"size:255"`
	OpUserID    uint    `json:"opUserId"`

	WorkOrder *MfgWorkOrder           `json:"workOrder,omitempty" gorm:"foreignKey:WorkOrderID"`
	Warehouse *inventory.InvWarehouse `json:"warehouse,omitempty" gorm:"foreignKey:WarehouseID"`
	Items     []MfgIssueItem          `json:"items,omitempty" gorm:"foreignKey:IssueID"`
}

func (MfgIssue) TableName() string { return "mfg_issue" }

// MfgIssueItem 领料单行
type MfgIssueItem struct {
	global.GVA_MODEL
	IssueID    uint    `json:"issueId" gorm:"index;not null"`
	ProductID  uint    `json:"productId" gorm:"not null"`
	LocationID uint    `json:"locationId" gorm:"comment:出库货位"`
	Quantity   float64 `json:"quantity" gorm:"type:decimal(18,4);not null"`
	UnitCost   float64 `json:"unitCost" gorm:"type:decimal(18,4);default:0;comment:出库时取得的均价"`
	Amount     float64 `json:"amount" gorm:"type:decimal(18,4);default:0"`
	Remark     string  `json:"remark" gorm:"size:255"`

	Product  *inventory.InvProduct  `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	Location *inventory.InvLocation `json:"location,omitempty" gorm:"foreignKey:LocationID"`
}

func (MfgIssueItem) TableName() string { return "mfg_issue_item" }

// ==================== 4. 完工入库单 ====================

// MfgReceipt 完工入库单
//
// 一张完工单只入库一个成品（=工单上的成品）。
// UnitCost 在 Confirm 时由 service 层计算：
//   unitCost = (该工单累计领料金额 + 本次领料金额已含) / 本次完工数量
type MfgReceipt struct {
	global.GVA_MODEL
	DocNo       string  `json:"docNo" gorm:"size:64;uniqueIndex;not null"`
	WorkOrderID uint    `json:"workOrderId" gorm:"index;not null"`
	WarehouseID uint    `json:"warehouseId" gorm:"index;not null;comment:入库仓"`
	LocationID  uint    `json:"locationId" gorm:"comment:入库货位"`
	BizDate     string  `json:"bizDate" gorm:"size:10;not null"`
	Quantity    float64 `json:"quantity" gorm:"type:decimal(18,4);not null;comment:完工数量"`
	UnitCost    float64 `json:"unitCost" gorm:"type:decimal(18,4);default:0;comment:Confirm时计算"`
	TotalAmount float64 `json:"totalAmount" gorm:"type:decimal(18,4);default:0"`
	Status      int8    `json:"status" gorm:"default:1;index;comment:1草稿 3已执行 4作废"`
	Remark      string  `json:"remark" gorm:"size:255"`
	OpUserID    uint    `json:"opUserId"`

	WorkOrder *MfgWorkOrder           `json:"workOrder,omitempty" gorm:"foreignKey:WorkOrderID"`
	Warehouse *inventory.InvWarehouse `json:"warehouse,omitempty" gorm:"foreignKey:WarehouseID"`
	Location  *inventory.InvLocation  `json:"location,omitempty" gorm:"foreignKey:LocationID"`
	// 注：成品信息通过 WorkOrder.Product 获取，本表不直接挂 Product 关联
}

func (MfgReceipt) TableName() string { return "mfg_receipt" }
