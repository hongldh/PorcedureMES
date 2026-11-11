<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo">
        <el-form-item label="成品">
          <el-select v-model="searchInfo.productId" filterable clearable style="width:240px">
            <el-option v-for="p in productOpts" :key="p.ID" :label="`${p.code} - ${p.name}`" :value="p.ID" />
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
      <div class="gva-btn-list"><el-button type="primary" icon="plus" @click="openDialog('add')">新建工艺路线</el-button></div>
      <el-table :data="tableData" border>
        <el-table-column label="成品" min-width="220">
          <template #default="{ row }">{{ row.product?.code }} - {{ row.product?.name }}</template>
        </el-table-column>
        <el-table-column prop="version" label="版本号" width="120" />
        <el-table-column prop="effectiveDate" label="生效日期" width="120" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }"><el-tag :type="row.status === 1 ? 'success' : 'info'">{{ row.status === 1 ? '启用' : '停用' }}</el-tag></template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" />
        <el-table-column label="操作" width="260" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="openDialog('view', row)">查看</el-button>
            <el-button type="primary" link @click="openDialog('edit', row)">编辑</el-button>
            <el-button v-if="row.status === 1" type="warning" link @click="onToggleStatus(row, 2)">停用</el-button>
            <el-button v-else type="success" link @click="onToggleStatus(row, 1)">启用</el-button>
            <el-button type="danger" link @click="onDelete(row)">删除</el-button>
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
          <el-col :span="10"><el-form-item label="成品" prop="productId">
            <el-select v-model="formData.productId" filterable style="width:100%">
              <el-option v-for="p in productOpts" :key="p.ID" :label="`${p.code} - ${p.name}`" :value="p.ID" />
            </el-select>
          </el-form-item></el-col>
          <el-col :span="6"><el-form-item label="版本号" prop="version"><el-input v-model="formData.version" placeholder="v1.0" /></el-form-item></el-col>
          <el-col :span="8"><el-form-item label="生效日期">
            <el-date-picker v-model="formData.effectiveDate" type="date" value-format="YYYY-MM-DD" style="width:100%" />
          </el-form-item></el-col>
          <el-col :span="24"><el-form-item label="备注"><el-input v-model="formData.remark" /></el-form-item></el-col>
        </el-row>
        <el-divider>工序列表</el-divider>
        <el-table :data="formData.items" border size="small">
          <el-table-column type="index" label="#" width="50" />
          <el-table-column label="工序号" width="100">
            <template #default="{ row }"><el-input-number v-model="row.seq" :min="0" :step="10" :disabled="readonly" controls-position="right" style="width:100%" /></template>
          </el-table-column>
          <el-table-column label="工序编码" width="140">
            <template #default="{ row }"><el-input v-model="row.code" :disabled="readonly" placeholder="OP010" /></template>
          </el-table-column>
          <el-table-column label="工序名称" min-width="160">
            <template #default="{ row }"><el-input v-model="row.name" :disabled="readonly" placeholder="车削/铣削/装配" /></template>
          </el-table-column>
          <el-table-column label="默认工位" width="160">
            <template #default="{ row }">
              <el-select v-model="row.workstationId" clearable :disabled="readonly" style="width:100%">
                <el-option v-for="w in workstationOpts" :key="w.ID" :label="w.code" :value="w.ID" />
              </el-select>
            </template>
          </el-table-column>
          <el-table-column label="标准工时(h)" width="130">
            <template #default="{ row }"><el-input-number v-model="row.stdHours" :min="0" :precision="4" :disabled="readonly" controls-position="right" style="width:100%" /></template>
          </el-table-column>
          <el-table-column label="关键工序" width="90">
            <template #default="{ row }"><el-switch v-model="row.isKey" :disabled="readonly" /></template>
          </el-table-column>
          <el-table-column v-if="!readonly" label="操作" width="70">
            <template #default="{ $index }"><el-button type="danger" link @click="formData.items.splice($index, 1)">删除</el-button></template>
          </el-table-column>
        </el-table>
        <el-button v-if="!readonly" type="primary" link icon="plus" @click="addRow" style="margin-top:8px">新增工序</el-button>
      </el-form>
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
import { routingApi, workstationApi } from '@/api/manufacturing'
import { getProductList } from '@/api/inventory/product'

