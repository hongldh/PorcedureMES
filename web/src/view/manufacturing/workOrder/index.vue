<!--
  MES - 生产工单
-->
<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo">
        <el-form-item label="单号"><el-input v-model="searchInfo.docNo" clearable /></el-form-item>
        <el-form-item label="成品">
          <el-select v-model="searchInfo.productId" filterable clearable style="width:220px">
            <el-option v-for="p in productOpts" :key="p.ID" :label="`${p.code} - ${p.name}`" :value="p.ID" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchInfo.status" clearable placeholder="全部" style="width:120px">
            <el-option label="草稿" :value="1" /><el-option label="已下达" :value="2" />
            <el-option label="生产中" :value="3" /><el-option label="已完工" :value="4" />
            <el-option label="已关闭" :value="5" /><el-option label="已作废" :value="9" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSearch">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="gva-table-box">
      <div class="gva-btn-list"><el-button type="primary" icon="plus" @click="openDialog('add')">新建工单</el-button></div>
      <el-table :data="tableData" border>
        <el-table-column prop="docNo" label="单号" width="180" />
        <el-table-column label="成品" min-width="200">
          <template #default="{ row }">{{ row.product?.code }} - {{ row.product?.name }}</template>
        </el-table-column>
        <el-table-column prop="planQty" label="计划数" width="100" />
        <el-table-column prop="actualQty" label="已完工" width="100" />
        <el-table-column label="入库仓" width="140"><template #default="{ row }">{{ row.warehouse?.name || '-' }}</template></el-table-column>
        <el-table-column prop="planStart" label="计划开始" width="110" />
        <el-table-column prop="planEnd" label="计划结束" width="110" />
        <el-table-column label="状态" width="100"><template #default="{ row }"><el-tag :type="statusTag(row.status)">{{ statusText(row.status) }}</el-tag></template></el-table-column>
        <el-table-column label="操作" width="300" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="openDialog('view', row)">查看</el-button>
            <el-button v-if="row.status === 1" type="primary" link @click="openDialog('edit', row)">编辑</el-button>
            <el-button v-if="row.status === 1" type="success" link @click="onRelease(row)">下达</el-button>
            <el-button v-if="row.status === 1" type="danger" link @click="onDelete(row)">删除</el-button>
            <el-button v-if="row.status >= 2 && row.status <= 4" type="warning" link @click="onClose(row)">关闭</el-button>
            <el-button v-if="row.status !== 9 && row.status !== 5" type="danger" link @click="onVoid(row)">作废</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination v-model:current-page="page" v-model:page-size="pageSize" :total="total" :page-sizes="[10,20,50,100]"
          layout="total, sizes, prev, pager, next, jumper" @current-change="getList" @size-change="getList" />
      </div>
    </div>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="900px" :close-on-click-modal="false">
      <el-form ref="formRef" :model="formData" :rules="rules" label-width="100px" :disabled="readonly">
        <el-row :gutter="16">
          <el-col :span="8"><el-form-item label="单号"><el-input v-model="formData.docNo" disabled placeholder="保存后生成" /></el-form-item></el-col>
          <el-col :span="8"><el-form-item label="成品" prop="productId">
            <el-select v-model="formData.productId" filterable style="width:100%" @change="onProductChange">
              <el-option v-for="p in productOpts" :key="p.ID" :label="`${p.code} - ${p.name}`" :value="p.ID" />
            </el-select>
          </el-form-item></el-col>
          <el-col :span="8"><el-form-item label="BOM 版本" prop="bomId">
            <el-select v-model="formData.bomId" style="width:100%" placeholder="先选成品">
              <el-option v-for="b in bomOpts" :key="b.ID" :label="b.version" :value="b.ID" />
            </el-select>
          </el-form-item></el-col>
          <el-col :span="8"><el-form-item label="工艺路线" prop="routingId">
            <el-select v-model="formData.routingId" clearable style="width:100%" placeholder="先选成品">
              <el-option v-for="r in routingOpts" :key="r.ID" :label="r.version" :value="r.ID" />
            </el-select>
          </el-form-item></el-col>
          <el-col :span="8"><el-form-item label="计划数量" prop="planQty"><el-input-number v-model="formData.planQty" :min="0.0001" :precision="4" controls-position="right" style="width:100%" /></el-form-item></el-col>
          <el-col :span="8"><el-form-item label="计划开始"><el-date-picker v-model="formData.planStart" type="date" value-format="YYYY-MM-DD" style="width:100%" /></el-form-item></el-col>
          <el-col :span="8"><el-form-item label="计划结束"><el-date-picker v-model="formData.planEnd" type="date" value-format="YYYY-MM-DD" style="width:100%" /></el-form-item></el-col>
          <el-col :span="8"><el-form-item label="入库仓">
            <el-select v-model="formData.warehouseId" style="width:100%">
              <el-option v-for="w in warehouseOpts" :key="w.ID" :label="w.name" :value="w.ID" />
            </el-select>
          </el-form-item></el-col>
          <el-col :span="16"><el-form-item label="备注"><el-input v-model="formData.remark" /></el-form-item></el-col>
        </el-row>
        <div v-if="dialogType === 'view' && formData.items && formData.items.length">
          <el-divider>物料需求（下达时已快照）</el-divider>
          <el-table :data="formData.items" border size="small">
            <el-table-column type="index" label="#" width="50" />
            <el-table-column label="物料" min-width="240"><template #default="{ row }">{{ row.product?.code }} - {{ row.product?.name }}</template></el-table-column>
            <el-table-column prop="requiredQty" label="需求量" width="140" />
            <el-table-column prop="issuedQty" label="已领数" width="140" />
            <el-table-column label="剩余" width="140">
              <template #default="{ row }">{{ Number((row.requiredQty - row.issuedQty) || 0).toFixed(4) }}</template>
            </el-table-column>
          </el-table>
        </div>
      </el-form>
      <div v-if="dialogType === 'view'">
        <el-divider>工序与派工</el-divider>
        <el-table :data="opsList" border size="small">
          <el-table-column prop="seq" label="#" width="50" />
          <el-table-column prop="code" label="工序编码" width="120" />
          <el-table-column prop="name" label="工序名称" min-width="140" />
          <el-table-column label="工位" width="120"><template #default="{ row }">{{ row.workstation?.code || '-' }}</template></el-table-column>
          <el-table-column label="指派" width="140"><template #default="{ row }">{{ row.assignedUser?.nickName || row.assignedUser?.userName || row.team?.name || '-' }}</template></el-table-column>
          <el-table-column label="计划数" width="110"><template #default="{ row }">{{ Number(row.planQty || 0).toFixed(4) }}</template></el-table-column>
          <el-table-column label="合格" width="110"><template #default="{ row }">{{ Number(row.doneQty || 0).toFixed(4) }}</template></el-table-column>
          <el-table-column label="不良" width="110"><template #default="{ row }">{{ Number(row.rejectQty || 0).toFixed(4) }}</template></el-table-column>
          <el-table-column label="工时" width="100"><template #default="{ row }">{{ Number(row.actualHours || 0).toFixed(2) }}</template></el-table-column>
          <el-table-column label="状态" width="100"><template #default="{ row }"><el-tag :type="opStatusTag(row.status)">{{ opStatusText(row.status) }}</el-tag></template></el-table-column>
          <el-table-column label="操作" width="80" fixed="right">
            <template #default="{ row }">
              <el-button v-if="row.status !== 4 && (workOrderStatus === 2 || workOrderStatus === 3)" type="primary" link @click="openDispatchDialog(row)">{{ row.status === 1 ? '派工' : '改派' }}</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
      <el-dialog v-model="dispatchVisible" title="工序派工" width="560px" append-to-body :close-on-click-modal="false">
        <el-form :model="dispatchForm" label-width="100px">
          <el-form-item label="员工">
            <el-select v-model="dispatchForm.assignedUserId" filterable clearable style="width:100%">
              <el-option v-for="u in userOpts" :key="u.ID" :label="`${u.nickName || u.userName}（${u.userName}）`" :value="u.ID" />
            </el-select>
          </el-form-item>
          <el-form-item label="班组">
            <el-select v-model="dispatchForm.teamId" clearable style="width:100%">
              <el-option v-for="t in teamOpts" :key="t.ID" :label="`${t.code} - ${t.name}`" :value="t.ID" />
            </el-select>
          </el-form-item>
          <el-form-item label="工位">
            <el-select v-model="dispatchForm.workstationId" clearable style="width:100%">
              <el-option v-for="w in workstationOpts" :key="w.ID" :label="`${w.code} - ${w.name}`" :value="w.ID" />
            </el-select>
          </el-form-item>
        </el-form>
        <template #footer>
          <el-button @click="dispatchVisible = false">关闭</el-button>
          <el-button type="primary" :loading="dispatchSubmitting" @click="onDispatch">保存</el-button>
        </template>
      </el-dialog>
      <template #footer>
        <el-button @click="dialogVisible = false">关闭</el-button>
        <el-button v-if="!readonly" type="primary" @click="onSubmit">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { workOrderApi, bomApi, routingApi, dispatchApi, workstationApi, teamApi } from '@/api/manufacturing'
