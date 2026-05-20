<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo">
        <el-form-item label="单号"><el-input v-model="searchInfo.docNo" clearable /></el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchInfo.status" placeholder="全部" clearable style="width:120px">
            <el-option label="草稿" :value="1" /><el-option label="已执行" :value="3" /><el-option label="作废" :value="4" />
          </el-select>
        </el-form-item>
        <el-form-item label="日期"><el-date-picker v-model="dateRange" type="daterange" value-format="YYYY-MM-DD" /></el-form-item>
        <el-form-item><el-button type="primary" icon="search" @click="onSearch">查询</el-button><el-button icon="refresh" @click="onReset">重置</el-button></el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list"><el-button type="primary" icon="plus" @click="openDialog('add')">新建调拨单</el-button></div>
      <el-table :data="tableData" border>
        <el-table-column prop="docNo" label="单号" width="180" />
        <el-table-column prop="bizDate" label="日期" width="120" />
        <el-table-column label="源仓" width="140"><template #default="{ row }">{{ row.fromWh?.name || '-' }}</template></el-table-column>
        <el-table-column label="目标仓" width="140"><template #default="{ row }">{{ row.toWh?.name || '-' }}</template></el-table-column>
        <el-table-column label="状态" width="90"><template #default="{ row }"><el-tag :type="statusTag(row.status)">{{ statusText(row.status) }}</el-tag></template></el-table-column>
        <el-table-column prop="remark" label="备注" />
        <el-table-column label="操作" width="240" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="openDialog('view', row)">查看</el-button>
            <el-button v-if="row.status === 1" type="primary" link @click="openDialog('edit', row)">编辑</el-button>
            <el-button v-if="row.status === 1" type="success" link @click="onConfirm(row)">确认调拨</el-button>
            <el-button v-if="row.status === 1" type="danger" link @click="onDelete(row)">删除</el-button>
            <el-button v-if="row.status === 3" type="warning" link @click="onVoid(row)">作废</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination v-model:current-page="page" v-model:page-size="pageSize" :total="total" :page-sizes="[10,20,50,100]"
          layout="total, sizes, prev, pager, next, jumper" @current-change="getList" @size-change="getList" />
      </div>
    </div>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="1000px" :close-on-click-modal="false">
      <el-form ref="formRef" :model="formData" :rules="rules" label-width="90px" :disabled="readonly">
        <el-row :gutter="16">
          <el-col :span="6"><el-form-item label="单号"><el-input v-model="formData.docNo" disabled placeholder="保存后生成" /></el-form-item></el-col>
          <el-col :span="6"><el-form-item label="单据日期" prop="bizDate"><el-date-picker v-model="formData.bizDate" type="date" value-format="YYYY-MM-DD" style="width:100%" /></el-form-item></el-col>
          <el-col :span="6"><el-form-item label="源仓库" prop="fromWarehouse">
            <el-select v-model="formData.fromWarehouse" style="width:100%" @change="(v) => onWhChange('from', v)">
              <el-option v-for="w in warehouseOpts" :key="w.ID" :label="w.name" :value="w.ID" />
            </el-select>
          </el-form-item></el-col>
          <el-col :span="6"><el-form-item label="目标仓库" prop="toWarehouse">
            <el-select v-model="formData.toWarehouse" style="width:100%" @change="(v) => onWhChange('to', v)">
              <el-option v-for="w in warehouseOpts" :key="w.ID" :label="w.name" :value="w.ID" />
            </el-select>
          </el-form-item></el-col>
          <el-col :span="24"><el-form-item label="备注"><el-input v-model="formData.remark" /></el-form-item></el-col>
        </el-row>
        <el-divider>明细</el-divider>
        <el-table :data="formData.items" border size="small">
          <el-table-column type="index" label="#" width="50" />
          <el-table-column label="商品" min-width="220">
            <template #default="{ row }">
              <el-select v-model="row.productId" filterable :disabled="readonly" style="width:100%">
                <el-option v-for="p in productOpts" :key="p.ID" :label="`${p.code} - ${p.name}`" :value="p.ID" />
              </el-select>
            </template>
          </el-table-column>
          <el-table-column label="源货位" width="160">
            <template #default="{ row }"><el-select v-model="row.fromLocationId" clearable :disabled="readonly" style="width:100%">
              <el-option v-for="l in fromLocOpts" :key="l.ID" :label="l.code" :value="l.ID" /></el-select></template>
          </el-table-column>
          <el-table-column label="目标货位" width="160">
            <template #default="{ row }"><el-select v-model="row.toLocationId" clearable :disabled="readonly" style="width:100%">
              <el-option v-for="l in toLocOpts" :key="l.ID" :label="l.code" :value="l.ID" /></el-select></template>
          </el-table-column>
          <el-table-column label="数量" width="120"><template #default="{ row }"><el-input-number v-model="row.quantity" :min="0" :precision="4" :disabled="readonly" controls-position="right" style="width:100%" /></template></el-table-column>
          <el-table-column label="备注" min-width="120"><template #default="{ row }"><el-input v-model="row.remark" :disabled="readonly" /></template></el-table-column>
          <el-table-column v-if="!readonly" label="操作" width="70"><template #default="{ $index }"><el-button type="danger" link @click="formData.items.splice($index, 1)">删除</el-button></template></el-table-column>
        </el-table>
        <el-button v-if="!readonly" type="primary" link icon="plus" @click="addRow" style="margin-top:8px">新增明细行</el-button>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">关闭</el-button>
        <el-button v-if="!readonly" type="primary" @click="onSubmit">保存草稿</el-button>
        <el-button v-if="!readonly && formData.ID" type="success" @click="onSubmitAndConfirm">保存并确认调拨</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { transferApi } from '@/api/inventory/documents'
