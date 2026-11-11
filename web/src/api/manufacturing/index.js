import service from '@/utils/request'

// ===== BOM =====
export const bomApi = {
  create:  (data)   => service({ url: '/mfg/bom', method: 'post',   data }),
  update:  (data)   => service({ url: '/mfg/bom', method: 'put',    data }),
  remove:  (data)   => service({ url: '/mfg/bom', method: 'delete', data }),
  setStatus: (data) => service({ url: '/mfg/bom/status', method: 'post', data }),
  findOne: (params) => service({ url: '/mfg/bom', method: 'get', params }),
  getList: (params) => service({ url: '/mfg/bom/list', method: 'get', params }),
  getActiveByProduct: (params) => service({ url: '/mfg/bom/activeByProduct', method: 'get', params })
}

// ===== 工单 =====
export const workOrderApi = {
  create:  (data) => service({ url: '/mfg/workOrder', method: 'post',   data }),
  update:  (data) => service({ url: '/mfg/workOrder', method: 'put',    data }),
  remove:  (data) => service({ url: '/mfg/workOrder', method: 'delete', data }),
  release: (data) => service({ url: '/mfg/workOrder/release', method: 'post', data }),
  close:   (data) => service({ url: '/mfg/workOrder/close',   method: 'post', data }),
  void:    (data) => service({ url: '/mfg/workOrder/void',    method: 'post', data }),
  findOne: (params) => service({ url: '/mfg/workOrder', method: 'get', params }),
  getList: (params) => service({ url: '/mfg/workOrder/list', method: 'get', params })
}

// ===== 领料单 =====
export const issueApi = {
  create:  (data) => service({ url: '/mfg/issue', method: 'post',   data }),
  update:  (data) => service({ url: '/mfg/issue', method: 'put',    data }),
  remove:  (data) => service({ url: '/mfg/issue', method: 'delete', data }),
  confirm: (data) => service({ url: '/mfg/issue/confirm', method: 'post', data }),
  void:    (data) => service({ url: '/mfg/issue/void',    method: 'post', data }),
  findOne: (params) => service({ url: '/mfg/issue', method: 'get', params }),
  getList: (params) => service({ url: '/mfg/issue/list', method: 'get', params })
}

// ===== 完工入库单 =====
export const receiptApi = {
  create:  (data) => service({ url: '/mfg/receipt', method: 'post',   data }),
  update:  (data) => service({ url: '/mfg/receipt', method: 'put',    data }),
  remove:  (data) => service({ url: '/mfg/receipt', method: 'delete', data }),
  confirm: (data) => service({ url: '/mfg/receipt/confirm', method: 'post', data }),
  void:    (data) => service({ url: '/mfg/receipt/void',    method: 'post', data }),
  findOne: (params) => service({ url: '/mfg/receipt', method: 'get', params }),
  getList: (params) => service({ url: '/mfg/receipt/list', method: 'get', params })
}

// ============================================================
// 阶段 5a：工艺主数据
// ============================================================

// ===== 工位 =====
export const workstationApi = {
  create:  (data) => service({ url: '/mfg/workstation', method: 'post',   data }),
  update:  (data) => service({ url: '/mfg/workstation', method: 'put',    data }),
  remove:  (data) => service({ url: '/mfg/workstation', method: 'delete', data }),
  findOne: (params) => service({ url: '/mfg/workstation', method: 'get', params }),
  getList: (params) => service({ url: '/mfg/workstation/list', method: 'get', params }),
  getAll:  () => service({ url: '/mfg/workstation/all', method: 'get' })
}

// ===== 班组 =====
export const teamApi = {
  create:  (data) => service({ url: '/mfg/team', method: 'post',   data }),
  update:  (data) => service({ url: '/mfg/team', method: 'put',    data }),
  remove:  (data) => service({ url: '/mfg/team', method: 'delete', data }),
  findOne: (params) => service({ url: '/mfg/team', method: 'get', params }),
  getList: (params) => service({ url: '/mfg/team/list', method: 'get', params })
}

// ===== 工艺路线 =====
export const routingApi = {
  create:  (data) => service({ url: '/mfg/routing', method: 'post',   data }),
  update:  (data) => service({ url: '/mfg/routing', method: 'put',    data }),
  remove:  (data) => service({ url: '/mfg/routing', method: 'delete', data }),
  setStatus: (data) => service({ url: '/mfg/routing/status', method: 'post', data }),
  findOne: (params) => service({ url: '/mfg/routing', method: 'get', params }),
  getList: (params) => service({ url: '/mfg/routing/list', method: 'get', params }),
  getActiveByProduct: (params) => service({ url: '/mfg/routing/activeByProduct', method: 'get', params })
}

// ============================================================
// 阶段 5b：派工 / 报工
// ============================================================

// ===== 工单工序 + 派工 =====
export const dispatchApi = {
  getOps:    (params) => service({ url: '/mfg/workOrder/ops', method: 'get', params }),
  dispatch:  (data)   => service({ url: '/mfg/op/dispatch', method: 'post', data })
}

// ===== 报工 =====
export const workReportApi = {
  create:   (data)   => service({ url: '/mfg/report', method: 'post',   data }),
  void:     (data)   => service({ url: '/mfg/report/void', method: 'post', data }),
  getList:  (params) => service({ url: '/mfg/report/list', method: 'get', params }),
  myOps:    ()       => service({ url: '/mfg/report/myOps', method: 'get' })
}
