import service from '@/utils/request'

// === 通用单据 API 工厂 ===
// 所有业务单据共享 7 个端点：create/update/delete/confirm/void/findOne/getList
// 这里用工厂函数减少重复（学习要点：JS 高阶函数）
const makeDocApi = (base) => ({
  create:  (data)   => service({ url: base, method: 'post',   data }),
  update:  (data)   => service({ url: base, method: 'put',    data }),
  remove:  (data)   => service({ url: base, method: 'delete', data }),
  confirm: (data)   => service({ url: base + '/confirm', method: 'post', data }),
  void:    (data)   => service({ url: base + '/void',    method: 'post', data }),
  findOne: (params) => service({ url: base, method: 'get', params }),
  getList: (params) => service({ url: base + '/list', method: 'get', params })
})

export const purchaseOrderApi  = makeDocApi('/inv/purchaseOrder')
export const purchaseReturnApi = makeDocApi('/inv/purchaseReturn')
export const salesOrderApi     = makeDocApi('/inv/salesOrder')
export const salesReturnApi    = makeDocApi('/inv/salesReturn')
export const transferApi       = makeDocApi('/inv/transfer')
export const stocktakeApi      = makeDocApi('/inv/stocktake')
