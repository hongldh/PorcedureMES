package inventory

import "github.com/flipped-aurora/gin-vue-admin/server/global"

// 单据状态常量
const (
	DocStatusDraft   int8 = 1 // 草稿（可编辑/删除）
	DocStatusDone    int8 = 3 // 已执行（已写库存流水，不可改）
	DocStatusVoid    int8 = 4 // 已作废
)

// ==================== 采购入库 ====================

// InvPurchaseOrder 采购入库单（单头）
type InvPurchaseOrder struct {
	global.GVA_MODEL
	DocNo       string  `json:"docNo" gorm:"size:64;uniqueIndex;not null;comment:单号"`
	SupplierID  uint    `json:"supplierId" gorm:"index;not null;comment:供应商ID"`
	WarehouseID uint    `json:"warehouseId" gorm:"index;not null;comment:入库仓库ID"`
	BizDate     string  `json:"bizDate" gorm:"size:10;not null;comment:业务日期 YYYY-MM-DD"`
	Status      int8    `json:"status" gorm:"default:1;index;comment:状态 1草稿 3已执行 4作废"`
	TotalAmount float64 `json:"totalAmount" gorm:"type:decimal(18,4);default:0;comment:合计金额"`
	Remark      string  `json:"remark" gorm:"size:255;comment:备注"`
	OpUserID    uint    `json:"opUserId" gorm:"comment:操作人"`

	Supplier  *InvSupplier             `json:"supplier,omitempty" gorm:"foreignKey:SupplierID"`
	Warehouse *InvWarehouse            `json:"warehouse,omitempty" gorm:"foreignKey:WarehouseID"`
	Items     []InvPurchaseOrderItem   `json:"items,omitempty" gorm:"foreignKey:OrderID"`
}

func (InvPurchaseOrder) TableName() string { return "inv_purchase_order" }

// InvPurchaseOrderItem 采购入库单明细
type InvPurchaseOrderItem struct {
	global.GVA_MODEL
	OrderID    uint    `json:"orderId" gorm:"index;not null"`
	ProductID  uint    `json:"productId" gorm:"not null"`
	LocationID uint    `json:"locationId" gorm:"not null;comment:入库货位"`
	Quantity   float64 `json:"quantity" gorm:"type:decimal(18,4);not null"`
	UnitPrice  float64 `json:"unitPrice" gorm:"type:decimal(18,4);default:0"`
	Amount     float64 `json:"amount" gorm:"type:decimal(18,4);default:0"`
	Remark     string  `json:"remark" gorm:"size:255"`

	Product  *InvProduct  `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	Location *InvLocation `json:"location,omitempty" gorm:"foreignKey:LocationID"`
}

func (InvPurchaseOrderItem) TableName() string { return "inv_purchase_order_item" }

// ==================== 采购退货 ====================

type InvPurchaseReturn struct {
	global.GVA_MODEL
	DocNo       string  `json:"docNo" gorm:"size:64;uniqueIndex;not null"`
	SupplierID  uint    `json:"supplierId" gorm:"index;not null"`
	WarehouseID uint    `json:"warehouseId" gorm:"index;not null"`
	BizDate     string  `json:"bizDate" gorm:"size:10;not null"`
	Status      int8    `json:"status" gorm:"default:1;index"`
	TotalAmount float64 `json:"totalAmount" gorm:"type:decimal(18,4);default:0"`
	Remark      string  `json:"remark" gorm:"size:255"`
	OpUserID    uint    `json:"opUserId"`

	Supplier  *InvSupplier             `json:"supplier,omitempty" gorm:"foreignKey:SupplierID"`
	Warehouse *InvWarehouse            `json:"warehouse,omitempty" gorm:"foreignKey:WarehouseID"`
	Items     []InvPurchaseReturnItem  `json:"items,omitempty" gorm:"foreignKey:OrderID"`
}

func (InvPurchaseReturn) TableName() string { return "inv_purchase_return" }

type InvPurchaseReturnItem struct {
	global.GVA_MODEL
	OrderID    uint    `json:"orderId" gorm:"index;not null"`
	ProductID  uint    `json:"productId" gorm:"not null"`
	LocationID uint    `json:"locationId" gorm:"not null"`
	Quantity   float64 `json:"quantity" gorm:"type:decimal(18,4);not null"`
	UnitPrice  float64 `json:"unitPrice" gorm:"type:decimal(18,4);default:0"`
	Amount     float64 `json:"amount" gorm:"type:decimal(18,4);default:0"`
	Remark     string  `json:"remark" gorm:"size:255"`

	Product  *InvProduct  `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	Location *InvLocation `json:"location,omitempty" gorm:"foreignKey:LocationID"`
}

