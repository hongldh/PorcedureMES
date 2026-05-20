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
        <el-form-item><el-button type="primary" icon="search" @click="onSearch">查询</el-button><el-button icon="refresh" @click="onReset">重置</el-button></el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list"><el-button type="primary" icon="plus" @click="openDialog('add')">新建盘点单</el-button></div>
      <el-table :data="tableData" border>
        <el-table-column prop="docNo" label="单号" width="180" />
        <el-table-column prop="bizDate" label="日期" width="120" />
        <el-table-column label="仓库" width="160"><template #default="{ row }">{{ row.warehouse?.name || '-' }}</template></el-table-column>
        <el-table-column label="状态" width="90"><template #default="{ row }"><el-tag :type="statusTag(row.status)">{{ statusText(row.status) }}</el-tag></template></el-table-column>
        <el-table-column prop="remark" label="备注" />
        <el-table-column label="操作" width="240" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="openDialog('view', row)">查看</el-button>
            <el-button v-if="row.status === 1" type="primary" link @click="openDialog('edit', row)">编辑</el-button>
            <el-button v-if="row.status === 1" type="success" link @click="onConfirm(row)">确认盘点</el-button>
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
          <el-col :span="8"><el-form-item label="单号"><el-input v-model="formData.docNo" disabled /></el-form-item></el-col>
          <el-col :span="8"><el-form-item label="单据日期" prop="bizDate"><el-date-picker v-model="formData.bizDate" type="date" value-format="YYYY-MM-DD" style="width:100%" /></el-form-item></el-col>
          <el-col :span="8"><el-form-item label="仓库" prop="warehouseId">
            <el-select v-model="formData.warehouseId" style="width:100%" @change="onWhChange">
              <el-option v-for="w in warehouseOpts" :key="w.ID" :label="w.name" :value="w.ID" />
            </el-select>
          </el-form-item></el-col>
          <el-col :span="24"><el-form-item label="备注"><el-input v-model="formData.remark" /></el-form-item></el-col>
        </el-row>
        <el-divider>盘点明细（提示：账面数在确认时由系统快照填充）</el-divider>
        <el-table :data="formData.items" border size="small">
          <el-table-column type="index" label="#" width="50" />
          <el-table-column label="商品" min-width="220">
            <template #default="{ row }">
              <el-select v-model="row.productId" filterable :disabled="readonly" style="width:100%">
                <el-option v-for="p in productOpts" :key="p.ID" :label="`${p.code} - ${p.name}`" :value="p.ID" />
              </el-select>
            </template>
          </el-table-column>
          <el-table-column label="货位" width="160">
            <template #default="{ row }"><el-select v-model="row.locationId" clearable :disabled="readonly" style="width:100%">
              <el-option v-for="l in locationOpts" :key="l.ID" :label="l.code" :value="l.ID" /></el-select></template>
          </el-table-column>
          <el-table-column label="账面数" width="110"><template #default="{ row }">{{ Number(row.bookQty || 0).toFixed(4) }}</template></el-table-column>
          <el-table-column label="实盘数" width="130">
            <template #default="{ row }"><el-input-number v-model="row.actualQty" :min="0" :precision="4" :disabled="readonly" controls-position="right" style="width:100%" @change="row.diffQty = Number(((row.actualQty||0)-(row.bookQty||0)).toFixed(4))" /></template>
          </el-table-column>
          <el-table-column label="差异" width="110">
            <template #default="{ row }">
              <span :style="{color: (row.diffQty||0)>0?'#67c23a':((row.diffQty||0)<0?'#f56c6c':'')}">{{ Number(row.diffQty || 0).toFixed(4) }}</span>
            </template>
          </el-table-column>
          <el-table-column label="备注" min-width="120"><template #default="{ row }"><el-input v-model="row.remark" :disabled="readonly" /></template></el-table-column>
          <el-table-column v-if="!readonly" label="操作" width="70"><template #default="{ $index }"><el-button type="danger" link @click="formData.items.splice($index, 1)">删除</el-button></template></el-table-column>
        </el-table>
        <el-button v-if="!readonly" type="primary" link icon="plus" @click="addRow" style="margin-top:8px">新增明细行</el-button>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">关闭</el-button>
        <el-button v-if="!readonly" type="primary" @click="onSubmit">保存草稿</el-button>
        <el-button v-if="!readonly && formData.ID" type="success" @click="onSubmitAndConfirm">保存并确认盘点</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { stocktakeApi } from '@/api/inventory/documents'
import { getProductList } from '@/api/inventory/product'
import { getWarehouseAll, getLocationByWarehouse } from '@/api/inventory/warehouse'

defineOptions({ name: 'InvStocktake' })

