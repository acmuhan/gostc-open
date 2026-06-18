<template>
  <div class="commerce-page">
    <n-card title="操作日志">
      <template #header-extra>
        <n-space>
          <n-input v-model:value="filter.module" placeholder="模块" clearable style="width:140px" @clear="resetAndLoad" @keyup.enter="resetAndLoad" />
          <n-input v-model:value="filter.action" placeholder="操作" clearable style="width:120px" @clear="resetAndLoad" @keyup.enter="resetAndLoad" />
          <n-input v-model:value="filter.userCode" placeholder="操作人编号" clearable style="width:160px" @clear="resetAndLoad" @keyup.enter="resetAndLoad" />
          <n-button type="primary" @click="resetAndLoad">搜索</n-button>
        </n-space>
      </template>
      <n-data-table :columns="columns" :data="list" :loading="loading" :row-key="r => r.id" />
      <n-pagination class="pager" v-model:page="page" :page-count="pageCount" @update:page="load" />
    </n-card>
  </div>
</template>
<script setup>
import { computed, h, onMounted, ref } from 'vue'
import { NTag } from 'naive-ui'
import { apiAdminSystemAuditLogPage } from '../../../api/admin/system_audit_log.js'

const filter = ref({ module: '', action: '', userCode: '' })
const list = ref([])
const page = ref(1)
const size = 20
const total = ref(0)
const loading = ref(false)
const pageCount = computed(() => Math.ceil(total.value / size) || 1)

const columns = [
  { title: 'ID', key: 'id', width: 60 },
  { title: '操作人', key: 'user_code', width: 120, ellipsis: true },
  { title: '模块', key: 'module', width: 140 },
  { title: '操作', key: 'action', width: 100 },
  { title: '请求路径', key: 'path', width: 220, ellipsis: true },
  {
    title: '结果', key: 'resp_code', width: 80,
    render: (row) => h(NTag, { type: row.resp_code === 0 ? 'success' : 'error', size: 'small' }, () => row.resp_code === 0 ? '成功' : '失败')
  },
  { title: 'IP', key: 'ip', width: 130 },
  { title: '耗时(ms)', key: 'ms', width: 80 },
  {
    title: '时间', key: 'created_at', width: 170,
    render: (row) => row.created_at ? new Date(row.created_at).toLocaleString() : '-'
  },
]

async function load() {
  loading.value = true
  try {
    const res = await apiAdminSystemAuditLogPage({
      page: page.value, size,
      userCode: filter.value.userCode || '',
      module: filter.value.module || '',
      action: filter.value.action || '',
    })
    list.value = res.data?.list || []
    total.value = res.data?.total || 0
  } finally {
    loading.value = false
  }
}

function resetAndLoad() {
  page.value = 1
  load()
}

onMounted(load)
</script>
<style scoped>
.pager { margin-top: 16px; display: flex; justify-content: flex-end; }
@media (max-width: 520px) {
  :deep(.n-data-table) { overflow-x: auto }
  :deep(.n-card-header) { flex-wrap: wrap; gap: 8px }
  :deep(.n-card__content) { padding: 12px !important }
  :deep(.n-space) { flex-wrap: wrap !important }
}
</style>