func (InvPurchaseReturnItem) TableName() string { return "inv_purchase_return_item" }

// ==================== 销售出库 ====================

type InvSalesOrder struct {
	global.GVA_MODEL
	DocNo       string  `json:"docNo" gorm:"size:64;uniqueIndex;not null"`
	CustomerID  uint    `json:"customerId" gorm:"index;not null"`
	WarehouseID uint    `json:"warehouseId" gorm:"index;not null"`
	BizDate     string  `json:"bizDate" gorm:"size:10;not null"`
	Status      int8    `json:"status" gorm:"default:1;index"`
	TotalAmount float64 `json:"totalAmount" gorm:"type:decimal(18,4);default:0"`
	Remark      string  `json:"remark" gorm:"size:255"`
	OpUserID    uint    `json:"opUserId"`

	Customer  *InvCustomer        `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
	Warehouse *InvWarehouse       `json:"warehouse,omitempty" gorm:"foreignKey:WarehouseID"`
	Items     []InvSalesOrderItem `json:"items,omitempty" gorm:"foreignKey:OrderID"`
}

func (InvSalesOrder) TableName() string { return "inv_sales_order" }

type InvSalesOrderItem struct {
	global.GVA_MODEL
	OrderID    uint    `json:"orderId" gorm:"index;not null"`
	ProductID  uint    `json:"productId" gorm:"not null"`
	LocationID uint    `json:"locationId" gorm:"not null"`
	Quantity   float64 `json:"quantity" gorm:"type:decimal(18,4);not null"`
	UnitPrice  float64 `json:"unitPrice" gorm:"type:decimal(18,4);default:0"`
	Amount     float64 `json:"amount" gorm:"type:decimal(18,4);default:0"`
	Remark     string  `json:"remark" gorm:"size:255"`

	Product  *InvProduct  `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	Location *InvLocation `json:"location,omitempty" gorm:"foreignKey:LocationID"`
}

func (InvSalesOrderItem) TableName() string { return "inv_sales_order_item" }

// ==================== 销售退货 ====================

type InvSalesReturn struct {
	global.GVA_MODEL
	DocNo       string  `json:"docNo" gorm:"size:64;uniqueIndex;not null"`
	CustomerID  uint    `json:"customerId" gorm:"index;not null"`
	WarehouseID uint    `json:"warehouseId" gorm:"index;not null"`
	BizDate     string  `json:"bizDate" gorm:"size:10;not null"`
	Status      int8    `json:"status" gorm:"default:1;index"`
	TotalAmount float64 `json:"totalAmount" gorm:"type:decimal(18,4);default:0"`
	Remark      string  `json:"remark" gorm:"size:255"`
	OpUserID    uint    `json:"opUserId"`

	Customer  *InvCustomer         `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
	Warehouse *InvWarehouse        `json:"warehouse,omitempty" gorm:"foreignKey:WarehouseID"`
	Items     []InvSalesReturnItem `json:"items,omitempty" gorm:"foreignKey:OrderID"`
}

func (InvSalesReturn) TableName() string { return "inv_sales_return" }

type InvSalesReturnItem struct {
	global.GVA_MODEL
	OrderID    uint    `json:"orderId" gorm:"index;not null"`
	ProductID  uint    `json:"productId" gorm:"not null"`
	LocationID uint    `json:"locationId" gorm:"not null"`
	Quantity   float64 `json:"quantity" gorm:"type:decimal(18,4);not null"`
	UnitPrice  float64 `json:"unitPrice" gorm:"type:decimal(18,4);default:0"`
	Amount     float64 `json:"amount" gorm:"type:decimal(18,4);default:0"`
	Remark     string  `json:"remark" gorm:"size:255"`

	Product  *InvProduct  `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	Location *InvLocation `json:"location,omitempty" gorm:"foreignKey:LocationID"`
}

