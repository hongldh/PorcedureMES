import service from '@/utils/request'

// 库存查询：余额、流水
export const getStockBalance = (params) => service({ url: '/inv/stock/balance', method: 'get', params })
export const getStockFlow    = (params) => service({ url: '/inv/stock/flow',    method: 'get', params })
