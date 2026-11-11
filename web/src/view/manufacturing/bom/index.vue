<!--
  MES - BOM 物料清单
-->
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
            <el-option label="启用" :value="1" />
            <el-option label="停用" :value="2" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSearch">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openDialog('add')">新建 BOM</el-button>
      </div>
      <el-table :data="tableData" border>
        <el-table-column label="成品" min-width="220">
          <template #default="{ row }">{{ row.product?.code }} - {{ row.product?.name }}</template>
        </el-table-column>
        <el-table-column prop="version" label="版本号" width="120" />
        <el-table-column prop="effectiveDate" label="生效日期" width="120" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'">{{ row.status === 1 ? '启用' : '停用' }}</el-tag>
          </template>
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
        <el-pagination v-model:current-page="page" v-model:page-size="pageSize" :total="total"
          :page-sizes="[10,20,50,100]" layout="total, sizes, prev, pager, next, jumper"
          @current-change="getList" @size-change="getList" />
      </div>
    </div>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="900px" :close-on-click-modal="false">
      <el-form ref="formRef" :model="formData" :rules="rules" label-width="90px" :disabled="readonly">
        <el-row :gutter="16">
          <el-col :span="10"><el-form-item label="成品" prop="productId">
            <el-select v-model="formData.productId" filterable style="width:100%">
              <el-option v-for="p in productOpts" :key="p.ID" :label="`${p.code} - ${p.name}`" :value="p.ID" />
            </el-select>
          </el-form-item></el-col>
          <el-col :span="6"><el-form-item label="版本号" prop="version">
            <el-input v-model="formData.version" placeholder="v1.0" />
          </el-form-item></el-col>
          <el-col :span="8"><el-form-item label="生效日期">
            <el-date-picker v-model="formData.effectiveDate" type="date" value-format="YYYY-MM-DD" style="width:100%" />
          </el-form-item></el-col>
          <el-col :span="24"><el-form-item label="备注"><el-input v-model="formData.remark" /></el-form-item></el-col>
        </el-row>
        <el-divider>子项物料</el-divider>
        <el-table :data="formData.items" border size="small">
          <el-table-column type="index" label="#" width="50" />
          <el-table-column label="物料" min-width="240">
            <template #default="{ row }">
              <el-select v-model="row.productId" filterable :disabled="readonly" style="width:100%">
                <el-option v-for="p in productOpts" :key="p.ID" :label="`${p.code} - ${p.name}`" :value="p.ID" />
              </el-select>
            </template>
          </el-table-column>
          <el-table-column label="单位用量" width="140">
            <template #default="{ row }"><el-input-number v-model="row.quantity" :min="0" :precision="4" :disabled="readonly" controls-position="right" style="width:100%" /></template>
          </el-table-column>
          <el-table-column label="损耗率(%)" width="140">
            <template #default="{ row }">
              <el-input-number :model-value="(row.lossRate || 0) * 100" :min="0" :max="100" :precision="2" :disabled="readonly"
                controls-position="right" style="width:100%" @update:model-value="v => row.lossRate = (v || 0) / 100" />
            </template>
          </el-table-column>
          <el-table-column label="备注" min-width="120">
            <template #default="{ row }"><el-input v-model="row.remark" :disabled="readonly" /></template>
          </el-table-column>
          <el-table-column v-if="!readonly" label="操作" width="70">
            <template #default="{ $index }"><el-button type="danger" link @click="formData.items.splice($index, 1)">删除</el-button></template>
          </el-table-column>
        </el-table>
        <el-button v-if="!readonly" type="primary" link icon="plus" @click="addRow" style="margin-top:8px">新增子项</el-button>
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
import { bomApi } from '@/api/manufacturing'
import { getProductList } from '@/api/inventory/product'

defineOptions({ name: 'MfgBom' })

const searchInfo = reactive({ productId: undefined, status: undefined })
const page = ref(1); const pageSize = ref(10); const total = ref(0); const tableData = ref([])
const productOpts = ref([])

const getList = async () => {
  const res = await bomApi.getList({ page: page.value, pageSize: pageSize.value, ...searchInfo })
  if (res.code === 0) { tableData.value = res.data.list || []; total.value = res.data.total }
}
const onSearch = () => { page.value = 1; getList() }
const onReset = () => { searchInfo.productId = undefined; searchInfo.status = undefined; onSearch() }

const dialogVisible = ref(false); const dialogType = ref('add'); const formRef = ref(null)
const dialogTitle = computed(() => ({ add: '新建 BOM', edit: '编辑 BOM', view: '查看 BOM' }[dialogType.value]))
const readonly = computed(() => dialogType.value === 'view')

const newRow = () => ({ productId: undefined, quantity: 1, lossRate: 0, remark: '' })
const defaultForm = () => ({ ID: 0, productId: undefined, version: 'v1.0',
  effectiveDate: new Date().toISOString().slice(0, 10), remark: '', items: [newRow()] })
const formData = reactive(defaultForm())
const rules = {
  productId: [{ required: true, message: '必选' }],
  version: [{ required: true, message: '必填' }]
}
const addRow = () => formData.items.push(newRow())

const openDialog = async (type, row) => {
  dialogType.value = type
  Object.assign(formData, defaultForm()); formData.items = [newRow()]
  if (row) {
    const res = await bomApi.findOne({ id: row.ID })
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
  const items = formData.items.filter(r => r.productId && r.quantity > 0)
  if (!items.length) { ElMessage.error('请至少添加一条子项'); return }
  const payload = { ID: formData.ID, productId: formData.productId, version: formData.version,
    effectiveDate: formData.effectiveDate, remark: formData.remark,
    items: items.map(r => ({ ID: r.ID || 0, productId: r.productId, quantity: r.quantity,
      lossRate: r.lossRate || 0, remark: r.remark })) }
  const fn = formData.ID ? bomApi.update : bomApi.create
  const res = await fn(payload)
  if (res.code === 0) { ElMessage.success('保存成功'); dialogVisible.value = false; getList() }
}
const onToggleStatus = (row, status) => {
  ElMessageBox.confirm(`${status === 1 ? '启用' : '停用'} BOM？`, '提示', { type: 'warning' }).then(async () => {
    const res = await bomApi.setStatus({ id: row.ID, status })
    if (res.code === 0) { ElMessage.success('已切换'); getList() }
  }).catch(() => {})
}
const onDelete = (row) => {
  ElMessageBox.confirm(`删除 BOM？`, '提示', { type: 'warning' }).then(async () => {
    const res = await bomApi.remove({ id: row.ID })
    if (res.code === 0) { ElMessage.success('已删除'); getList() }
  }).catch(() => {})
}

onMounted(async () => {
  const p = await getProductList({ page: 1, pageSize: 9999 })
  if (p.code === 0) productOpts.value = p.data.list || []
  getList()
})
</script>
