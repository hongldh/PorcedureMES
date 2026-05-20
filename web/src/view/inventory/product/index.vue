<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo">
        <el-form-item label="关键字"><el-input v-model="searchInfo.keyword" placeholder="编码/名称/条码" clearable /></el-form-item>
        <el-form-item label="分类">
          <el-select v-model="searchInfo.categoryId" placeholder="全部" clearable style="width:160px">
            <el-option v-for="c in categoryOpts" :key="c.ID" :label="c.name" :value="c.ID" />
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
        <el-button type="primary" icon="plus" @click="openDialog('add')">新增</el-button>
        <el-button type="danger" icon="delete" :disabled="!selectedIds.length" @click="onBatchDelete">批量删除</el-button>
      </div>
      <el-table :data="tableData" @selection-change="onSelectionChange" border>
        <el-table-column type="selection" width="55" />
        <el-table-column prop="code" label="SKU编码" width="160" />
        <el-table-column prop="name" label="商品名称" />
        <el-table-column prop="spec" label="规格" width="120" />
        <el-table-column label="分类" width="120">
          <template #default="{ row }">{{ row.category?.name || '-' }}</template>
        </el-table-column>
        <el-table-column label="单位" width="80">
          <template #default="{ row }">{{ row.unit?.name || '-' }}</template>
        </el-table-column>
        <el-table-column prop="purchasePrice" label="采购价" width="100" />
        <el-table-column prop="salePrice" label="销售价" width="100" />
        <el-table-column prop="safetyStock" label="安全库存" width="100" />
        <el-table-column label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'">{{ row.status === 1 ? '启用' : '停用' }}</el-tag>
          </template>
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
    <el-dialog v-model="dialogVisible" :title="dialogType === 'add' ? '新增商品' : '编辑商品'" width="640px">
      <el-form ref="formRef" :model="formData" :rules="rules" label-width="100px">
        <el-row :gutter="16">
          <el-col :span="12"><el-form-item label="SKU编码" prop="code"><el-input v-model="formData.code" /></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="商品名称" prop="name"><el-input v-model="formData.name" /></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="规格"><el-input v-model="formData.spec" /></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="条形码"><el-input v-model="formData.barcode" /></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="分类">
            <el-select v-model="formData.categoryId" clearable style="width:100%">
              <el-option v-for="c in categoryOpts" :key="c.ID" :label="c.name" :value="c.ID" />
            </el-select>
          </el-form-item></el-col>
          <el-col :span="12"><el-form-item label="单位">
            <el-select v-model="formData.unitId" clearable style="width:100%">
              <el-option v-for="u in unitOpts" :key="u.ID" :label="u.name" :value="u.ID" />
            </el-select>
          </el-form-item></el-col>
          <el-col :span="8"><el-form-item label="采购价"><el-input-number v-model="formData.purchasePrice" :min="0" :precision="4" controls-position="right" /></el-form-item></el-col>
          <el-col :span="8"><el-form-item label="销售价"><el-input-number v-model="formData.salePrice" :min="0" :precision="4" controls-position="right" /></el-form-item></el-col>
          <el-col :span="8"><el-form-item label="安全库存"><el-input-number v-model="formData.safetyStock" :min="0" :precision="4" controls-position="right" /></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="状态">
            <el-radio-group v-model="formData.status"><el-radio :value="1">启用</el-radio><el-radio :value="2">停用</el-radio></el-radio-group>
          </el-form-item></el-col>
          <el-col :span="24"><el-form-item label="备注"><el-input v-model="formData.remark" type="textarea" /></el-form-item></el-col>
        </el-row>
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
import { createProduct, updateProduct, deleteProduct, deleteProductBatch, getProductList } from '@/api/inventory/product'
import { getCategoryAll } from '@/api/inventory/category'
import { getUnitAll } from '@/api/inventory/unit'

defineOptions({ name: 'InvProduct' })

const searchInfo = reactive({ keyword: '', categoryId: undefined })
const page = ref(1); const pageSize = ref(10); const total = ref(0)
const tableData = ref([]); const selectedIds = ref([])
const categoryOpts = ref([]); const unitOpts = ref([])
const dialogVisible = ref(false); const dialogType = ref('add'); const formRef = ref(null)
const defaultForm = () => ({ ID: 0, code: '', name: '', spec: '', barcode: '', categoryId: undefined, unitId: undefined, purchasePrice: 0, salePrice: 0, safetyStock: 0, status: 1, remark: '' })
const formData = reactive(defaultForm())
const rules = {
  code: [{ required: true, message: '必填', trigger: 'blur' }],
  name: [{ required: true, message: '必填', trigger: 'blur' }],
}
const getList = async () => {
  const res = await getProductList({ page: page.value, pageSize: pageSize.value, ...searchInfo })
  if (res.code === 0) { tableData.value = res.data.list || []; total.value = res.data.total }
}
const onSearch = () => { page.value = 1; getList() }
const onReset = () => { searchInfo.keyword = ''; searchInfo.categoryId = undefined; onSearch() }
const onSelectionChange = (rows) => { selectedIds.value = rows.map(r => r.ID) }
const openDialog = (type, row) => {
  dialogType.value = type; Object.assign(formData, defaultForm())
  if (type === 'edit' && row) Object.assign(formData, row)
  dialogVisible.value = true
}
const onSubmit = async () => {
  await formRef.value.validate()
  const fn = dialogType.value === 'add' ? createProduct : updateProduct
  const res = await fn(formData)
  if (res.code === 0) { ElMessage.success(res.msg); dialogVisible.value = false; getList() }
}
const onDelete = (row) => {
  ElMessageBox.confirm(`确认删除 [${row.name}] ?`, '提示', { type: 'warning' }).then(async () => {
    const res = await deleteProduct({ id: row.ID }); if (res.code === 0) { ElMessage.success('删除成功'); getList() }
  }).catch(() => {})
}
const onBatchDelete = () => {
  ElMessageBox.confirm(`删除 ${selectedIds.value.length} 项?`, '提示', { type: 'warning' }).then(async () => {
    const res = await deleteProductBatch({ ids: selectedIds.value })
    if (res.code === 0) { ElMessage.success('成功'); getList() }
  }).catch(() => {})
}
const loadOpts = async () => {
  const [c, u] = await Promise.all([getCategoryAll(), getUnitAll()])
  if (c.code === 0) categoryOpts.value = c.data.list || []
  if (u.code === 0) unitOpts.value = u.data.list || []
}
onMounted(() => { loadOpts(); getList() })
</script>
