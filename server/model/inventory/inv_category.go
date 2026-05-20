package inventory

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// InvCategory 商品分类（树形结构）
//
// === Go 学习要点：GORM 标签详解 ===
// GORM 通过结构体字段后的 `tag` 来描述列属性。这里把所有常见 tag 一次讲透：
//
// - column           : 自定义数据库列名（不写就用 snake_case 默认转换）
// - type             : 指定数据库列类型（不同方言可能不同；如 varchar(64)）
// - size             : 简化版 type，仅控制长度
// - primaryKey       : 主键
// - uniqueIndex      : 唯一索引（可命名：uniqueIndex:idx_xxx）
// - index            : 普通索引
// - not null         : 非空约束
// - default          : 默认值
// - comment          : 列注释（最终会写到数据库 DDL 注释里）
// - <-:create        : 仅创建时写入，更新时忽略
// - -                : 完全忽略（不映射到数据库）
//
// JSON 标签：控制 HTTP 序列化字段名
// FORM 标签：控制 Gin 从 query/form 绑定时的字段名
// BINDING 标签：Gin validator 校验规则（required/min/max/...）
type InvCategory struct {
	global.GVA_MODEL                                                                       // 嵌入基类（ID、CreatedAt、UpdatedAt、DeletedAt）
	ParentID         uint   `json:"parentId" form:"parentId" gorm:"index;default:0;comment:父分类ID(0=根)"` // 父级ID，0 代表根节点
	Code             string `json:"code" form:"code" gorm:"size:64;uniqueIndex;not null;comment:分类编码"`   // 分类编码（业务唯一）
	Name             string `json:"name" form:"name" gorm:"size:128;not null;comment:分类名称"`              // 分类名
	Sort             int    `json:"sort" form:"sort" gorm:"default:0;comment:排序"`                        // 同级排序
	Status           int8   `json:"status" form:"status" gorm:"default:1;comment:状态 1启用 2停用"`            // 状态
	Remark           string `json:"remark" form:"remark" gorm:"size:255;comment:备注"`                     // 备注
}

// TableName 显式指定数据库表名
// 不写时 GORM 默认会用 "inv_categories"（复数）；我们统一用单数表名。
func (InvCategory) TableName() string {
	return "inv_category"
}
