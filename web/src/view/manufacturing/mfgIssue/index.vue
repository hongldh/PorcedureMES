<!--
  MES - 领料单
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
      <div class="gva-btn-list"><el-button type="primary" icon="plus" @click="openDialog('add')">新建领料单</el-button></div>
      <el-table :data="tableData" border>
        <el-table-column prop="docNo" label="单号" width="180" />
        <el-table-column prop="bizDate" label="日期" width="120" />
        <el-table-column label="工单" width="180"><template #default="{ row }">{{ row.workOrder?.docNo || '-' }}</template></el-table-column>
        <el-table-column label="领料仓" width="140"><template #default="{ row }">{{ row.warehouse?.name || '-' }}</template></el-table-column>
        <el-table-column prop="totalAmount" label="总金额" width="120" />
        <el-table-column label="状态" width="90"><template #default="{ row }"><el-tag :type="statusTag(row.status)">{{ statusText(row.status) }}</el-tag></template></el-table-column>
        <el-table-column prop="remark" label="备注" />
        <el-table-column label="操作" width="240" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="openDialog('view', row)">查看</el-button>
            <el-button v-if="row.status === 1" type="primary" link @click="openDialog('edit', row)">编辑</el-button>
            <el-button v-if="row.status === 1" type="success" link @click="onConfirm(row)">确认领料</el-button>
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
      <el-form ref="formRef" :model="formData" :rules="rules" label-width="100px" :disabled="readonly">
        <el-row :gutter="16">
          <el-col :span="6"><el-form-item label="单号"><el-input v-model="formData.docNo" disabled placeholder="保存后生成" /></el-form-item></el-col>
          <el-col :span="6"><el-form-item label="日期" prop="bizDate"><el-date-picker v-model="formData.bizDate" type="date" value-format="YYYY-MM-DD" style="width:100%" /></el-form-item></el-col>
          <el-col :span="6"><el-form-item label="工单" prop="workOrderId">
            <el-select v-model="formData.workOrderId" filterable :disabled="!!formData.ID" style="width:100%" @change="onWoChange">
              <el-option v-for="w in woOpts" :key="w.ID" :label="`${w.docNo} (${w.product?.code})`" :value="w.ID" />
            </el-select>
          </el-form-item></el-col>
          <el-col :span="6"><el-form-item label="领料仓" prop="warehouseId">
            <el-select v-model="formData.warehouseId" style="width:100%" @change="onWhChange">
              <el-option v-for="w in warehouseOpts" :key="w.ID" :label="w.name" :value="w.ID" />
            </el-select>
          </el-form-item></el-col>
          <el-col :span="24"><el-form-item label="备注"><el-input v-model="formData.remark" /></el-form-item></el-col>
        </el-row>
        <el-divider>领料明细</el-divider>
        <div v-if="woRequirement.length" style="margin-bottom:8px;color:#909399;font-size:12px">
          工单需求参考（点击 + 加入领料）：
          <el-button v-for="r in woRequirement" :key="r.productId" link
            :disabled="readonly || (r.requiredQty - r.issuedQty) <= 0"
            @click="addRow(r)">
            {{ r.product?.code }} (剩{{ Number(r.requiredQty - r.issuedQty).toFixed(2) }})
          </el-button>
        </div>
        <el-table :data="formData.items" border size="small">
          <el-table-column type="index" label="#" width="50" />
          <el-table-column label="物料" min-width="220">
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
          <el-table-column label="领料数量" width="140">
            <template #default="{ row }"><el-input-number v-model="row.quantity" :min="0" :precision="4" :disabled="readonly" controls-position="right" style="width:100%" /></template>
          </el-table-column>
          <el-table-column label="出库均价" width="120"><template #default="{ row }">{{ Number(row.unitCost || 0).toFixed(4) }}</template></el-table-column>
          <el-table-column label="金额" width="120"><template #default="{ row }">{{ Number(row.amount || 0).toFixed(2) }}</template></el-table-column>
          <el-table-column label="备注" min-width="120"><template #default="{ row }"><el-input v-model="row.remark" :disabled="readonly" /></template></el-table-column>
          <el-table-column v-if="!readonly" label="操作" width="70"><template #default="{ $index }"><el-button type="danger" link @click="formData.items.splice($index, 1)">删除</el-button></template></el-table-column>
        </el-table>
        <el-button v-if="!readonly" type="primary" link icon="plus" @click="addRow()" style="margin-top:8px">新增明细</el-button>
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
import { issueApi, workOrderApi } from '@/api/manufacturing'
import { getProductList } from '@/api/inventory/product'
import { getWarehouseAll, getLocationByWarehouse } from '@/api/inventory/warehouse'

defineOptions({ name: 'MfgIssue' })

const statusText = (s) => ({ 1: '草稿', 3: '已执行', 4: '作废' }[s] || '-')
const statusTag = (s) => ({ 1: 'info', 3: 'success', 4: 'danger' }[s] || 'info')

const searchInfo = reactive({ docNo: '', status: undefined })
const page = ref(1); const pageSize = ref(10); const total = ref(0); const tableData = ref([])
const productOpts = ref([]); const warehouseOpts = ref([]); const locationOpts = ref([])
const woOpts = ref([]); const woRequirement = ref([])

const getList = async () => {
  const res = await issueApi.getList({ page: page.value, pageSize: pageSize.value, ...searchInfo })
  if (res.code === 0) { tableData.value = res.data.list || []; total.value = res.data.total }
}
const onSearch = () => { page.value = 1; getList() }
const onReset = () => { Object.assign(searchInfo, { docNo: '', status: undefined }); onSearch() }