const searchInfo = reactive({ docNo: '', status: undefined })
const dateRange = ref([])
const page = ref(1); const pageSize = ref(10); const total = ref(0); const tableData = ref([])
const statusText = (s) => ({ 1: '草稿', 3: '已执行', 4: '作废' }[s] || '-')
const statusTag = (s) => ({ 1: 'info', 3: 'success', 4: 'danger' }[s] || 'info')
const getList = async () => {
  const res = await stocktakeApi.getList({ page: page.value, pageSize: pageSize.value, ...searchInfo })
  if (res.code === 0) { tableData.value = res.data.list || []; total.value = res.data.total }
}
const onSearch = () => { page.value = 1; getList() }
const onReset = () => { searchInfo.docNo = ''; searchInfo.status = undefined; onSearch() }

const warehouseOpts = ref([]); const productOpts = ref([]); const locationOpts = ref([])
const dialogVisible = ref(false); const dialogType = ref('add'); const formRef = ref(null)
const dialogTitle = computed(() => ({ add: '新建盘点单', edit: '编辑盘点单', view: '查看盘点单' }[dialogType.value]))
const readonly = computed(() => dialogType.value === 'view')
const newRow = () => ({ productId: undefined, locationId: undefined, bookQty: 0, actualQty: 0, diffQty: 0, remark: '' })
const defaultForm = () => ({ ID: 0, docNo: '', bizDate: new Date().toISOString().slice(0, 10), warehouseId: undefined, remark: '', items: [newRow()] })
const formData = reactive(defaultForm())
const rules = { bizDate: [{ required: true, message: '必填' }], warehouseId: [{ required: true, message: '必选' }] }
const addRow = () => formData.items.push(newRow())
const loadLocations = async (whId) => {
  if (!whId) { locationOpts.value = []; return }
  const res = await getLocationByWarehouse({ warehouseId: whId })
  if (res.code === 0) locationOpts.value = res.data.list || []
}
const onWhChange = async (whId) => {
  formData.items.forEach(r => r.locationId = undefined)
  await loadLocations(whId)
}
const openDialog = async (type, row) => {
  dialogType.value = type; Object.assign(formData, defaultForm()); formData.items = [newRow()]
  if (row) {
    const res = await stocktakeApi.findOne({ id: row.ID })
    if (res.code === 0) {
      const d = res.data.record || res.data.redata || res.data
      Object.assign(formData, { ID: d.ID, docNo: d.docNo, bizDate: d.bizDate, warehouseId: d.warehouseId, remark: d.remark, items: d.items || [newRow()] })
      if (d.warehouseId) await loadLocations(d.warehouseId)
    }
  }
  dialogVisible.value = true
}
const buildPayload = () => ({ ID: formData.ID, bizDate: formData.bizDate, warehouseId: formData.warehouseId, remark: formData.remark,
  items: formData.items.filter(r => r.productId).map(r => ({ ID: r.ID || 0, productId: r.productId, locationId: r.locationId || 0, actualQty: r.actualQty || 0, remark: r.remark })) })
const onSubmit = async () => {
  await formRef.value.validate()
  const payload = buildPayload()
  if (!payload.items.length) { ElMessage.error('请至少添加一条明细'); return false }
  const fn = formData.ID ? stocktakeApi.update : stocktakeApi.create
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
  if (ok && formData.ID) { const res = await stocktakeApi.confirm({ id: formData.ID }); if (res.code === 0) { ElMessage.success('已确认盘点'); dialogVisible.value = false; getList() } }
}
const onConfirm = (row) => { ElMessageBox.confirm(`确认盘点 [${row.docNo}]? 差异将生成调整流水。`, '提示', { type: 'warning' }).then(async () => {
  const res = await stocktakeApi.confirm({ id: row.ID }); if (res.code === 0) { ElMessage.success('成功'); getList() } }).catch(() => {}) }
const onVoid = (row) => { ElMessageBox.confirm(`作废 [${row.docNo}]?`, '提示', { type: 'warning' }).then(async () => {
  const res = await stocktakeApi.void({ id: row.ID }); if (res.code === 0) { ElMessage.success('已作废'); getList() } }).catch(() => {}) }
const onDelete = (row) => { ElMessageBox.confirm(`删除草稿 [${row.docNo}]?`, '提示', { type: 'warning' }).then(async () => {
  const res = await stocktakeApi.remove({ id: row.ID }); if (res.code === 0) { ElMessage.success('已删除'); getList() } }).catch(() => {}) }

onMounted(async () => {
  const [w, p] = await Promise.all([getWarehouseAll(), getProductList({ page: 1, pageSize: 9999 })])
  if (w.code === 0) warehouseOpts.value = w.data.list || []
  if (p.code === 0) productOpts.value = p.data.list || []
  getList()
})
</script>
