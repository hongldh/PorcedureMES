<!--
  MES - 完工入库单
-->
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
      <div class="gva-btn-list"><el-button type="primary" icon="plus" @click="openDialog('add')">新建完工入库单</el-button></div>
      <el-table :data="tableData" border>
        <el-table-column prop="docNo" label="单号" width="180" />
        <el-table-column prop="bizDate" label="日期" width="120" />
        <el-table-column label="工单" width="180"><template #default="{ row }">{{ row.workOrder?.docNo || '-' }}</template></el-table-column>
        <el-table-column label="成品" min-width="200"><template #default="{ row }">{{ row.workOrder?.product?.code }} - {{ row.workOrder?.product?.name }}</template></el-table-column>
        <el-table-column label="入库仓" width="140"><template #default="{ row }">{{ row.warehouse?.name || '-' }}</template></el-table-column>
        <el-table-column prop="quantity" label="完工数" width="100" />
        <el-table-column prop="unitCost" label="单位成本" width="120" />
        <el-table-column prop="totalAmount" label="总金额" width="120" />
        <el-table-column label="状态" width="90"><template #default="{ row }"><el-tag :type="statusTag(row.status)">{{ statusText(row.status) }}</el-tag></template></el-table-column>
        <el-table-column label="操作" width="240" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="openDialog('view', row)">查看</el-button>
            <el-button v-if="row.status === 1" type="primary" link @click="openDialog('edit', row)">编辑</el-button>
            <el-button v-if="row.status === 1" type="success" link @click="onConfirm(row)">确认入库</el-button>
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

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="720px" :close-on-click-modal="false">
      <el-form ref="formRef" :model="formData" :rules="rules" label-width="100px" :disabled="readonly">
        <el-row :gutter="16">
          <el-col :span="12"><el-form-item label="单号"><el-input v-model="formData.docNo" disabled placeholder="保存后生成" /></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="日期" prop="bizDate"><el-date-picker v-model="formData.bizDate" type="date" value-format="YYYY-MM-DD" style="width:100%" /></el-form-item></el-col>
          <el-col :span="24"><el-form-item label="工单" prop="workOrderId">
            <el-select v-model="formData.workOrderId" filterable :disabled="!!formData.ID" style="width:100%" @change="onWoChange">
              <el-option v-for="w in woOpts" :key="w.ID"
                :label="`${w.docNo} - ${w.product?.code} 计划${w.planQty} 已完${w.actualQty}`"
                :value="w.ID" />
            </el-select>
          </el-form-item></el-col>
          <el-col :span="12"><el-form-item label="入库仓" prop="warehouseId">
            <el-select v-model="formData.warehouseId" style="width:100%" @change="onWhChange">
              <el-option v-for="w in warehouseOpts" :key="w.ID" :label="w.name" :value="w.ID" />
            </el-select>
          </el-form-item></el-col>
          <el-col :span="12"><el-form-item label="入库货位">
            <el-select v-model="formData.locationId" clearable style="width:100%">
              <el-option v-for="l in locationOpts" :key="l.ID" :label="l.code" :value="l.ID" />
            </el-select>
          </el-form-item></el-col>
          <el-col :span="12"><el-form-item label="完工数量" prop="quantity">
            <el-input-number v-model="formData.quantity" :min="0.0001" :precision="4" controls-position="right" style="width:100%" />
          </el-form-item></el-col>
          <el-col :span="12" v-if="selectedWo">
            <el-form-item label="可完工">
              <el-input :model-value="Number((selectedWo.planQty - selectedWo.actualQty) || 0).toFixed(4)" disabled />
            </el-form-item>
          </el-col>
          <el-col :span="24"><el-form-item label="备注"><el-input v-model="formData.remark" /></el-form-item></el-col>
        </el-row>
        <el-alert v-if="!readonly" type="info" :closable="false" style="margin-top:8px">
          单位成本将在"确认入库"时自动计算：累计领料金额 ÷ 完工总数量
        </el-alert>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">关闭</el-button>
        <el-button v-if="!readonly" type="primary" @click="onSubmit">保存草稿</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { receiptApi, workOrderApi } from '@/api/manufacturing'
import { getWarehouseAll, getLocationByWarehouse } from '@/api/inventory/warehouse'

defineOptions({ name: 'MfgReceipt' })

const statusText = (s) => ({ 1: '草稿', 3: '已执行', 4: '作废' }[s] || '-')
const statusTag = (s) => ({ 1: 'info', 3: 'success', 4: 'danger' }[s] || 'info')

const searchInfo = reactive({ docNo: '', status: undefined })
const page = ref(1); const pageSize = ref(10); const total = ref(0); const tableData = ref([])
const warehouseOpts = ref([]); const locationOpts = ref([])
const woOpts = ref([]); const selectedWo = ref(null)