const dialogVisible = ref(false); const dialogType = ref('add'); const formRef = ref(null)
const dialogTitle = computed(() => ({ add: '新建领料单', edit: '编辑领料单', view: '查看领料单' }[dialogType.value]))
const readonly = computed(() => dialogType.value === 'view')
const newRow = () => ({ productId: undefined, locationId: undefined, quantity: 0, unitCost: 0, amount: 0, remark: '' })
const defaultForm = () => ({ ID: 0, docNo: '', bizDate: new Date().toISOString().slice(0, 10),
  workOrderId: undefined, warehouseId: undefined, remark: '', items: [] })
const formData = reactive(defaultForm())
const rules = {
  bizDate: [{ required: true, message: '必填' }],
  workOrderId: [{ required: true, message: '必选' }],
  warehouseId: [{ required: true, message: '必选' }]
}

const loadLocations = async (whId) => {
  if (!whId) { locationOpts.value = []; return }
  const res = await getLocationByWarehouse({ warehouseId: whId })
  if (res.code === 0) locationOpts.value = res.data.list || []
}
const onWhChange = async (whId) => {
  formData.items.forEach(r => r.locationId = undefined)
  await loadLocations(whId)
}

const onWoChange = async (woId) => {
  woRequirement.value = []; formData.items = []
  if (!woId) return
  const res = await workOrderApi.findOne({ id: woId })
  if (res.code === 0) {
    const d = res.data.record || res.data
    woRequirement.value = (d.items || []).filter(r => r.requiredQty > r.issuedQty)
    if (!formData.warehouseId && d.warehouseId) formData.warehouseId = d.warehouseId
  }
}

// 添加行：传 wo 需求时自动带物料并按剩余量预填
const addRow = (req) => {
  const r = newRow()
  if (req && req.productId) {
    r.productId = req.productId
    r.quantity = Number((req.requiredQty - req.issuedQty).toFixed(4))
  }
  formData.items.push(r)
}

const openDialog = async (type, row) => {
  dialogType.value = type
  Object.assign(formData, defaultForm())
  woRequirement.value = []; locationOpts.value = []
  // 加载可用工单（已下达 + 生产中）
  const r1 = await workOrderApi.getList({ page: 1, pageSize: 9999, status: 2 })
  const r2 = await workOrderApi.getList({ page: 1, pageSize: 9999, status: 3 })
  woOpts.value = [...(r1.code === 0 ? r1.data.list || [] : []), ...(r2.code === 0 ? r2.data.list || [] : [])]

  if (row) {
    const res = await issueApi.findOne({ id: row.ID })
    if (res.code === 0) {
      const d = res.data.record || res.data
      Object.assign(formData, { ID: d.ID, docNo: d.docNo, bizDate: d.bizDate,
        workOrderId: d.workOrderId, warehouseId: d.warehouseId, remark: d.remark,
        items: d.items || [] })
      if (d.warehouseId) await loadLocations(d.warehouseId)
      if (d.workOrderId) {
        const wo = await workOrderApi.findOne({ id: d.workOrderId })
        if (wo.code === 0) {
          const w = wo.data.record || wo.data
          // 把工单加入下拉以便显示
          if (!woOpts.value.find(x => x.ID === w.ID)) woOpts.value.push(w)
          woRequirement.value = (w.items || []).filter(r => r.requiredQty > r.issuedQty)
        }
      }
    }
  }
  dialogVisible.value = true
}
const onSubmit = async () => {
  try { await formRef.value.validate() } catch { ElMessage.error('请补全必填项'); return }
  const items = formData.items.filter(r => r.productId && r.quantity > 0)
  if (!items.length) { ElMessage.error('请至少添加一条明细'); return }
  const payload = { ID: formData.ID, bizDate: formData.bizDate, workOrderId: formData.workOrderId,
    warehouseId: formData.warehouseId, remark: formData.remark,
    items: items.map(r => ({ ID: r.ID || 0, productId: r.productId, locationId: r.locationId || 0,
      quantity: r.quantity, remark: r.remark })) }
  const fn = formData.ID ? issueApi.update : issueApi.create
  const res = await fn(payload)
  if (res.code === 0) { ElMessage.success('保存成功'); dialogVisible.value = false; getList() }
}
const onConfirm = (row) => {
  ElMessageBox.confirm(`确认领料 [${row.docNo}]? 将出库并占用工单需求量。`, '提示', { type: 'warning' }).then(async () => {
    const res = await issueApi.confirm({ id: row.ID })
    if (res.code === 0) { ElMessage.success('已确认'); getList() }
  }).catch(() => {})
}
const onVoid = (row) => {
  ElMessageBox.confirm(`作废 [${row.docNo}]? 将反向回库并扣减工单已领数。`, '提示', { type: 'warning' }).then(async () => {
    const res = await issueApi.void({ id: row.ID })
    if (res.code === 0) { ElMessage.success('已作废'); getList() }
  }).catch(() => {})
}
const onDelete = (row) => {
  ElMessageBox.confirm(`删除草稿 [${row.docNo}]?`, '提示', { type: 'warning' }).then(async () => {
    const res = await issueApi.remove({ id: row.ID })
    if (res.code === 0) { ElMessage.success('已删除'); getList() }
  }).catch(() => {})
}

onMounted(async () => {
  const [p, w] = await Promise.all([getProductList({ page: 1, pageSize: 9999 }), getWarehouseAll()])
  if (p.code === 0) productOpts.value = p.data.list || []
  if (w.code === 0) warehouseOpts.value = w.data.list || []
  getList()
})
</script>