import { getUserList } from '@/api/user'
import { getProductList } from '@/api/inventory/product'
import { getWarehouseAll } from '@/api/inventory/warehouse'

defineOptions({ name: 'MfgWorkOrder' })

const statusText = (s) => ({ 1: '草稿', 2: '已下达', 3: '生产中', 4: '已完工', 5: '已关闭', 9: '已作废' }[s] || '-')
const statusTag = (s) => ({ 1: 'info', 2: 'primary', 3: 'warning', 4: 'success', 5: '', 9: 'danger' }[s] || 'info')
const opStatusText = (s) => ({ 1: '未派工', 2: '已派工', 3: '进行中', 4: '已完工' }[s] || '-')
const opStatusTag = (s) => ({ 1: 'info', 2: 'primary', 3: 'warning', 4: 'success' }[s] || 'info')

const searchInfo = reactive({ docNo: '', productId: undefined, status: undefined })
const page = ref(1); const pageSize = ref(10); const total = ref(0); const tableData = ref([])
const productOpts = ref([]); const warehouseOpts = ref([]); const bomOpts = ref([]); const routingOpts = ref([])
const opsList = ref([]); const workOrderStatus = ref(0)
const userOpts = ref([]); const teamOpts = ref([]); const workstationOpts = ref([])
const dispatchVisible = ref(false); const dispatchSubmitting = ref(false)
const dispatchForm = reactive({ opId: 0, assignedUserId: undefined, teamId: undefined, workstationId: undefined })

