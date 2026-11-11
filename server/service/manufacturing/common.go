// Package manufacturing MES 业务服务层
package manufacturing

import (
	"errors"
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/inventory"
	"github.com/flipped-aurora/gin-vue-admin/server/model/manufacturing"
	invsvc "github.com/flipped-aurora/gin-vue-admin/server/service/inventory"
	"gorm.io/gorm"
)

// 复用进销存的 StockEngine 和 SequenceService
var (
	stockEngine = &invsvc.StockEngine{}
	seqService  = &invsvc.SequenceService{}
)

// ==================== 公共校验 ====================

func mustWoDraft(status int8) error {
	if status != manufacturing.WoStatusDraft {
		return errors.New("仅草稿状态可编辑/删除")
	}
	return nil
}

func mustDocDraft(status int8) error {
	if status != manufacturing.DocStatusDraft {
		return errors.New("仅草稿状态可编辑/删除")
	}
	return nil
}

func mustDocDone(status int8) error {
	if status != manufacturing.DocStatusDone {
		return errors.New("仅已执行状态可作废")
	}
	return nil
}

// 通用：单据搜索
type DocSearch struct {
	request.PageInfo
	Status    int8   `json:"status" form:"status"`
	DateStart string `json:"dateStart" form:"dateStart"`
	DateEnd   string `json:"dateEnd" form:"dateEnd"`
	DocNo     string `json:"docNo" form:"docNo"`
}

func applyDocFilter(db *gorm.DB, s DocSearch) *gorm.DB {
	if s.Status > 0 {
		db = db.Where("status = ?", s.Status)
	}
	if s.DateStart != "" {
		db = db.Where("biz_date >= ?", s.DateStart)
	}
	if s.DateEnd != "" {
		db = db.Where("biz_date <= ?", s.DateEnd)
	}
	if s.DocNo != "" {
		db = db.Where("doc_no LIKE ?", "%"+s.DocNo+"%")
	}
	return db
}

// 工单专用搜索（无 biz_date，有 product_id）
type WoSearch struct {
	request.PageInfo
	Status    int8   `json:"status" form:"status"`
	DocNo     string `json:"docNo" form:"docNo"`
	ProductID uint   `json:"productId" form:"productId"`
}

// 同时获取仓库的 InvLocation 引用以编译通过（防止 manufacturing 包 import 但不使用）
var _ = inventory.InvLocation{}
var _ = fmt.Sprintf
