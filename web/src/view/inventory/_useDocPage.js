// 通用单据页 composable：列表搜索 + 草稿编辑 + 确认/作废/删除
// 学习要点：Vue3 Composition API + 高阶函数复用
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'

export function useDocPage(api, options = {}) {
  const { partyKey, partyApi, defaultPriceField = 'purchasePrice' } = options
  // partyKey: 'supplierId' | 'customerId' | null
  // partyApi: () => Promise<list>

  const searchInfo = reactive({ docNo: '', status: undefined })
  const dateRange = ref([])
  const page = ref(1); const pageSize = ref(10); const total = ref(0)
  const tableData = ref([])

  const statusText = (s) => ({ 1: '草稿', 3: '已执行', 4: '作废' }[s] || '-')
  const statusTag = (s) => ({ 1: 'info', 3: 'success', 4: 'danger' }[s] || 'info')

  const getList = async () => {
    const params = { page: page.value, pageSize: pageSize.value, ...searchInfo }
    if (dateRange.value?.length === 2) { params.dateStart = dateRange.value[0]; params.dateEnd = dateRange.value[1] }
    const res = await api.getList(params)
    if (res.code === 0) { tableData.value = res.data.list || []; total.value = res.data.total }
  }
  const onSearch = () => { page.value = 1; getList() }
  const onReset = () => { searchInfo.docNo = ''; searchInfo.status = undefined; dateRange.value = []; onSearch() }

  // === 选项 ===
  const partyOpts = ref([]); const warehouseOpts = ref([])
  const productOpts = ref([]); const locationOpts = ref([])

  const dialogVisible = ref(false); const dialogType = ref('add')
  const formRef = ref(null)   // ★ 必须在 composable 内部声明，onSubmit 闭包才能引用
  const readonly = computed(() => dialogType.value === 'view')

  const newRow = () => ({ productId: undefined, locationId: undefined, quantity: 0, unitPrice: 0, amount: 0, remark: '' })
  const buildDefault = () => {
    const f = { ID: 0, docNo: '', bizDate: new Date().toISOString().slice(0, 10),
      warehouseId: undefined, remark: '', items: [newRow()] }
    if (partyKey) f[partyKey] = undefined
    return f
  }
  const formData = reactive(buildDefault())

  const calcRow = (row) => { row.amount = Number(((row.quantity || 0) * (row.unitPrice || 0)).toFixed(4)) }
  const totalQuantity = computed(() => formData.items.reduce((s, r) => s + Number(r.quantity || 0), 0))
  const totalAmount = computed(() => formData.items.reduce((s, r) => s + Number(r.amount || 0), 0))
  const onProductChange = (row, pid) => {
    const p = productOpts.value.find(x => x.ID === pid)
    if (p && !row.unitPrice) { row.unitPrice = p[defaultPriceField] || 0; calcRow(row) }
  }
  const addRow = () => formData.items.push(newRow())
  const removeRow = (i) => formData.items.splice(i, 1)

  const openDialog = async (type, row, onWarehouseChange) => {
    dialogType.value = type
    Object.assign(formData, buildDefault()); formData.items = [newRow()]
    if (row) {
      const res = await api.findOne({ id: row.ID })
      if (res.code === 0) {
        const d = res.data.record || res.data.redata || res.data
        Object.assign(formData, {
          ID: d.ID, docNo: d.docNo, bizDate: d.bizDate,
          warehouseId: d.warehouseId, remark: d.remark,
          items: (d.items || []).map(it => ({ ...it, amount: Number(it.quantity) * Number(it.unitPrice) }))
        })
        if (partyKey) formData[partyKey] = d[partyKey]
        if (d.warehouseId && onWarehouseChange) await onWarehouseChange(d.warehouseId)
      }
    }
    dialogVisible.value = true
  }

  const buildPayload = () => {
    const p = { ID: formData.ID, bizDate: formData.bizDate, warehouseId: formData.warehouseId, remark: formData.remark,
      items: formData.items.filter(r => r.productId && r.quantity > 0).map(r => ({
        ID: r.ID || 0, productId: r.productId, locationId: r.locationId || 0,
        quantity: r.quantity, unitPrice: r.unitPrice, remark: r.remark
      })) }
    if (partyKey) p[partyKey] = formData[partyKey]
    return p
  }
  const onSubmit = async () => {
    try {
      await formRef.value.validate()
    } catch (e) {
      ElMessage.error('请补全必填项')
      return false
    }
    const payload = buildPayload()
    if (!payload.items.length) { ElMessage.error('请至少添加一条明细'); return false }
    const fn = formData.ID ? api.update : api.create
    const res = await fn(payload)
    if (res.code === 0) {
      ElMessage.success('保存成功')
      // create 接口若返回 id，回填
      const newId = res.data?.ID || res.data?.id || res.data?.redata?.ID
      if (newId && !formData.ID) formData.ID = newId
      getList(); return true
    }
    return false
  }
  const onSubmitAndConfirm = async () => {
    const ok = await onSubmit()
    if (ok && formData.ID) {
      const res = await api.confirm({ id: formData.ID })
      if (res.code === 0) { ElMessage.success('已确认'); dialogVisible.value = false; getList() }
    }
  }
  const onConfirm = (row, label = '确认') => {
    ElMessageBox.confirm(`${label} [${row.docNo}]? 库存将变更。`, '提示', { type: 'warning' }).then(async () => {
      const res = await api.confirm({ id: row.ID })
      if (res.code === 0) { ElMessage.success('成功'); getList() }
    }).catch(() => {})
  }
  const onVoid = (row) => {
    ElMessageBox.confirm(`作废 [${row.docNo}]? 将生成反向流水。`, '提示', { type: 'warning' }).then(async () => {
      const res = await api.void({ id: row.ID })
      if (res.code === 0) { ElMessage.success('已作废'); getList() }
    }).catch(() => {})
  }
  const onDelete = (row) => {
    ElMessageBox.confirm(`删除草稿 [${row.docNo}]?`, '提示', { type: 'warning' }).then(async () => {
      const res = await api.remove({ id: row.ID })
      if (res.code === 0) { ElMessage.success('已删除'); getList() }
    }).catch(() => {})
  }

  const loadParty = async () => {
    if (partyApi) { const r = await partyApi(); if (r.code === 0) partyOpts.value = r.data.list || [] }
  }

  return {
    searchInfo, dateRange, page, pageSize, total, tableData, statusText, statusTag,
    getList, onSearch, onReset,
    partyOpts, warehouseOpts, productOpts, locationOpts, loadParty,
    dialogVisible, dialogType, formRef, formData, readonly,
    // ★ formRef 现已在 return 中导出，外层 _DocPage.vue 只需将 <el-form ref="..."> 绑到它
    calcRow, totalQuantity, totalAmount, onProductChange, addRow, removeRow,
    openDialog, onSubmit, onSubmitAndConfirm, onConfirm, onVoid, onDelete
  }
}
