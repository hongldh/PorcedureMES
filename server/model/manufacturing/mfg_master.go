// MES 阶段 5a 主数据模型：工位 / 班组 / 班组成员 / 工艺路线 / 工序
//
// === 学习要点 ===
// 1. 工位与班组是基础主数据，互不耦合（工位是物理位置，班组是人员组织）
// 2. 工艺路线与 BOM 同结构：头多版本 + 行明细，工单引用时复制快照
// 3. 班组成员用关联表，避免多对多关联表自动命名导致索引重复
package manufacturing

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/inventory"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
)

// ==================== 工位类型 ====================
const (
	WorkstationManual    int8 = 1 // 人工
	WorkstationEquipment int8 = 2 // 设备
	WorkstationOutsource int8 = 3 // 外协
)

// ==================== 工艺路线状态 ====================
const (
	RoutingStatusActive   int8 = 1 // 启用
	RoutingStatusInactive int8 = 2 // 停用
)

// ==================== 工位 / 班组共享状态 ====================
const (
	BasicStatusEnabled  int8 = 1 // 启用
	BasicStatusDisabled int8 = 2 // 停用
)

// ============================================================
// 1. 工位
// ============================================================

// MfgWorkstation 工位（车间内的物理工作单元）
type MfgWorkstation struct {
	global.GVA_MODEL
	Code        string  `json:"code" gorm:"size:64;uniqueIndex;not null;comment:工位编码"`
	Name        string  `json:"name" gorm:"size:128;not null"`
	Type        int8    `json:"type" gorm:"default:1;index;comment:1人工 2设备 3外协"`
	Capacity    float64 `json:"capacity" gorm:"type:decimal(18,4);default:0;comment:每小时标称产能"`
	WarehouseID uint    `json:"warehouseId" gorm:"index;comment:所属仓库（可空）"`
	Status      int8    `json:"status" gorm:"default:1;index"`
	Remark      string  `json:"remark" gorm:"size:255"`

	Warehouse *inventory.InvWarehouse `json:"warehouse,omitempty" gorm:"foreignKey:WarehouseID"`
}

func (MfgWorkstation) TableName() string { return "mfg_workstation" }

// ============================================================
// 2. 班组
// ============================================================

// MfgTeam 班组
type MfgTeam struct {
	global.GVA_MODEL
	Code     string `json:"code" gorm:"size:64;uniqueIndex;not null;comment:班组编码"`
	Name     string `json:"name" gorm:"size:128;not null"`
	LeaderID uint   `json:"leaderId" gorm:"comment:班组长（sys_user.ID）"`
	Status   int8   `json:"status" gorm:"default:1;index"`
	Remark   string `json:"remark" gorm:"size:255"`

	Leader  *system.SysUser  `json:"leader,omitempty" gorm:"foreignKey:LeaderID"`
	Members []MfgTeamMember  `json:"members,omitempty" gorm:"foreignKey:TeamID"`
}

func (MfgTeam) TableName() string { return "mfg_team" }

// MfgTeamMember 班组成员（人 ↔ 班组的关联表）
type MfgTeamMember struct {
	global.GVA_MODEL
	TeamID uint `json:"teamId" gorm:"index;not null"`
	UserID uint `json:"userId" gorm:"index;not null;comment:sys_user.ID"`

	User *system.SysUser `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

func (MfgTeamMember) TableName() string { return "mfg_team_member" }

// ============================================================
// 3. 工艺路线 + 工序
// ============================================================

// MfgRouting 工艺路线头
//
// 与 BOM 同结构：同一成品多版本，工单引用时复制行到工单工序快照（5b 实现）
type MfgRouting struct {
	global.GVA_MODEL
	ProductID     uint   `json:"productId" gorm:"index;not null;comment:成品ID"`
	Version       string `json:"version" gorm:"size:32;not null;comment:版本号 v1.0"`
	EffectiveDate string `json:"effectiveDate" gorm:"size:10"`
	Status        int8   `json:"status" gorm:"default:1;index"`
	Remark        string `json:"remark" gorm:"size:255"`

	Product *inventory.InvProduct `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	Items   []MfgRoutingItem      `json:"items,omitempty" gorm:"foreignKey:RoutingID"`
}

func (MfgRouting) TableName() string { return "mfg_routing" }

// MfgRoutingItem 工艺路线行 = 工序定义
type MfgRoutingItem struct {
	global.GVA_MODEL
	RoutingID     uint    `json:"routingId" gorm:"index;not null"`
	Seq           int     `json:"seq" gorm:"not null;comment:工序号 10/20/30 便于插入"`
	Code          string  `json:"code" gorm:"size:64;not null;comment:工序编码 OP010"`
	Name          string  `json:"name" gorm:"size:128;not null;comment:工序名称"`
	WorkstationID uint    `json:"workstationId" gorm:"comment:默认工位"`
	StdHours      float64 `json:"stdHours" gorm:"type:decimal(10,4);default:0;comment:标准工时(小时/单位)"`
	IsKey         bool    `json:"isKey" gorm:"default:false;comment:是否关键工序"`
	Description   string  `json:"description" gorm:"type:text;comment:作业指导"`
	Remark        string  `json:"remark" gorm:"size:255"`

	Workstation *MfgWorkstation `json:"workstation,omitempty" gorm:"foreignKey:WorkstationID"`
}

func (MfgRoutingItem) TableName() string { return "mfg_routing_item" }
