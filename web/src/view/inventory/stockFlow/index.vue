<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo">
        <el-form-item label="商品">
          <el-select v-model="searchInfo.productId" filterable clearable style="width:240px">
            <el-option v-for="p in productOpts" :key="p.ID" :label="`${p.code} - ${p.name}`" :value="p.ID" />
          </el-select>
        </el-form-item>
        <el-form-item label="仓库">
          <el-select v-model="searchInfo.warehouseId" clearable style="width:160px">
            <el-option v-for="w in warehouseOpts" :key="w.ID" :label="w.name" :value="w.ID" />
          </el-select>
        </el-form-item>
        <el-form-item><el-button type="primary" icon="search" @click="onSearch">查询</el-button><el-button icon="refresh" @click="onReset">重置</el-button></el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <el-table :data="tableData" border>
        <el-table-column prop="CreatedAt" label="时间" width="180">
          <template #default="{ row }">{{ row.CreatedAt ? new Date(row.CreatedAt).toLocaleString() : '' }}</template>
        </el-table-column>
        <el-table-column prop="bizType" label="业务类型" width="140" />
        <el-table-column prop="bizNo" label="业务单号" width="180" />
        <el-table-column label="方向" width="80">
          <template #default="{ row }">
            <el-tag :type="row.direction === 1 ? 'success' : 'danger'">{{ row.direction === 1 ? '入' : '出' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="商品" min-width="200"><template #default="{ row }">{{ row.product?.code }} - {{ row.product?.name }}</template></el-table-column>
        <el-table-column label="仓库" width="120"><template #default="{ row }">{{ row.warehouse?.name || '-' }}</template></el-table-column>
        <el-table-column prop="quantity" label="数量" width="100" />
        <el-table-column prop="unitCost" label="成本单价" width="100" />
        <el-table-column prop="qtyBefore" label="变动前" width="100" />
        <el-table-column prop="qtyAfter" label="变动后" width="100" />
        <el-table-column prop="remark" label="备注" />
      </el-table>
      <div class="gva-pagination">
        <el-pagination v-model:current-page="page" v-model:page-size="pageSize" :total="total" :page-sizes="[10,20,50,100]"
          layout="total, sizes, prev, pager, next, jumper" @current-change="getList" @size-change="getList" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { getStockFlow } from '@/api/inventory/stock'
import { getProductList } from '@/api/inventory/product'
import { getWarehouseAll } from '@/api/inventory/warehouse'

defineOptions({ name: 'InvStockFlow' })
const searchInfo = reactive({ productId: undefined, warehouseId: undefined })
const page = ref(1); const pageSize = ref(20); const total = ref(0); const tableData = ref([])
const productOpts = ref([]); const warehouseOpts = ref([])
const getList = async () => {
  const res = await getStockFlow({ page: page.value, pageSize: pageSize.value, ...searchInfo })
  if (res.code === 0) { tableData.value = res.data.list || []; total.value = res.data.total }
}
const onSearch = () => { page.value = 1; getList() }
const onReset = () => { searchInfo.productId = undefined; searchInfo.warehouseId = undefined; onSearch() }
onMounted(async () => {
  const [w, p] = await Promise.all([getWarehouseAll(), getProductList({ page: 1, pageSize: 9999 })])
  if (w.code === 0) warehouseOpts.value = w.data.list || []
  if (p.code === 0) productOpts.value = p.data.list || []
  getList()
})
</script>