func (InvSalesReturnItem) TableName() string { return "inv_sales_return_item" }

// ==================== 调拨 ====================

type InvTransfer struct {
	global.GVA_MODEL
	DocNo         string `json:"docNo" gorm:"size:64;uniqueIndex;not null"`
	FromWarehouse uint   `json:"fromWarehouse" gorm:"not null;comment:源仓库"`
	ToWarehouse   uint   `json:"toWarehouse" gorm:"not null;comment:目标仓库"`
	BizDate       string `json:"bizDate" gorm:"size:10;not null"`
	Status        int8   `json:"status" gorm:"default:1;index"`
	Remark        string `json:"remark" gorm:"size:255"`
	OpUserID      uint   `json:"opUserId"`

	FromWh *InvWarehouse     `json:"fromWh,omitempty" gorm:"foreignKey:FromWarehouse"`
	ToWh   *InvWarehouse     `json:"toWh,omitempty" gorm:"foreignKey:ToWarehouse"`
	Items  []InvTransferItem `json:"items,omitempty" gorm:"foreignKey:OrderID"`
}

func (InvTransfer) TableName() string { return "inv_transfer" }

type InvTransferItem struct {
	global.GVA_MODEL
	OrderID        uint    `json:"orderId" gorm:"index;not null"`
	ProductID      uint    `json:"productId" gorm:"not null"`
	FromLocationID uint    `json:"fromLocationId" gorm:"not null;comment:源货位"`
	ToLocationID   uint    `json:"toLocationId" gorm:"not null;comment:目标货位"`
	Quantity       float64 `json:"quantity" gorm:"type:decimal(18,4);not null"`
	Remark         string  `json:"remark" gorm:"size:255"`

	Product *InvProduct `json:"product,omitempty" gorm:"foreignKey:ProductID"`
}

func (InvTransferItem) TableName() string { return "inv_transfer_item" }

// ==================== 盘点 ====================

type InvStocktake struct {
	global.GVA_MODEL
	DocNo       string `json:"docNo" gorm:"size:64;uniqueIndex;not null"`
	WarehouseID uint   `json:"warehouseId" gorm:"index;not null"`
	BizDate     string `json:"bizDate" gorm:"size:10;not null"`
	Status      int8   `json:"status" gorm:"default:1;index"`
	Remark      string `json:"remark" gorm:"size:255"`
	OpUserID    uint   `json:"opUserId"`

	Warehouse *InvWarehouse      `json:"warehouse,omitempty" gorm:"foreignKey:WarehouseID"`
	Items     []InvStocktakeItem `json:"items,omitempty" gorm:"foreignKey:OrderID"`
}

func (InvStocktake) TableName() string { return "inv_stocktake" }

// InvStocktakeItem 盘点明细
//   BookQty   = 盘点时的账面数（确认时由系统快照写入）
//   ActualQty = 用户填入的实盘数
//   DiffQty   = ActualQty - BookQty （正数=盘盈，负数=盘亏）
type InvStocktakeItem struct {
	global.GVA_MODEL
	OrderID    uint    `json:"orderId" gorm:"index;not null"`
	ProductID  uint    `json:"productId" gorm:"not null"`
	LocationID uint    `json:"locationId" gorm:"not null"`
	BookQty    float64 `json:"bookQty" gorm:"type:decimal(18,4);default:0;comment:账面数"`
	ActualQty  float64 `json:"actualQty" gorm:"type:decimal(18,4);default:0;comment:实盘数"`
	DiffQty    float64 `json:"diffQty" gorm:"type:decimal(18,4);default:0;comment:差异(盘盈正/盘亏负)"`
	Remark     string  `json:"remark" gorm:"size:255"`

	Product  *InvProduct  `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	Location *InvLocation `json:"location,omitempty" gorm:"foreignKey:LocationID"`
}

func (InvStocktakeItem) TableName() string { return "inv_stocktake_item" }