import { getProductList } from '@/api/inventory/product'
import { getWarehouseAll, getLocationByWarehouse } from '@/api/inventory/warehouse'

defineOptions({ name: 'InvTransfer' })

const searchInfo = reactive({ docNo: '', status: undefined })
const dateRange = ref([])
const page = ref(1); const pageSize = ref(10); const total = ref(0); const tableData = ref([])
const statusText = (s) => ({ 1: '草稿', 3: '已执行', 4: '作废' }[s] || '-')
const statusTag = (s) => ({ 1: 'info', 3: 'success', 4: 'danger' }[s] || 'info')
const getList = async () => {
  const params = { page: page.value, pageSize: pageSize.value, ...searchInfo }
  if (dateRange.value?.length === 2) { params.dateStart = dateRange.value[0]; params.dateEnd = dateRange.value[1] }
  const res = await transferApi.getList(params)
  if (res.code === 0) { tableData.value = res.data.list || []; total.value = res.data.total }
}
const onSearch = () => { page.value = 1; getList() }
const onReset = () => { searchInfo.docNo = ''; searchInfo.status = undefined; dateRange.value = []; onSearch() }

const warehouseOpts = ref([]); const productOpts = ref([])
const fromLocOpts = ref([]); const toLocOpts = ref([])

const dialogVisible = ref(false); const dialogType = ref('add'); const formRef = ref(null)
const dialogTitle = computed(() => ({ add: '新建调拨单', edit: '编辑调拨单', view: '查看调拨单' }[dialogType.value]))
const readonly = computed(() => dialogType.value === 'view')

