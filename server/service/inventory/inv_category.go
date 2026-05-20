package inventory

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/inventory"
)

// CategoryService 商品分类服务
//
// === Go 学习要点 ===
// 1. 空结构体 struct{} —— Go 习惯用空结构体作为"方法的命名空间"，零内存开销
// 2. 方法接收者 (s *CategoryService) —— 类似 Java 的 this，但要显式声明
// 3. 命名返回值 (err error) —— 函数签名里直接命名返回值变量，函数体里可直接 return
// 4. global.GVA_DB —— GORM 全局单例，所有 service 都用这个查库
type CategoryService struct{}

// Create 创建分类
func (s *CategoryService) Create(e *inventory.InvCategory) error {
	return global.GVA_DB.Create(e).Error
}

// Update 全字段更新（GORM 的 Save = 主键存在则 update，否则 insert）
func (s *CategoryService) Update(e *inventory.InvCategory) error {
	return global.GVA_DB.Save(e).Error
}

// Delete 软删除（GVA_MODEL 含 DeletedAt，GORM 自动改 update + 设置 deleted_at）
func (s *CategoryService) Delete(id uint) error {
	return global.GVA_DB.Delete(&inventory.InvCategory{}, id).Error
}

// DeleteBatch 批量软删除
func (s *CategoryService) DeleteBatch(ids []uint) error {
	return global.GVA_DB.Delete(&[]inventory.InvCategory{}, ids).Error
}

// GetByID 根据 ID 查询单条
func (s *CategoryService) GetByID(id uint) (e inventory.InvCategory, err error) {
	err = global.GVA_DB.First(&e, id).Error
	return
}

// GetList 分页查询 + 关键字搜索
//
// 注意：interface{} 是空接口，可以承载任意类型，但调用方需要类型断言才能用
//      这里用 []InvCategory 更明确，但 GVA 习惯返回 interface{} 保持一致性
func (s *CategoryService) GetList(info request.PageInfo) (list []inventory.InvCategory, total int64, err error) {
	db := global.GVA_DB.Model(&inventory.InvCategory{})
	if info.Keyword != "" {
		db = db.Where("code LIKE ? OR name LIKE ?", "%"+info.Keyword+"%", "%"+info.Keyword+"%")
	}
	if err = db.Count(&total).Error; err != nil {
		return
	}
	err = db.Scopes(info.Paginate()).Order("sort asc, id desc").Find(&list).Error
	return
}

// GetAll 获取全部（用于下拉框、树形）
func (s *CategoryService) GetAll() (list []inventory.InvCategory, err error) {
	err = global.GVA_DB.Where("status = ?", 1).Order("sort asc, id desc").Find(&list).Error
	return
}
