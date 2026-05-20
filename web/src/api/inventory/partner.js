import service from '@/utils/request'

// 供应商
const sup = '/inv/supplier'
export const createSupplier = (data) => service({ url: sup, method: 'post', data })
export const updateSupplier = (data) => service({ url: sup, method: 'put', data })
export const deleteSupplier = (data) => service({ url: sup, method: 'delete', data })
export const deleteSupplierBatch = (data) => service({ url: sup + '/batch', method: 'delete', data })
export const findSupplier = (params) => service({ url: sup, method: 'get', params })
export const getSupplierList = (params) => service({ url: sup + '/list', method: 'get', params })
export const getSupplierAll = () => service({ url: sup + '/all', method: 'get' })

// 客户（注意：与 GVA 自带 customer.js 不同，这是进销存独立的客户表）
const cus = '/inv/customer'
export const createInvCustomer = (data) => service({ url: cus, method: 'post', data })
export const updateInvCustomer = (data) => service({ url: cus, method: 'put', data })
export const deleteInvCustomer = (data) => service({ url: cus, method: 'delete', data })
export const deleteInvCustomerBatch = (data) => service({ url: cus + '/batch', method: 'delete', data })
export const findInvCustomer = (params) => service({ url: cus, method: 'get', params })
export const getInvCustomerList = (params) => service({ url: cus + '/list', method: 'get', params })
export const getInvCustomerAll = () => service({ url: cus + '/all', method: 'get' })