const newRow = () => ({ productId: undefined, fromLocationId: undefined, toLocationId: undefined, quantity: 0, remark: '' })
const defaultForm = () => ({ ID: 0, docNo: '', bizDate: new Date().toISOString().slice(0, 10), fromWarehouse: undefined, toWarehouse: undefined, remark: '', items: [newRow()] })
const formData = reactive(defaultForm())
const rules = {
  bizDate: [{ required: true, message: '必填' }],
  fromWarehouse: [{ required: true, message: '必选' }],
  toWarehouse: [{ required: true, message: '必选' }]
}
const addRow = () => formData.items.push(newRow())
const loadLocations = async (which, whId) => {
  const opts = which === 'from' ? fromLocOpts : toLocOpts
  if (!whId) { opts.value = []; return }
  const res = await getLocationByWarehouse({ warehouseId: whId })
  if (res.code === 0) opts.value = res.data.list || []
}
const onWhChange = async (which, whId) => {
  formData.items.forEach(r => { if (which === 'from') r.fromLocationId = undefined; else r.toLocationId = undefined })
  await loadLocations(which, whId)
}

const openDialog = async (type, row) => {
  dialogType.value = type; Object.assign(formData, defaultForm()); formData.items = [newRow()]
  if (row) {
    const res = await transferApi.findOne({ id: row.ID })
    if (res.code === 0) {
      const d = res.data.record || res.data.redata || res.data
      Object.assign(formData, { ID: d.ID, docNo: d.docNo, bizDate: d.bizDate, fromWarehouse: d.fromWarehouse, toWarehouse: d.toWarehouse, remark: d.remark, items: d.items || [newRow()] })
      if (d.fromWarehouse) await loadLocations('from', d.fromWarehouse)
      if (d.toWarehouse) await loadLocations('to', d.toWarehouse)
    }
  }
  dialogVisible.value = true
}
const buildPayload = () => ({ ID: formData.ID, bizDate: formData.bizDate, fromWarehouse: formData.fromWarehouse, toWarehouse: formData.toWarehouse, remark: formData.remark,
  items: formData.items.filter(r => r.productId && r.quantity > 0).map(r => ({ ID: r.ID || 0, productId: r.productId, fromLocationId: r.fromLocationId || 0, toLocationId: r.toLocationId || 0, quantity: r.quantity, remark: r.remark })) })
const onSubmit = async () => {
  await formRef.value.validate()
  const payload = buildPayload()
  if (!payload.items.length) { ElMessage.error('请至少添加一条明细'); return false }
  if (payload.fromWarehouse === payload.toWarehouse) { ElMessage.error('源仓与目标仓不可相同'); return false }
  const fn = formData.ID ? transferApi.update : transferApi.create
  const res = await fn(payload)
  if (res.code === 0) {
    ElMessage.success('保存成功')
    const newId = res.data?.ID || res.data?.id || res.data?.redata?.ID
    if (newId && !formData.ID) formData.ID = newId
    getList(); return true
  }
  return false
}
const onSubmitAndConfirm = async () => {
  const ok = await onSubmit()
  if (ok && formData.ID) { const res = await transferApi.confirm({ id: formData.ID }); if (res.code === 0) { ElMessage.success('已确认调拨'); dialogVisible.value = false; getList() } }
}
const onConfirm = (row) => { ElMessageBox.confirm(`确认调拨 [${row.docNo}]?`, '提示', { type: 'warning' }).then(async () => {
  const res = await transferApi.confirm({ id: row.ID }); if (res.code === 0) { ElMessage.success('成功'); getList() } }).catch(() => {}) }
const onVoid = (row) => { ElMessageBox.confirm(`作废 [${row.docNo}]?`, '提示', { type: 'warning' }).then(async () => {
  const res = await transferApi.void({ id: row.ID }); if (res.code === 0) { ElMessage.success('已作废'); getList() } }).catch(() => {}) }
const onDelete = (row) => { ElMessageBox.confirm(`删除草稿 [${row.docNo}]?`, '提示', { type: 'warning' }).then(async () => {
  const res = await transferApi.remove({ id: row.ID }); if (res.code === 0) { ElMessage.success('已删除'); getList() } }).catch(() => {}) }

onMounted(async () => {
  const [w, p] = await Promise.all([getWarehouseAll(), getProductList({ page: 1, pageSize: 9999 })])
  if (w.code === 0) warehouseOpts.value = w.data.list || []
  if (p.code === 0) productOpts.value = p.data.list || []
  getList()
})
</script>
