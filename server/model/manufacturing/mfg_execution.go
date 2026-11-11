// MES 阶段 5b：工单工序（快照 + 执行态）+ 报工记录
//
// === 学习要点 ===
// 1. 工序快照：下达时从 mfg_routing_item 复制到 mfg_work_order_op，之后修改工艺路线不影响在制工单
// 2. 报工"即时生效"：没有草稿态，提交即扣工序进度；错了用作废冲回
// 3. 工序顺序约束在 service 层实现，表只记事实
package manufacturing

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
)

// 工单工序状态
const (
	OpStatusPending    int8 = 1 // 未派工
	OpStatusDispatched int8 = 2 // 已派工
	OpStatusInProgress int8 = 3 // 进行中（有报工）
	OpStatusFinished   int8 = 4 // 已完工（合格数 ≥ 计划数）
)

// 报工记录状态
const (
	ReportStatusDone int8 = 3 // 已生效
	ReportStatusVoid int8 = 4 // 已作废
)

// MfgWorkOrderOp 工单工序（工艺路线快照 + 执行状态）
type MfgWorkOrderOp struct {
	global.GVA_MODEL
	WorkOrderID uint `json:"workOrderId" gorm:"index;not null"`

	// --- 快照自 MfgRoutingItem ---
	Seq           int     `json:"seq" gorm:"not null;comment:工序号"`
	Code          string  `json:"code" gorm:"size:64;not null"`
	Name          string  `json:"name" gorm:"size:128;not null"`
	WorkstationID uint    `json:"workstationId" gorm:"comment:默认工位"`
	StdHours      float64 `json:"stdHours" gorm:"type:decimal(10,4);default:0"`
	IsKey         bool    `json:"isKey" gorm:"default:false"`

	// --- 执行态 ---
	PlanQty     float64 `json:"planQty" gorm:"type:decimal(18,4);not null;comment:= 工单计划数"`
	DoneQty     float64 `json:"doneQty" gorm:"type:decimal(18,4);default:0;comment:累计合格"`
	RejectQty   float64 `json:"rejectQty" gorm:"type:decimal(18,4);default:0;comment:累计不良"`
	ActualHours float64 `json:"actualHours" gorm:"type:decimal(10,4);default:0;comment:累计工时"`

	// --- 派工 ---
	AssignedUserID uint `json:"assignedUserId" gorm:"index;comment:指派员工（可空）"`
	TeamID         uint `json:"teamId" gorm:"index;comment:指派班组（可空）"`
	Status         int8 `json:"status" gorm:"default:1;index;comment:1未派工 2已派工 3进行中 4已完工"`

	Workstation  *MfgWorkstation `json:"workstation,omitempty" gorm:"foreignKey:WorkstationID"`
	AssignedUser *system.SysUser `json:"assignedUser,omitempty" gorm:"foreignKey:AssignedUserID"`
	Team         *MfgTeam        `json:"team,omitempty" gorm:"foreignKey:TeamID"`
}

func (MfgWorkOrderOp) TableName() string { return "mfg_work_order_op" }

// MfgWorkReport 报工记录
type MfgWorkReport struct {
	global.GVA_MODEL
	DocNo         string  `json:"docNo" gorm:"size:64;uniqueIndex;not null"`
	WorkOrderID   uint    `json:"workOrderId" gorm:"index;not null"`
	OpID          uint    `json:"opId" gorm:"index;not null;comment:mfg_work_order_op.ID"`
	UserID        uint    `json:"userId" gorm:"index;not null;comment:报工人"`
	WorkstationID uint    `json:"workstationId" gorm:"comment:实际工位"`
	GoodQty       float64 `json:"goodQty" gorm:"type:decimal(18,4);default:0;comment:合格数"`
	RejectQty     float64 `json:"rejectQty" gorm:"type:decimal(18,4);default:0;comment:不良数"`
	Hours         float64 `json:"hours" gorm:"type:decimal(10,4);default:0;comment:工时(小时)"`
	ReportTime    string  `json:"reportTime" gorm:"size:19;comment:报工时间 YYYY-MM-DD HH:mm:ss"`
	Status        int8    `json:"status" gorm:"default:3;index;comment:3已生效 4作废"`
	Remark        string  `json:"remark" gorm:"size:255"`

	WorkOrder   *MfgWorkOrder   `json:"workOrder,omitempty" gorm:"foreignKey:WorkOrderID"`
	Op          *MfgWorkOrderOp `json:"op,omitempty" gorm:"foreignKey:OpID"`
	User        *system.SysUser `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Workstation *MfgWorkstation `json:"workstation,omitempty" gorm:"foreignKey:WorkstationID"`
}

func (MfgWorkReport) TableName() string { return "mfg_work_report" }
