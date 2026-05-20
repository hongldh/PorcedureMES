/**
 * 进销存 - 商品分类 API
 *
 * === 学习要点 ===
 * GVA 前端 API 文件结构非常简单：
 *   - 引入 service（已封装的 axios 实例，自动加 token + 错误处理）
 *   - 每个后端接口对应一个 export 函数
 *   - 列表/查询用 params（GET），写操作用 data（POST/PUT/DELETE body）
 */
import service from '@/utils/request'

const base = '/inv/category'

export const createCategory = (data) => service({ url: base, method: 'post', data })
export const updateCategory = (data) => service({ url: base, method: 'put', data })
export const deleteCategory = (data) => service({ url: base, method: 'delete', data })
export const deleteCategoryBatch = (data) => service({ url: base + '/batch', method: 'delete', data })
export const findCategory = (params) => service({ url: base, method: 'get', params })
export const getCategoryList = (params) => service({ url: base + '/list', method: 'get', params })
export const getCategoryAll = () => service({ url: base + '/all', method: 'get' })