const getOps = async (id) => {
  const res = await dispatchApi.getOps({ id })
  if (res.code === 0) opsList.value = res.data.list || []
}
const openDispatchDialog = async (row) => {
  Object.assign(dispatchForm, { opId: row.ID, assignedUserId: row.assignedUserId || undefined,
    teamId: row.teamId || undefined, workstationId: row.workstationId || undefined })
  const [u, t, w] = await Promise.all([getUserList({ page: 1, pageSize: 9999 }),
    teamApi.getList({ page: 1, pageSize: 9999 }), workstationApi.getAll()])
  if (u.code === 0) userOpts.value = u.data.list || []
  if (t.code === 0) teamOpts.value = t.data.list || []
  if (w.code === 0) workstationOpts.value = w.data.list || []
  dispatchVisible.value = true
}
const onDispatch = async () => {
  if (dispatchSubmitting.value) return
  if (!dispatchForm.assignedUserId && !dispatchForm.teamId) { ElMessage.error('员工和班组至少选择一个'); return }
  dispatchSubmitting.value = true
  try {
    const res = await dispatchApi.dispatch({ opId: dispatchForm.opId, assignedUserId: dispatchForm.assignedUserId || 0,
      teamId: dispatchForm.teamId || 0, workstationId: dispatchForm.workstationId || 0 })
    if (res.code === 0) { ElMessage.success('派工成功'); dispatchVisible.value = false; await getOps(formData.ID) }
  } finally { dispatchSubmitting.value = false }
}

