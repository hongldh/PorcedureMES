import service from '@/utils/request'
const base = '/inv/product'
export const createProduct = (data) => service({ url: base, method: 'post', data })
export const updateProduct = (data) => service({ url: base, method: 'put', data })
export const deleteProduct = (data) => service({ url: base, method: 'delete', data })
export const deleteProductBatch = (data) => service({ url: base + '/batch', method: 'delete', data })
export const findProduct = (params) => service({ url: base, method: 'get', params })
export const getProductList = (params) => service({ url: base + '/list', method: 'get', params })
export const getProductAll = () => service({ url: base + '/all', method: 'get' })
