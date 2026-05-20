<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo">
        <el-form-item label="关键字"><el-input v-model="searchInfo.keyword" placeholder="编码/名称" clearable /></el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSearch">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openDialog('add')">新增</el-button>
        <el-button type="danger" icon="delete" :disabled="!selectedIds.length" @click="onBatchDelete">批量删除</el-button>
      </div>
      <el-table :data="tableData" @selection-change="onSelectionChange" border>
        <el-table-column type="selection" width="55" />
        <el-table-column prop="code" label="编码" width="180" />
        <el-table-column prop="name" label="名称" />
        <el-table-column label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'">{{ row.status === 1 ? '启用' : '停用' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" />
        <el-table-column label="操作" width="160" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="openDialog('edit', row)">编辑</el-button>
            <el-button type="danger" link @click="onDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination v-model:current-page="page" v-model:page-size="pageSize" :total="total"
          :page-sizes="[10,20,50,100]" layout="total, sizes, prev, pager, next, jumper"
          @current-change="getList" @size-change="getList" />
      </div>
    </div>
    <el-dialog v-model="dialogVisible" :title="dialogType === 'add' ? '新增单位' : '编辑单位'" width="480px">
      <el-form ref="formRef" :model="formData" :rules="rules" label-width="80px">
        <el-form-item label="编码" prop="code"><el-input v-model="formData.code" /></el-form-item>
        <el-form-item label="名称" prop="name"><el-input v-model="formData.name" /></el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="formData.status">
            <el-radio :value="1">启用</el-radio><el-radio :value="2">停用</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="备注"><el-input v-model="formData.remark" type="textarea" /></el-form-item>
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
import { createUnit, updateUnit, deleteUnit, deleteUnitBatch, getUnitList } from '@/api/inventory/unit'

defineOptions({ name: 'InvUnit' })

const searchInfo = reactive({ keyword: '' })
const page = ref(1); const pageSize = ref(10); const total = ref(0)
const tableData = ref([]); const selectedIds = ref([])
const dialogVisible = ref(false); const dialogType = ref('add'); const formRef = ref(null)
const defaultForm = () => ({ ID: 0, code: '', name: '', status: 1, remark: '' })
const formData = reactive(defaultForm())
const rules = {
  code: [{ required: true, message: '请输入编码', trigger: 'blur' }],
  name: [{ required: true, message: '请输入名称', trigger: 'blur' }],
}
const getList = async () => {
  const res = await getUnitList({ page: page.value, pageSize: pageSize.value, keyword: searchInfo.keyword })
  if (res.code === 0) { tableData.value = res.data.list || []; total.value = res.data.total }
}
const onSearch = () => { page.value = 1; getList() }
const onReset = () => { searchInfo.keyword = ''; onSearch() }
const onSelectionChange = (rows) => { selectedIds.value = rows.map(r => r.ID) }
const openDialog = (type, row) => {
  dialogType.value = type
  Object.assign(formData, defaultForm())
  if (type === 'edit' && row) Object.assign(formData, row)
  dialogVisible.value = true
}
const onSubmit = async () => {
  await formRef.value.validate()
  const fn = dialogType.value === 'add' ? createUnit : updateUnit
  const res = await fn(formData)
  if (res.code === 0) { ElMessage.success(res.msg); dialogVisible.value = false; getList() }
}
const onDelete = (row) => {
  ElMessageBox.confirm(`确认删除 [${row.name}] ?`, '提示', { type: 'warning' }).then(async () => {
    const res = await deleteUnit({ id: row.ID }); if (res.code === 0) { ElMessage.success('删除成功'); getList() }
  }).catch(() => {})
}
const onBatchDelete = () => {
  ElMessageBox.confirm(`确认删除选中 ${selectedIds.value.length} 项?`, '提示', { type: 'warning' }).then(async () => {
    const res = await deleteUnitBatch({ ids: selectedIds.value })
    if (res.code === 0) { ElMessage.success('批量删除成功'); getList() }
  }).catch(() => {})
}
onMounted(getList)
</script>
