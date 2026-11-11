<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo">
        <el-form-item label="关键字"><el-input v-model="searchInfo.keyword" placeholder="编码/名称" clearable /></el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchInfo.status" clearable placeholder="全部" style="width:120px">
            <el-option label="启用" :value="1" /><el-option label="停用" :value="2" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSearch">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="gva-table-box">
      <div class="gva-btn-list"><el-button type="primary" icon="plus" @click="openDialog('add')">新增班组</el-button></div>
      <el-table :data="tableData" border>
        <el-table-column prop="code" label="编码" width="180" />
        <el-table-column prop="name" label="名称" />
        <el-table-column label="班组长" width="160"><template #default="{ row }">{{ row.leader?.nickName || row.leader?.userName || '-' }}</template></el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }"><el-tag :type="row.status === 1 ? 'success' : 'info'">{{ row.status === 1 ? '启用' : '停用' }}</el-tag></template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" />
        <el-table-column label="操作" width="160" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="openDialog('edit', row)">编辑</el-button>
            <el-button type="danger" link @click="onDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination v-model:current-page="page" v-model:page-size="pageSize" :total="total" :page-sizes="[10,20,50,100]"
          layout="total, sizes, prev, pager, next, jumper" @current-change="getList" @size-change="getList" />
      </div>
    </div>

    <el-dialog v-model="dialogVisible" :title="dialogType === 'add' ? '新增班组' : '编辑班组'" width="640px">
      <el-form ref="formRef" :model="formData" :rules="rules" label-width="100px">
        <el-form-item label="编码" prop="code"><el-input v-model="formData.code" placeholder="如 TEAM01" /></el-form-item>
        <el-form-item label="名称" prop="name"><el-input v-model="formData.name" /></el-form-item>
        <el-form-item label="班组长">
          <el-select v-model="formData.leaderId" filterable clearable style="width:100%">
            <el-option v-for="u in userOpts" :key="u.ID" :label="`${u.nickName || u.userName}（${u.userName}）`" :value="u.ID" />
          </el-select>
        </el-form-item>
        <el-form-item label="组员">
          <el-select v-model="memberIds" multiple filterable style="width:100%">
            <el-option v-for="u in userOpts" :key="u.ID" :label="`${u.nickName || u.userName}（${u.userName}）`" :value="u.ID" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="formData.status">
            <el-radio :value="1">启用</el-radio><el-radio :value="2">停用</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="备注"><el-input v-model="formData.remark" type="textarea" :rows="2" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">关闭</el-button>
        <el-button type="primary" @click="onSubmit">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { teamApi } from '@/api/manufacturing'
import { getUserList } from '@/api/user'

defineOptions({ name: 'MfgTeam' })

const searchInfo = reactive({ keyword: '', status: undefined })
const page = ref(1); const pageSize = ref(10); const total = ref(0); const tableData = ref([])
const userOpts = ref([])

const getList = async () => {
  const res = await teamApi.getList({ page: page.value, pageSize: pageSize.value, ...searchInfo })
  if (res.code === 0) { tableData.value = res.data.list || []; total.value = res.data.total }
}
const onSearch = () => { page.value = 1; getList() }
const onReset = () => { Object.assign(searchInfo, { keyword: '', status: undefined }); onSearch() }

const dialogVisible = ref(false); const dialogType = ref('add'); const formRef = ref(null)
const memberIds = ref([])
const defaultForm = () => ({ ID: 0, code: '', name: '', leaderId: undefined, status: 1, remark: '' })
const formData = reactive(defaultForm())
const rules = {
  code: [{ required: true, message: '必填' }],
  name: [{ required: true, message: '必填' }]
}

const openDialog = async (type, row) => {
  dialogType.value = type
  Object.assign(formData, defaultForm())
  memberIds.value = []
  if (row) {
    const res = await teamApi.findOne({ id: row.ID })
    if (res.code === 0) {
      const d = res.data.record || res.data
      Object.assign(formData, { ID: d.ID, code: d.code, name: d.name, leaderId: d.leaderId, status: d.status, remark: d.remark })
      memberIds.value = (d.members || []).map(m => m.userId)
    }
  }
  dialogVisible.value = true
}
const onSubmit = async () => {
  try { await formRef.value.validate() } catch { ElMessage.error('请补全必填项'); return }
  const payload = { ...formData, members: memberIds.value.map(uid => ({ userId: uid })) }
  const fn = formData.ID ? teamApi.update : teamApi.create
  const res = await fn(payload)
  if (res.code === 0) { ElMessage.success('保存成功'); dialogVisible.value = false; getList() }
}
const onDelete = (row) => {
  ElMessageBox.confirm(`删除班组 [${row.code}]?`, '提示', { type: 'warning' }).then(async () => {
    const res = await teamApi.remove({ id: row.ID })
    if (res.code === 0) { ElMessage.success('已删除'); getList() }
  }).catch(() => {})
}

onMounted(async () => {
  const u = await getUserList({ page: 1, pageSize: 9999 })
  if (u.code === 0) userOpts.value = u.data.list || []
  getList()
})
</script>
