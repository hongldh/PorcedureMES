<!--
  进销存 - 商品分类管理页
  ====================================
  这是阶段 1 的"模板页面"，加详细注释，其余 6 个页面参考此结构。

  GVA 前端页面通用套路（重点学习）：
  1. 顶部 search 区：el-form 收集查询条件 + 搜索/重置按钮
  2. 中间 toolbar：新增 + 批量删除按钮
  3. 主体 el-table：列展示 + 操作列
  4. 底部 el-pagination：分页
  5. 弹窗 el-dialog：新增/编辑表单（用同一个 dialog，根据 dialogType 区分）
  6. 调用 @/api/inventory/xxx 的接口
-->
<template>
  <div>
    <!-- 搜索区 -->
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo">
        <el-form-item label="关键字">
          <el-input v-model="searchInfo.keyword" placeholder="编码/名称" clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSearch">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 表格区 -->
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openDialog('add')">新增</el-button>
        <el-button type="danger" icon="delete" :disabled="!selectedIds.length" @click="onBatchDelete">批量删除</el-button>
      </div>
      <el-table :data="tableData" @selection-change="onSelectionChange" border>
        <el-table-column type="selection" width="55" />
        <el-table-column prop="code" label="编码" width="180" />
        <el-table-column prop="name" label="名称" />
        <el-table-column prop="parentId" label="父级ID" width="100" />
        <el-table-column prop="sort" label="排序" width="80" />
        <el-table-column label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'">
              {{ row.status === 1 ? '启用' : '停用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" />
        <el-table-column prop="CreatedAt" label="创建时间" width="180">
          <template #default="{ row }">{{ formatDate(row.CreatedAt) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="160" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="openDialog('edit', row)">编辑</el-button>
            <el-button type="danger" link @click="onDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="getList"
          @size-change="getList"
        />
      </div>
    </div>

    <!-- 编辑弹窗 -->
    <el-dialog v-model="dialogVisible" :title="dialogType === 'add' ? '新增分类' : '编辑分类'" width="520px">
      <el-form ref="formRef" :model="formData" :rules="rules" label-width="100px">
        <el-form-item label="编码" prop="code">
          <el-input v-model="formData.code" placeholder="唯一编码" />
        </el-form-item>
        <el-form-item label="名称" prop="name">
          <el-input v-model="formData.name" />
        </el-form-item>
        <el-form-item label="父级ID">
          <el-input-number v-model="formData.parentId" :min="0" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="formData.sort" :min="0" />
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="formData.status">
            <el-radio :value="1">启用</el-radio>
            <el-radio :value="2">停用</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="formData.remark" type="textarea" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="onSubmit">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  createCategory, updateCategory, deleteCategory, deleteCategoryBatch, getCategoryList
} from '@/api/inventory/category'
import { formatDate } from '@/utils/format'

defineOptions({ name: 'InvCategory' })

// 搜索 & 分页
const searchInfo = reactive({ keyword: '' })
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const tableData = ref([])
const selectedIds = ref([])

// 弹窗
const dialogVisible = ref(false)
const dialogType = ref('add')
const formRef = ref(null)
const defaultForm = () => ({ ID: 0, code: '', name: '', parentId: 0, sort: 0, status: 1, remark: '' })
const formData = reactive(defaultForm())
const rules = {
  code: [{ required: true, message: '请输入编码', trigger: 'blur' }],
  name: [{ required: true, message: '请输入名称', trigger: 'blur' }],
}

// === 列表 ===
const getList = async () => {
  const res = await getCategoryList({ page: page.value, pageSize: pageSize.value, keyword: searchInfo.keyword })
  if (res.code === 0) {
    tableData.value = res.data.list || []
    total.value = res.data.total
  }
}
const onSearch = () => { page.value = 1; getList() }
const onReset = () => { searchInfo.keyword = ''; onSearch() }
const onSelectionChange = (rows) => { selectedIds.value = rows.map(r => r.ID) }

// === 弹窗 ===
const openDialog = (type, row) => {
  dialogType.value = type
  Object.assign(formData, defaultForm())
  if (type === 'edit' && row) Object.assign(formData, row)
  dialogVisible.value = true
}
const onSubmit = async () => {
  await formRef.value.validate()
  const fn = dialogType.value === 'add' ? createCategory : updateCategory
  const res = await fn(formData)
  if (res.code === 0) {
    ElMessage.success(res.msg || '操作成功')
    dialogVisible.value = false
    getList()
  }
}

// === 删除 ===
const onDelete = (row) => {
  ElMessageBox.confirm(`确认删除分类 [${row.name}] ?`, '提示', { type: 'warning' }).then(async () => {
    const res = await deleteCategory({ id: row.ID })
    if (res.code === 0) { ElMessage.success('删除成功'); getList() }
  }).catch(() => {})
}
const onBatchDelete = () => {
  ElMessageBox.confirm(`确认删除选中的 ${selectedIds.value.length} 项?`, '提示', { type: 'warning' }).then(async () => {
    const res = await deleteCategoryBatch({ ids: selectedIds.value })
    if (res.code === 0) { ElMessage.success('批量删除成功'); getList() }
  }).catch(() => {})
}

onMounted(getList)
</script>
