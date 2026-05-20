package inventory

import (
	"fmt"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/model/inventory"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// SequenceService 单号生成器
//
// 用法：在事务里调用 NextDocNo(tx, BizPurchaseIn) 拿到形如 "CGRK20260425001" 的单号
//
// === Go 学习要点 ===
// 1. 用 map[string]string 做"业务类型 → 单号前缀"配置表，比 switch 更易扩展
// 2. clause.Locking{Strength:"UPDATE"} 触发 SELECT ... FOR UPDATE 行锁，并发安全
// 3. 必须在调用方传入的 tx 里执行，跟主业务一个事务，保证回滚一致
type SequenceService struct{}

// 业务类型对应的单号前缀
var bizPrefix = map[string]string{
	inventory.BizPurchaseIn:      "CGRK", // 采购入库
	inventory.BizPurchaseReturn:  "CGTH", // 采购退货
	inventory.BizSalesOut:        "XSCK", // 销售出库
	inventory.BizSalesReturn:     "XSTH", // 销售退货
	inventory.BizTransferOut:     "DBDD", // 调拨（出/入共用）
	inventory.BizTransferIn:      "DBDD",
	inventory.BizStocktakeAdjust: "PDDD", // 盘点
}

// NextDocNo 取下一个单号（必须在事务中调用）
func (s *SequenceService) NextDocNo(tx *gorm.DB, bizType string) (string, error) {
	prefix, ok := bizPrefix[bizType]
	if !ok {
		return "", fmt.Errorf("未知业务类型: %s", bizType)
	}
	dateStr := time.Now().Format("20060102")

	var seq inventory.InvDocSequence
	// SELECT ... FOR UPDATE，没有就插入一条
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("biz_type = ? AND date_str = ?", bizType, dateStr).
		First(&seq).Error
	if err == gorm.ErrRecordNotFound {
		seq = inventory.InvDocSequence{BizType: bizType, DateStr: dateStr, Seq: 0}
		if err = tx.Create(&seq).Error; err != nil {
			return "", err
		}
	} else if err != nil {
		return "", err
	}
	seq.Seq++
	if err = tx.Save(&seq).Error; err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%s%03d", prefix, dateStr, seq.Seq), nil
}
