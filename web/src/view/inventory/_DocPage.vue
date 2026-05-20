<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo">
        <el-form-item label="单号"><el-input v-model="searchInfo.docNo" clearable /></el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchInfo.status" placeholder="全部" clearable style="width:120px">
            <el-option label="草稿" :value="1" />
            <el-option label="已执行" :value="3" />
            <el-option label="作废" :value="4" />
          </el-select>
        </el-form-item>
        <el-form-item label="日期">
          <el-date-picker v-model="dateRange" type="daterange" value-format="YYYY-MM-DD" start-placeholder="开始" end-placeholder="结束" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSearch">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="onOpenAdd">新建{{ title }}</el-button>
      </div>
      <el-table :data="tableData" border>
        <el-table-column prop="docNo" label="单号" width="180" />
        <el-table-column prop="bizDate" label="日期" width="120" />
        <el-table-column :label="partyLabel" width="160">
          <template #default="{ row }">{{ row[partyRel]?.name || '-' }}</template>
        </el-table-column>
        <el-table-column label="仓库" width="140">
          <template #default="{ row }">{{ row.warehouse?.name || '-' }}</template>
        </el-table-column>
        <el-table-column prop="totalAmount" label="总金额" width="120" />
        <el-table-column label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="statusTag(row.status)">{{ statusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" />
        <el-table-column label="操作" width="240" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="onOpenView(row)">查看</el-button>
            <el-button v-if="row.status === 1" type="primary" link @click="onOpenEdit(row)">编辑</el-button>
            <el-button v-if="row.status === 1" type="success" link @click="onConfirm(row, confirmLabel)">{{ confirmLabel }}</el-button>
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
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="960px" :close-on-click-modal="false">
      <el-form ref="formRef" :model="formData" :rules="rules" label-width="90px" :disabled="readonly">
        <el-row :gutter="16">
          <el-col :span="8"><el-form-item label="单号"><el-input v-model="formData.docNo" disabled placeholder="保存后生成" /></el-form-item></el-col>
          <el-col :span="8"><el-form-item label="单据日期" prop="bizDate"><el-date-picker v-model="formData.bizDate" type="date" value-format="YYYY-MM-DD" style="width:100%" /></el-form-item></el-col>
          <el-col :span="8"><el-form-item :label="partyLabel" :prop="partyKey">
            <el-select v-model="formData[partyKey]" filterable style="width:100%">
              <el-option v-for="o in partyOpts" :key="o.ID" :label="o.name" :value="o.ID" />
            </el-select>
          </el-form-item></el-col>
          <el-col :span="8"><el-form-item label="仓库" prop="warehouseId">
            <el-select v-model="formData.warehouseId" style="width:100%" @change="onWhChange">
              <el-option v-for="w in warehouseOpts" :key="w.ID" :label="w.name" :value="w.ID" />
            </el-select>
          </el-form-item></el-col>
          <el-col :span="16"><el-form-item label="备注"><el-input v-model="formData.remark" /></el-form-item></el-col>
        </el-row>
        <el-divider>明细</el-divider>
        <el-table :data="formData.items" border size="small">
          <el-table-column type="index" label="#" width="50" />
          <el-table-column label="商品" min-width="220">
            <template #default="{ row }">
              <el-select v-model="row.productId" filterable :disabled="readonly" style="width:100%" @change="(v) => onProductChange(row, v)">
                <el-option v-for="p in productOpts" :key="p.ID" :label="`${p.code} - ${p.name}`" :value="p.ID" />
              </el-select>
            </template>
          </el-table-column>
          <el-table-column label="货位" width="160">
            <template #default="{ row }">
              <el-select v-model="row.locationId" clearable :disabled="readonly" style="width:100%">
                <el-option v-for="l in locationOpts" :key="l.ID" :label="l.code" :value="l.ID" />
              </el-select>
            </template>
          </el-table-column>
          <el-table-column label="数量" width="120"><template #default="{ row }"><el-input-number v-model="row.quantity" :min="0" :precision="4" :disabled="readonly" controls-position="right" style="width:100%" @change="calcRow(row)" /></template></el-table-column>
          <el-table-column :label="priceLabel" width="120"><template #default="{ row }"><el-input-number v-model="row.unitPrice" :min="0" :precision="4" :disabled="readonly" controls-position="right" style="width:100%" @change="calcRow(row)" /></template></el-table-column>
          <el-table-column label="金额" width="120"><template #default="{ row }">{{ (row.amount || 0).toFixed(2) }}</template></el-table-column>
          <el-table-column label="备注" min-width="120"><template #default="{ row }"><el-input v-model="row.remark" :disabled="readonly" /></template></el-table-column>
          <el-table-column v-if="!readonly" label="操作" width="70" fixed="right"><template #default="{ $index }"><el-button type="danger" link @click="removeRow($index)">删除</el-button></template></el-table-column>
        </el-table>
        <el-button v-if="!readonly" type="primary" link icon="plus" @click="addRow" style="margin-top:8px">新增明细行</el-button>
        <div style="margin-top:12px;text-align:right">
          <span style="margin-right:24px">总数量：<b>{{ totalQuantity.toFixed(4) }}</b></span>
          <span>总金额：<b>{{ totalAmount.toFixed(2) }}</b></span>
        </div>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">关闭</el-button>
        <el-button v-if="!readonly" type="primary" @click="onSubmit">保存草稿</el-button>
        <el-button v-if="!readonly && formData.ID" type="success" @click="onSubmitAndConfirm">保存并{{ confirmLabel }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { getProductList } from '@/api/inventory/product'
import { getWarehouseAll, getLocationByWarehouse } from '@/api/inventory/warehouse'
import { useDocPage } from './_useDocPage.js'

const props = defineProps({
  title: { type: String, required: true },
  api: { type: Object, required: true },
  partyKey: { type: String, required: true },
  partyLabel: { type: String, required: true },
  partyRel: { type: String, required: true },
  partyApi: { type: Function, required: true },
  priceLabel: { type: String, default: '单价' },
  confirmLabel: { type: String, default: '确认' },
  defaultPriceField: { type: String, default: 'purchasePrice' }
})

const ctx = useDocPage(props.api, { partyKey: props.partyKey, partyApi: props.partyApi, defaultPriceField: props.defaultPriceField })
const {
  searchInfo, dateRange, page, pageSize, total, tableData,
  statusText, statusTag, getList, onSearch, onReset,
  partyOpts, warehouseOpts, productOpts, locationOpts,
  dialogVisible, dialogType, formRef, formData, readonly,
  calcRow, totalQuantity, totalAmount, onProductChange, addRow, removeRow,
  openDialog, onConfirm, onVoid, onDelete
} = ctx

const dialogTitle = computed(() => ({ add: `新建${props.title}`, edit: `编辑${props.title}`, view: `查看${props.title}` }[dialogType.value]))
const { partyKey, partyLabel, partyRel, priceLabel, confirmLabel, title } = props

const rules = {
  bizDate: [{ required: true, message: '必填' }],
  [props.partyKey]: [{ required: true, message: '必选' }],
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
const onOpenAdd = () => openDialog('add', null, onWhChange)
const onOpenEdit = (row) => openDialog('edit', row, loadLocations)
const onOpenView = (row) => openDialog('view', row, loadLocations)
const onSubmit = () => ctx.onSubmit()
const onSubmitAndConfirm = () => ctx.onSubmitAndConfirm()

onMounted(async () => {
  ctx.loadParty()
  const [w, p] = await Promise.all([getWarehouseAll(), getProductList({ page: 1, pageSize: 9999 })])
  if (w.code === 0) warehouseOpts.value = w.data.list || []
  if (p.code === 0) productOpts.value = p.data.list || []
  getList()
})
</script>
