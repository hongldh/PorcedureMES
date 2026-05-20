import service from '@/utils/request'

// 仓库
const wh = '/inv/warehouse'
export const createWarehouse = (data) => service({ url: wh, method: 'post', data })
export const updateWarehouse = (data) => service({ url: wh, method: 'put', data })
export const deleteWarehouse = (data) => service({ url: wh, method: 'delete', data })
export const deleteWarehouseBatch = (data) => service({ url: wh + '/batch', method: 'delete', data })
export const findWarehouse = (params) => service({ url: wh, method: 'get', params })
export const getWarehouseList = (params) => service({ url: wh + '/list', method: 'get', params })
export const getWarehouseAll = () => service({ url: wh + '/all', method: 'get' })

// 货位
const lc = '/inv/location'
export const createLocation = (data) => service({ url: lc, method: 'post', data })
export const updateLocation = (data) => service({ url: lc, method: 'put', data })
export const deleteLocation = (data) => service({ url: lc, method: 'delete', data })
export const deleteLocationBatch = (data) => service({ url: lc + '/batch', method: 'delete', data })
export const findLocation = (params) => service({ url: lc, method: 'get', params })
export const getLocationList = (params) => service({ url: lc + '/list', method: 'get', params })
export const getLocationByWarehouse = (params) => service({ url: lc + '/byWarehouse', method: 'get', params })
