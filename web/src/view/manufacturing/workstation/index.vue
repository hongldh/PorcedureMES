<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo">
        <el-form-item label="关键字"><el-input v-model="searchInfo.keyword" placeholder="编码/名称" clearable /></el-form-item>
        <el-form-item label="类型">
          <el-select v-model="searchInfo.type" clearable placeholder="全部" style="width:120px">
            <el-option label="人工" :value="1" /><el-option label="设备" :value="2" /><el-option label="外协" :value="3" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchInfo.status" clearable placeholder="全部" style="width:120px">
            <el-option label="启用" :value="1" /><el-option label="停用" :value="2" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSearch">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="gva-table-box">
      <div class="gva-btn-list"><el-button type="primary" icon="plus" @click="openDialog('add')">新增工位</el-button></div>
      <el-table :data="tableData" border>
        <el-table-column prop="code" label="编码" width="180" />
        <el-table-column prop="name" label="名称" />
        <el-table-column label="类型" width="100">
          <template #default="{ row }">
            <el-tag :type="typeTag(row.type)">{{ typeText(row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="capacity" label="标称产能(/h)" width="140" />
        <el-table-column label="所属仓库" width="160"><template #default="{ row }">{{ row.warehouse?.name || '-' }}</template></el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }"><el-tag :type="row.status === 1 ? 'success' : 'info'">{{ row.status === 1 ? '启用' : '停用' }}</el-tag></template>
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
        <el-pagination v-model:current-page="page" v-model:page-size="pageSize" :total="total" :page-sizes="[10,20,50,100]"
          layout="total, sizes, prev, pager, next, jumper" @current-change="getList" @size-change="getList" />
      </div>
    </div>

    <el-dialog v-model="dialogVisible" :title="dialogType === 'add' ? '新增工位' : '编辑工位'" width="560px">
      <el-form ref="formRef" :model="formData" :rules="rules" label-width="100px">
        <el-form-item label="编码" prop="code"><el-input v-model="formData.code" placeholder="如 WS001" /></el-form-item>
        <el-form-item label="名称" prop="name"><el-input v-model="formData.name" /></el-form-item>
        <el-form-item label="类型" prop="type">
          <el-select v-model="formData.type" style="width:100%">
            <el-option label="人工" :value="1" /><el-option label="设备" :value="2" /><el-option label="外协" :value="3" />
          </el-select>
        </el-form-item>
        <el-form-item label="标称产能"><el-input-number v-model="formData.capacity" :min="0" :precision="4" controls-position="right" style="width:100%" /></el-form-item>
        <el-form-item label="所属仓库">
          <el-select v-model="formData.warehouseId" clearable style="width:100%">
            <el-option v-for="w in warehouseOpts" :key="w.ID" :label="w.name" :value="w.ID" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="formData.status">
            <el-radio :value="1">启用</el-radio><el-radio :value="2">停用</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="备注"><el-input v-model="formData.remark" type="textarea" :rows="2" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">关闭</el-button>
        <el-button type="primary" @click="onSubmit">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { workstationApi } from '@/api/manufacturing'
import { getWarehouseAll } from '@/api/inventory/warehouse'

defineOptions({ name: 'MfgWorkstation' })

const typeText = (t) => ({ 1: '人工', 2: '设备', 3: '外协' }[t] || '-')
const typeTag  = (t) => ({ 1: '', 2: 'success', 3: 'warning' }[t] || '')

const searchInfo = reactive({ keyword: '', type: undefined, status: undefined })
const page = ref(1); const pageSize = ref(10); const total = ref(0); const tableData = ref([])
const warehouseOpts = ref([])

const getList = async () => {
  const res = await workstationApi.getList({ page: page.value, pageSize: pageSize.value, ...searchInfo })
  if (res.code === 0) { tableData.value = res.data.list || []; total.value = res.data.total }
}
const onSearch = () => { page.value = 1; getList() }
const onReset = () => { Object.assign(searchInfo, { keyword: '', type: undefined, status: undefined }); onSearch() }

const dialogVisible = ref(false); const dialogType = ref('add'); const formRef = ref(null)
const defaultForm = () => ({ ID: 0, code: '', name: '', type: 1, capacity: 0, warehouseId: undefined, status: 1, remark: '' })
const formData = reactive(defaultForm())
const rules = {
  code: [{ required: true, message: '必填' }],
  name: [{ required: true, message: '必填' }],
  type: [{ required: true, message: '必选' }]
}

const openDialog = (type, row) => {
  dialogType.value = type
  Object.assign(formData, defaultForm())
  if (row) Object.assign(formData, row)
  dialogVisible.value = true
}
const onSubmit = async () => {
  try { await formRef.value.validate() } catch { ElMessage.error('请补全必填项'); return }
  const fn = formData.ID ? workstationApi.update : workstationApi.create
  const res = await fn(formData)
  if (res.code === 0) { ElMessage.success('保存成功'); dialogVisible.value = false; getList() }
}
const onDelete = (row) => {
  ElMessageBox.confirm(`删除工位 [${row.code}]?`, '提示', { type: 'warning' }).then(async () => {
    const res = await workstationApi.remove({ id: row.ID })
    if (res.code === 0) { ElMessage.success('已删除'); getList() }
  }).catch(() => {})
}

onMounted(async () => {
  const w = await getWarehouseAll()
  if (w.code === 0) warehouseOpts.value = w.data.list || []
  getList()
})
</script>