const getList = async () => {
  const res = await receiptApi.getList({ page: page.value, pageSize: pageSize.value, ...searchInfo })
  if (res.code === 0) { tableData.value = res.data.list || []; total.value = res.data.total }
}
const onSearch = () => { page.value = 1; getList() }
const onReset = () => { Object.assign(searchInfo, { docNo: '', status: undefined }); onSearch() }

const dialogVisible = ref(false); const dialogType = ref('add'); const formRef = ref(null)
const dialogTitle = computed(() => ({ add: '新建完工入库单', edit: '编辑完工入库单', view: '查看完工入库单' }[dialogType.value]))
const readonly = computed(() => dialogType.value === 'view')
const defaultForm = () => ({ ID: 0, docNo: '', bizDate: new Date().toISOString().slice(0, 10),
  workOrderId: undefined, warehouseId: undefined, locationId: undefined, quantity: 1, remark: '' })
const formData = reactive(defaultForm())
const rules = {
  bizDate: [{ required: true, message: '必填' }],
  workOrderId: [{ required: true, message: '必选' }],
  warehouseId: [{ required: true, message: '必选' }],
  quantity: [{ required: true, message: '必填' }]
}

const loadLocations = async (whId) => {
  if (!whId) { locationOpts.value = []; return }
  const res = await getLocationByWarehouse({ warehouseId: whId })
  if (res.code === 0) locationOpts.value = res.data.list || []
}
const onWhChange = async (whId) => { formData.locationId = undefined; await loadLocations(whId) }
const onWoChange = (woId) => {
  selectedWo.value = woOpts.value.find(w => w.ID === woId) || null
  if (selectedWo.value && !formData.warehouseId && selectedWo.value.warehouseId) {
    formData.warehouseId = selectedWo.value.warehouseId
    onWhChange(formData.warehouseId)
  }
  if (selectedWo.value) {
    const remain = selectedWo.value.planQty - selectedWo.value.actualQty
    if (remain > 0 && (!formData.quantity || formData.quantity > remain)) {
      formData.quantity = Number(remain.toFixed(4))
    }
  }
}

const openDialog = async (type, row) => {
  dialogType.value = type
  Object.assign(formData, defaultForm())
  selectedWo.value = null; locationOpts.value = []
  // 加载可入库工单（已下达 + 生产中）
  const r1 = await workOrderApi.getList({ page: 1, pageSize: 9999, status: 2 })
  const r2 = await workOrderApi.getList({ page: 1, pageSize: 9999, status: 3 })
  woOpts.value = [...(r1.code === 0 ? r1.data.list || [] : []), ...(r2.code === 0 ? r2.data.list || [] : [])]

  if (row) {
    const res = await receiptApi.findOne({ id: row.ID })
    if (res.code === 0) {
      const d = res.data.record || res.data
      Object.assign(formData, { ID: d.ID, docNo: d.docNo, bizDate: d.bizDate,
        workOrderId: d.workOrderId, warehouseId: d.warehouseId, locationId: d.locationId,
        quantity: d.quantity, remark: d.remark })
      if (d.warehouseId) await loadLocations(d.warehouseId)
      if (d.workOrderId) {
        if (!woOpts.value.find(x => x.ID === d.workOrderId) && d.workOrder) {
          woOpts.value.push(d.workOrder)
        }
        selectedWo.value = woOpts.value.find(w => w.ID === d.workOrderId) || null
      }
    }
  }
  dialogVisible.value = true
}
const onSubmit = async () => {
  try { await formRef.value.validate() } catch { ElMessage.error('请补全必填项'); return }
  const fn = formData.ID ? receiptApi.update : receiptApi.create
  const res = await fn(formData)
  if (res.code === 0) { ElMessage.success('保存成功'); dialogVisible.value = false; getList() }
}
const onConfirm = (row) => {
  ElMessageBox.confirm(`确认入库 [${row.docNo}]? 将按累计领料金额自动计算单位成本。`, '提示', { type: 'warning' }).then(async () => {
    const res = await receiptApi.confirm({ id: row.ID })
    if (res.code === 0) { ElMessage.success('已确认'); getList() }
  }).catch(() => {})
}
const onVoid = (row) => {
  ElMessageBox.confirm(`作废 [${row.docNo}]? 将反向出库成品并扣减工单完工数。`, '提示', { type: 'warning' }).then(async () => {
    const res = await receiptApi.void({ id: row.ID })
    if (res.code === 0) { ElMessage.success('已作废'); getList() }
  }).catch(() => {})
}
const onDelete = (row) => {
  ElMessageBox.confirm(`删除草稿 [${row.docNo}]?`, '提示', { type: 'warning' }).then(async () => {
    const res = await receiptApi.remove({ id: row.ID })
    if (res.code === 0) { ElMessage.success('已删除'); getList() }
  }).catch(() => {})
}

onMounted(async () => {
  const w = await getWarehouseAll()
  if (w.code === 0) warehouseOpts.value = w.data.list || []
  getList()
})
</script>
