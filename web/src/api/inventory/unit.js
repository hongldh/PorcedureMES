import service from '@/utils/request'
const base = '/inv/unit'
export const createUnit = (data) => service({ url: base, method: 'post', data })
export const updateUnit = (data) => service({ url: base, method: 'put', data })
export const deleteUnit = (data) => service({ url: base, method: 'delete', data })
export const deleteUnitBatch = (data) => service({ url: base + '/batch', method: 'delete', data })
export const findUnit = (params) => service({ url: base, method: 'get', params })
export const getUnitList = (params) => service({ url: base + '/list', method: 'get', params })
export const getUnitAll = () => service({ url: base + '/all', method: 'get' })
