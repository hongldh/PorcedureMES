<template>
  <div>
    <div class="gva-table-box">
      <el-card shadow="never">
        <template #header>
          <el-row justify="space-between" align="middle">
            <span>我的待报工工序</span>
            <el-button icon="refresh" circle aria-label="刷新我的待报工工序" @click="getMyOps" />
          </el-row>
        </template>
        <el-table :data="myOpsList" border empty-text="暂无指派给您的待报工工序">
          <el-table-column prop="workOrderId" label="工单ID" width="100" />
          <el-table-column prop="seq" label="工序号" width="80" />
          <el-table-column prop="code" label="工序编码" width="120" />
          <el-table-column prop="name" label="工序名称" min-width="140" />
          <el-table-column label="工位" width="120"><template #default="{ row }">{{ row.workstation?.code || '-' }}</template></el-table-column>
          <el-table-column label="计划" width="110"><template #default="{ row }">{{ Number(row.planQty || 0).toFixed(4) }}</template></el-table-column>
          <el-table-column label="已合格" width="110"><template #default="{ row }">{{ Number(row.doneQty || 0).toFixed(4) }}</template></el-table-column>
          <el-table-column label="已不良" width="110"><template #default="{ row }">{{ Number(row.rejectQty || 0).toFixed(4) }}</template></el-table-column>
          <el-table-column label="剩余" width="110"><template #default="{ row }">{{ Number(remainingQty(row) || 0).toFixed(4) }}</template></el-table-column>
          <el-table-column label="状态" width="100"><template #default="{ row }"><el-tag :type="opStatusTag(row.status)">{{ opStatusText(row.status) }}</el-tag></template></el-table-column>
          <el-table-column label="操作" width="80" fixed="right">
            <template #default="{ row }"><el-button type="primary" link @click="openDialog('add', row)">报工</el-button></template>
          </el-table-column>
        </el-table>
      </el-card>
    </div>

    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo">
        <el-form-item label="单号"><el-input v-model="searchInfo.docNo" clearable /></el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchInfo.status" placeholder="全部" clearable style="width:120px">
            <el-option label="已生效" :value="3" /><el-option label="作废" :value="4" />
          </el-select>
        </el-form-item>
        <el-form-item><el-button type="primary" icon="search" @click="onSearch">查询</el-button><el-button icon="refresh" @click="onReset">重置</el-button></el-form-item>
      </el-form>
    </div>

    <div class="gva-table-box">
      <div class="gva-btn-list">报工记录</div>
      <el-table :data="tableData" border>
        <el-table-column prop="docNo" label="单号" width="180" />
        <el-table-column prop="reportTime" label="报工时间" width="180" />
        <el-table-column label="工单号" width="180"><template #default="{ row }">{{ row.workOrder?.docNo || '-' }}</template></el-table-column>
        <el-table-column label="成品" min-width="200"><template #default="{ row }">{{ row.workOrder?.product?.code || '-' }} - {{ row.workOrder?.product?.name || '-' }}</template></el-table-column>
        <el-table-column label="工序" min-width="180"><template #default="{ row }">{{ row.op?.code || '-' }} {{ row.op?.name || '-' }}</template></el-table-column>
        <el-table-column label="报工人" width="120"><template #default="{ row }">{{ row.user?.nickName || row.user?.userName || '-' }}</template></el-table-column>
        <el-table-column label="工位" width="120"><template #default="{ row }">{{ row.workstation?.code || '-' }}</template></el-table-column>
        <el-table-column label="合格数" width="110"><template #default="{ row }">{{ Number(row.goodQty || 0).toFixed(4) }}</template></el-table-column>
        <el-table-column label="不良数" width="110"><template #default="{ row }">{{ Number(row.rejectQty || 0).toFixed(4) }}</template></el-table-column>
        <el-table-column label="工时" width="100"><template #default="{ row }">{{ Number(row.hours || 0).toFixed(2) }}</template></el-table-column>
        <el-table-column label="状态" width="100"><template #default="{ row }"><el-tag :type="statusTag(row.status)">{{ statusText(row.status) }}</el-tag></template></el-table-column>
        <el-table-column prop="remark" label="备注" min-width="160" />
        <el-table-column label="操作" width="80" fixed="right">
          <template #default="{ row }"><el-button v-if="row.status === 3" type="danger" link @click="onVoid(row)">作废</el-button></template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination v-model:current-page="page" v-model:page-size="pageSize" :total="total" :page-sizes="[10,20,50,100]"
          layout="total, sizes, prev, pager, next, jumper" @current-change="getList" @size-change="getList" />
      </div>
    </div>

    <el-dialog v-model="dialogVisible" title="工序报工" width="640px" :close-on-click-modal="false">
      <el-form :model="formData" label-width="100px">
        <el-form-item label="工序">{{ selectedOp?.code }} - {{ selectedOp?.name }}</el-form-item>
        <el-form-item label="剩余数量">{{ Number(remainingQty(selectedOp) || 0).toFixed(4) }}</el-form-item>
        <el-form-item label="合格数"><el-input-number v-model="formData.goodQty" :min="0" :precision="4" controls-position="right" style="width:100%" /></el-form-item>
        <el-form-item label="不良数"><el-input-number v-model="formData.rejectQty" :min="0" :precision="4" controls-position="right" style="width:100%" /></el-form-item>
        <el-form-item label="工时(小时)"><el-input-number v-model="formData.hours" :min="0" :precision="2" controls-position="right" style="width:100%" /></el-form-item>
        <el-form-item label="工位">
          <el-select v-model="formData.workstationId" clearable style="width:100%">
            <el-option v-for="w in workstationOpts" :key="w.ID" :label="`${w.code} - ${w.name}`" :value="w.ID" />
          </el-select>
        </el-form-item>
        <el-form-item label="备注"><el-input v-model="formData.remark" type="textarea" :rows="2" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">关闭</el-button>
        <el-button type="primary" :loading="submitting" @click="onSubmit">报工</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { workReportApi, workstationApi } from '@/api/manufacturing'