const getList = async () => {
  const res = await workOrderApi.getList({ page: page.value, pageSize: pageSize.value, ...searchInfo })
  if (res.code === 0) { tableData.value = res.data.list || []; total.value = res.data.total }
}
const onSearch = () => { page.value = 1; getList() }
const onReset = () => { Object.assign(searchInfo, { docNo: '', productId: undefined, status: undefined }); onSearch() }

const dialogVisible = ref(false); const dialogType = ref('add'); const formRef = ref(null)
const dialogTitle = computed(() => ({ add: '新建工单', edit: '编辑工单', view: '查看工单' }[dialogType.value]))
const readonly = computed(() => dialogType.value === 'view')
const defaultForm = () => ({ ID: 0, docNo: '', productId: undefined, bomId: undefined, routingId: undefined, planQty: 1,
  planStart: '', planEnd: '', warehouseId: undefined, remark: '', items: [] })
const formData = reactive(defaultForm())
const rules = {
  productId: [{ required: true, message: '必选' }],
  bomId: [{ required: true, message: '必选' }],
  routingId: [{ required: true, message: '必选' }],
  planQty: [{ required: true, message: '必填' }]
}

const onProductChange = async (pid) => {
  formData.bomId = undefined; bomOpts.value = []
  formData.routingId = undefined; routingOpts.value = []
  if (!pid) return
  const [bRes, rRes] = await Promise.all([
    bomApi.getActiveByProduct({ id: pid }),
    routingApi.getActiveByProduct({ id: pid })
  ])
  if (bRes.code === 0) bomOpts.value = bRes.data.list || []
  if (rRes.code === 0) routingOpts.value = rRes.data.list || []
}

const openDialog = async (type, row) => {
  dialogType.value = type
  Object.assign(formData, defaultForm())
  opsList.value = []; workOrderStatus.value = 0; dispatchVisible.value = false
  if (row) {
    const [res] = await Promise.all([workOrderApi.findOne({ id: row.ID }), type === 'view' ? getOps(row.ID) : Promise.resolve()])
    if (res.code === 0) {
      const d = res.data.record || res.data
      workOrderStatus.value = d.status
      Object.assign(formData, { ID: d.ID, docNo: d.docNo, productId: d.productId, bomId: d.bomId, routingId: d.routingId,
        planQty: d.planQty, planStart: d.planStart, planEnd: d.planEnd, warehouseId: d.warehouseId,
        remark: d.remark, items: d.items || [] })
      if (d.productId) await onProductChange(d.productId)
      formData.bomId = d.bomId  // onProductChange 会清空，重置回来
      formData.routingId = d.routingId
    }
  }
  dialogVisible.value = true
}
const onSubmit = async () => {
  try { await formRef.value.validate() } catch { ElMessage.error('请补全必填项'); return }
  const fn = formData.ID ? workOrderApi.update : workOrderApi.create
  const res = await fn(formData)
  if (res.code === 0) { ElMessage.success('保存成功'); dialogVisible.value = false; getList() }
}
const onRelease = (row) => {
  ElMessageBox.confirm(`下达工单 [${row.docNo}]? 将按 BOM 展开物料需求、按工艺路线展开工序。`, '提示', { type: 'warning' }).then(async () => {
    const res = await workOrderApi.release({ id: row.ID })
    if (res.code === 0) { ElMessage.success('已下达'); getList() }
  }).catch(() => {})
}
const onClose = (row) => {
  ElMessageBox.confirm(`关闭工单 [${row.docNo}]?`, '提示', { type: 'warning' }).then(async () => {
    const res = await workOrderApi.close({ id: row.ID })
    if (res.code === 0) { ElMessage.success('已关闭'); getList() }
  }).catch(() => {})
}
const onVoid = (row) => {
  ElMessageBox.confirm(`作废工单 [${row.docNo}]?`, '提示', { type: 'warning' }).then(async () => {
    const res = await workOrderApi.void({ id: row.ID })
    if (res.code === 0) { ElMessage.success('已作废'); getList() }
  }).catch(() => {})
}
const onDelete = (row) => {
  ElMessageBox.confirm(`删除草稿 [${row.docNo}]?`, '提示', { type: 'warning' }).then(async () => {
    const res = await workOrderApi.remove({ id: row.ID })
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