defineOptions({ name: 'MfgRouting' })

const searchInfo = reactive({ productId: undefined, status: undefined })
const page = ref(1); const pageSize = ref(10); const total = ref(0); const tableData = ref([])
const productOpts = ref([]); const workstationOpts = ref([])

const getList = async () => {
  const res = await routingApi.getList({ page: page.value, pageSize: pageSize.value, ...searchInfo })
  if (res.code === 0) { tableData.value = res.data.list || []; total.value = res.data.total }
}
const onSearch = () => { page.value = 1; getList() }
const onReset = () => { searchInfo.productId = undefined; searchInfo.status = undefined; onSearch() }

const dialogVisible = ref(false); const dialogType = ref('add'); const formRef = ref(null)
const dialogTitle = computed(() => ({ add: '新建工艺路线', edit: '编辑工艺路线', view: '查看工艺路线' }[dialogType.value]))
const readonly = computed(() => dialogType.value === 'view')
const newRow = () => ({ seq: 10, code: '', name: '', workstationId: undefined, stdHours: 0, isKey: false })
const defaultForm = () => ({ ID: 0, productId: undefined, version: 'v1.0',
  effectiveDate: new Date().toISOString().slice(0, 10), remark: '', items: [newRow()] })
const formData = reactive(defaultForm())
const rules = {
  productId: [{ required: true, message: '必选' }],
  version: [{ required: true, message: '必填' }]
}
const addRow = () => {
  const last = formData.items[formData.items.length - 1]
  const seq = last ? Number(last.seq || 0) + 10 : 10
  formData.items.push({ ...newRow(), seq })
}

const openDialog = async (type, row) => {
  dialogType.value = type
  Object.assign(formData, defaultForm()); formData.items = [newRow()]
  if (row) {
    const res = await routingApi.findOne({ id: row.ID })
    if (res.code === 0) {
      const d = res.data.record || res.data
      Object.assign(formData, { ID: d.ID, productId: d.productId, version: d.version,
        effectiveDate: d.effectiveDate, remark: d.remark, items: d.items || [newRow()] })
    }
  }
  dialogVisible.value = true
}
const onSubmit = async () => {
  try { await formRef.value.validate() } catch { ElMessage.error('请补全必填项'); return }
  const items = formData.items.filter(r => r.code && r.name)
  if (!items.length) { ElMessage.error('请至少添加一道工序'); return }
  const payload = { ID: formData.ID, productId: formData.productId, version: formData.version,
    effectiveDate: formData.effectiveDate, remark: formData.remark,
    items: items.map(r => ({ ID: r.ID || 0, seq: r.seq, code: r.code, name: r.name,
      workstationId: r.workstationId || 0, stdHours: r.stdHours || 0,
      isKey: !!r.isKey, description: r.description || '', remark: r.remark || '' })) }
  const fn = formData.ID ? routingApi.update : routingApi.create
  const res = await fn(payload)
  if (res.code === 0) { ElMessage.success('保存成功'); dialogVisible.value = false; getList() }
}
const onToggleStatus = (row, status) => {
  ElMessageBox.confirm(`${status === 1 ? '启用' : '停用'} 工艺路线？`, '提示', { type: 'warning' }).then(async () => {
    const res = await routingApi.setStatus({ id: row.ID, status })
    if (res.code === 0) { ElMessage.success('已切换'); getList() }
  }).catch(() => {})
}
const onDelete = (row) => {
  ElMessageBox.confirm(`删除工艺路线?`, '提示', { type: 'warning' }).then(async () => {
    const res = await routingApi.remove({ id: row.ID })
    if (res.code === 0) { ElMessage.success('已删除'); getList() }
  }).catch(() => {})
}

onMounted(async () => {
  const [p, w] = await Promise.all([getProductList({ page: 1, pageSize: 9999 }), workstationApi.getAll()])
  if (p.code === 0) productOpts.value = p.data.list || []
  if (w.code === 0) workstationOpts.value = w.data.list || []
  getList()
})
</script>