defineOptions({ name: 'MfgWorkReport' })

const statusText = (s) => ({ 3: '已生效', 4: '作废' }[s] || '-')
const statusTag = (s) => ({ 3: 'success', 4: 'danger' }[s] || 'info')
const opStatusText = (s) => ({ 1: '未派工', 2: '已派工', 3: '进行中', 4: '已完工' }[s] || '-')
const opStatusTag = (s) => ({ 1: 'info', 2: 'primary', 3: 'warning', 4: 'success' }[s] || 'info')
const remainingQty = (op) => Number((Number(op?.planQty || 0) - Number(op?.doneQty || 0) - Number(op?.rejectQty || 0)).toFixed(4))

const searchInfo = reactive({ docNo: '', status: undefined })
const page = ref(1); const pageSize = ref(10); const total = ref(0); const tableData = ref([])
const myOpsList = ref([]); const workstationOpts = ref([])

const getList = async () => {
  const res = await workReportApi.getList({ page: page.value, pageSize: pageSize.value, ...searchInfo })
  if (res.code === 0) { tableData.value = res.data.list || []; total.value = res.data.total }
}
const getMyOps = async () => {
  const res = await workReportApi.myOps()
  if (res.code === 0) myOpsList.value = res.data.list || []
}
const onSearch = () => { page.value = 1; getList() }
const onReset = () => { Object.assign(searchInfo, { docNo: '', status: undefined }); onSearch() }

const dialogVisible = ref(false); const dialogType = ref('add'); const selectedOp = ref(null); const submitting = ref(false)
const defaultForm = () => ({ opId: 0, goodQty: 0, rejectQty: 0, hours: 0, workstationId: undefined, remark: '' })
const formData = reactive(defaultForm())

const openDialog = (type, row) => {
  dialogType.value = type
  selectedOp.value = row
  Object.assign(formData, defaultForm(), { opId: row.ID, workstationId: row.workstationId || undefined })
  dialogVisible.value = true
}
const onSubmit = async () => {
  if (submitting.value) return
  const goodQty = Number(formData.goodQty || 0); const rejectQty = Number(formData.rejectQty || 0)
  const hours = Number(formData.hours || 0)
  if (![goodQty, rejectQty, hours].every(n => Number.isFinite(n) && n >= 0)) { ElMessage.error('数量和工时必须为非负数'); return }
  const quantity = Number((goodQty + rejectQty).toFixed(4))
  if (quantity <= 0) { ElMessage.error('合格数与不良数之和必须大于 0'); return }
  if (quantity > remainingQty(selectedOp.value)) { ElMessage.error('合格数与不良数之和不能超过剩余数量'); return }
  submitting.value = true
  try {
    const res = await workReportApi.create({ opId: formData.opId, goodQty, rejectQty, hours,
      workstationId: formData.workstationId || 0, remark: formData.remark })
    if (res.code === 0) {
      ElMessage.success(`报工成功 ${res.data.docNo}`); dialogVisible.value = false
      await Promise.all([getMyOps(), getList()])
    }
  } finally { submitting.value = false }
}
const onVoid = (row) => {
  ElMessageBox.confirm(`作废报工 [${row.docNo}]?`, '提示', { type: 'warning' }).then(async () => {
    const res = await workReportApi.void({ id: row.ID })
    if (res.code === 0) { ElMessage.success('已作废'); await Promise.all([getMyOps(), getList()]) }
  }).catch(() => {})
}

onMounted(async () => {
  const w = await workstationApi.getAll()
  if (w.code === 0) workstationOpts.value = w.data.list || []
  await Promise.all([getMyOps(), getList()])
})
</script>
