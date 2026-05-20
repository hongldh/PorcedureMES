<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo">
        <el-form-item label="仓库">
          <el-select v-model="searchInfo.warehouseId" placeholder="全部仓库" clearable style="width:200px">
            <el-option v-for="w in warehouseOpts" :key="w.ID" :label="w.name" :value="w.ID" />
          </el-select>
        </el-form-item>
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
        <el-table-column label="仓库" width="160">
          <template #default="{ row }">{{ row.warehouse?.name || '-' }}</template>
        </el-table-column>
        <el-table-column prop="code" label="货位编码" width="160" />
        <el-table-column prop="name" label="货位名称" />
        <el-table-column prop="sort" label="排序" width="80" />
        <el-table-column label="状态" width="80">
          <template #default="{ row }"><el-tag :type="row.status === 1 ? 'success' : 'info'">{{ row.status === 1 ? '启用' : '停用' }}</el-tag></template>
        </el-table-column>
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
    <el-dialog v-model="dialogVisible" :title="dialogType === 'add' ? '新增货位' : '编辑货位'" width="500px">
      <el-form ref="formRef" :model="formData" :rules="rules" label-width="80px">
        <el-form-item label="仓库" prop="warehouseId">
          <el-select v-model="formData.warehouseId" style="width:100%">
            <el-option v-for="w in warehouseOpts" :key="w.ID" :label="w.name" :value="w.ID" />
          </el-select>
        </el-form-item>
        <el-form-item label="编码" prop="code"><el-input v-model="formData.code" /></el-form-item>
        <el-form-item label="名称" prop="name"><el-input v-model="formData.name" /></el-form-item>
        <el-form-item label="排序"><el-input-number v-model="formData.sort" :min="0" /></el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="formData.status"><el-radio :value="1">启用</el-radio><el-radio :value="2">停用</el-radio></el-radio-group>
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
import { createLocation, updateLocation, deleteLocation, deleteLocationBatch, getLocationList } from '@/api/inventory/warehouse'
import { getWarehouseAll } from '@/api/inventory/warehouse'

defineOptions({ name: 'InvLocation' })

const searchInfo = reactive({ keyword: '', warehouseId: undefined })
const page = ref(1); const pageSize = ref(10); const total = ref(0)
const tableData = ref([]); const selectedIds = ref([])
const warehouseOpts = ref([])
const dialogVisible = ref(false); const dialogType = ref('add'); const formRef = ref(null)
const defaultForm = () => ({ ID: 0, warehouseId: undefined, code: '', name: '', sort: 0, status: 1, remark: '' })
const formData = reactive(defaultForm())
const rules = {
  warehouseId: [{ required: true, message: '请选择仓库', trigger: 'change' }],
  code: [{ required: true, message: '必填', trigger: 'blur' }],
  name: [{ required: true, message: '必填', trigger: 'blur' }],
}
const getList = async () => {
  const res = await getLocationList({ page: page.value, pageSize: pageSize.value, ...searchInfo })
  if (res.code === 0) { tableData.value = res.data.list || []; total.value = res.data.total }
}
const onSearch = () => { page.value = 1; getList() }
const onReset = () => { searchInfo.keyword = ''; searchInfo.warehouseId = undefined; onSearch() }
const onSelectionChange = (rows) => { selectedIds.value = rows.map(r => r.ID) }
const openDialog = (type, row) => {
  dialogType.value = type; Object.assign(formData, defaultForm())
  if (type === 'edit' && row) Object.assign(formData, row)
  dialogVisible.value = true
}
const onSubmit = async () => {
  await formRef.value.validate()
  const fn = dialogType.value === 'add' ? createLocation : updateLocation
  const res = await fn(formData)
  if (res.code === 0) { ElMessage.success(res.msg); dialogVisible.value = false; getList() }
}
const onDelete = (row) => {
  ElMessageBox.confirm(`确认删除 [${row.name}] ?`, '提示', { type: 'warning' }).then(async () => {
    const res = await deleteLocation({ id: row.ID }); if (res.code === 0) { ElMessage.success('删除成功'); getList() }
  }).catch(() => {})
}
const onBatchDelete = () => {
  ElMessageBox.confirm(`删除 ${selectedIds.value.length} 项?`, '提示', { type: 'warning' }).then(async () => {
    const res = await deleteLocationBatch({ ids: selectedIds.value })
    if (res.code === 0) { ElMessage.success('成功'); getList() }
  }).catch(() => {})
}
onMounted(async () => {
  const w = await getWarehouseAll(); if (w.code === 0) warehouseOpts.value = w.data.list || []
  getList()
})
</script>
