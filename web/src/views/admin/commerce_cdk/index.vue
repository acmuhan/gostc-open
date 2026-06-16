<template>
  <div class="commerce-page">
    <n-card title="生成 CDK">
      <n-form inline :model="form">
        <n-form-item label="类型">
          <n-select v-model:value="form.type" :options="typeOptions" style="width:140px" />
        </n-form-item>
        <n-form-item label="面值">
          <n-input-number v-model:value="form.value" :min="0.01" />
        </n-form-item>
        <n-form-item label="数量">
          <n-input-number v-model:value="form.count" :min="1" :max="500" />
        </n-form-item>
        <n-form-item label="备注">
          <n-input v-model:value="form.remark" placeholder="可选" style="width:160px" />
        </n-form-item>
        <n-form-item>
          <n-button type="primary" :loading="creating" @click="create">生成</n-button>
        </n-form-item>
      </n-form>
      <template v-if="codes.length">
        <n-input type="textarea" :value="codes.join('\n')" :autosize="{minRows:3,maxRows:8}" readonly />
        <n-space style="margin-top:8px">
          <n-button size="small" @click="copyCodes">复制全部</n-button>
          <n-button size="small" @click="exportCodes">导出 TXT</n-button>
        </n-space>
      </template>
    </n-card>
    <n-card title="CDK 列表">
      <template #header-extra>
        <n-space>
          <n-select v-model:value="filter.status" :options="statusFilterOptions" placeholder="状态筛选" clearable style="width:130px" @update:value="resetAndLoad" />
          <n-input v-model:value="filter.batchNo" placeholder="批次号" clearable style="width:160px" @clear="resetAndLoad" @keyup.enter="resetAndLoad" />
        </n-space>
      </template>
      <n-space v-if="checkedKeys.length" style="margin-bottom:12px">
        <n-tag>已选 {{ checkedKeys.length }} 项</n-tag>
        <n-button size="small" type="warning" @click="batchDisable">批量禁用</n-button>
        <n-button size="small" type="success" @click="batchEnable">批量启用</n-button>
        <n-button size="small" type="error" @click="batchDelete">批量删除</n-button>
        <n-button size="small" @click="showRemarkModal = true">修改备注</n-button>
      </n-space>
      <n-data-table
        :columns="columns"
        :data="list"
        :loading="loading"
        :row-key="r => r.cdk_code"
        v-model:checked-row-keys="checkedKeys"
      />
      <n-pagination class="pager" v-model:page="page" :page-count="pageCount" @update:page="load" />
    </n-card>

    <n-modal v-model:show="showRemarkModal" preset="dialog" title="修改备注" positive-text="确认" negative-text="取消" @positive-click="submitRemark">
      <n-input v-model:value="remarkInput" type="textarea" placeholder="输入新备注" :autosize="{minRows:2,maxRows:4}" />
    </n-modal>
  </div>
</template>
<script setup>
import { computed, h, onMounted, ref } from 'vue'
import { NButton, NTag } from 'naive-ui'
import {
  apiAdminCommerceCdkCreate, apiAdminCommerceCdkPage, apiAdminCommerceCdkDisable,
  apiAdminCommerceCdkBatchDisable, apiAdminCommerceCdkBatchEnable,
  apiAdminCommerceCdkBatchDelete, apiAdminCommerceCdkUpdateRemark,
} from '../../../api/admin/commerce_cdk.js'

const message = window.$message
const dialog = window.$dialog

const form = ref({ type: 'balance', value: 10, count: 1, remark: '' })
const typeOptions = [{ label: '余额', value: 'balance' }, { label: '积分', value: 'points' }]
const statusFilterOptions = [
  { label: '未使用', value: 1 },
  { label: '已使用', value: 2 },
  { label: '已禁用', value: 3 },
]
const codes = ref([])
const list = ref([])
const page = ref(1)
const total = ref(0)
const loading = ref(false)
const creating = ref(false)
const checkedKeys = ref([])
const filter = ref({ status: null, batchNo: '' })
const showRemarkModal = ref(false)
const remarkInput = ref('')
const pageCount = computed(() => Math.max(1, Math.ceil(total.value / 10)))

const columns = [
  { type: 'selection' },
  { title: '兑换码', key: 'cdk_code' },
  { title: '类型', key: 'type', width: 60, render: r => r.type === 'balance' ? '余额' : '积分' },
  { title: '面值', key: 'value', width: 80 },
  { title: '批次', key: 'batch_no' },
  { title: '状态', key: 'status', width: 80, render: r => h(NTag, { type: r.status === 1 ? 'info' : r.status === 2 ? 'success' : 'default', size: 'small' }, () => r.status === 1 ? '未使用' : r.status === 2 ? '已使用' : '已禁用') },
  { title: '使用用户', key: 'user_code', ellipsis: { tooltip: true } },
  { title: '备注', key: 'remark', ellipsis: { tooltip: true } },
  { title: '操作', key: 'actions', width: 80, render: r => r.status === 1 ? h(NButton, { size: 'small', type: 'warning', secondary: true, onClick: () => disable(r.cdk_code) }, () => '禁用') : r.status === 3 ? h(NButton, { size: 'small', type: 'success', secondary: true, onClick: () => enable(r.cdk_code) }, () => '启用') : null },
]

function copyCodes() {
  navigator.clipboard.writeText(codes.value.join('\n')).then(() => message.success('已复制到剪贴板')).catch(() => message.error('复制失败'))
}

function exportCodes() {
  const blob = new Blob([codes.value.join('\n')], { type: 'text/plain' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `cdk_${Date.now()}.txt`
  a.click()
  URL.revokeObjectURL(url)
}

async function create() {
  creating.value = true
  try {
    const res = await apiAdminCommerceCdkCreate(form.value)
    codes.value = res.data || []
    message.success('生成成功')
    page.value = 1
    await load()
  } catch (_) {}
  creating.value = false
}

async function load() {
  loading.value = true
  try {
    const params = { page: page.value, size: 10 }
    if (filter.value.status) params.status = filter.value.status
    if (filter.value.batchNo) params.batchNo = filter.value.batchNo
    const res = await apiAdminCommerceCdkPage(params)
    list.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch (_) {}
  loading.value = false
  checkedKeys.value = []
}

function resetAndLoad() { page.value = 1; load() }

async function disable(code) {
  try {
    await apiAdminCommerceCdkDisable({ code })
    message.success('已禁用')
    await load()
  } catch (_) {}
}

async function enable(code) {
  try {
    await apiAdminCommerceCdkBatchEnable({ codes: [code] })
    message.success('已启用')
    await load()
  } catch (_) {}
}

async function batchDisable() {
  dialog.warning({ title: '确认批量禁用', content: `确认禁用选中的 ${checkedKeys.value.length} 个CDK？`, positiveText: '确认', negativeText: '取消', onPositiveClick: async () => {
    try { await apiAdminCommerceCdkBatchDisable({ codes: checkedKeys.value }); message.success('批量禁用成功'); await load() } catch (_) {}
  }})
}

async function batchEnable() {
  dialog.warning({ title: '确认批量启用', content: `确认启用选中的 ${checkedKeys.value.length} 个CDK？`, positiveText: '确认', negativeText: '取消', onPositiveClick: async () => {
    try { await apiAdminCommerceCdkBatchEnable({ codes: checkedKeys.value }); message.success('批量启用成功'); await load() } catch (_) {}
  }})
}

async function batchDelete() {
  dialog.error({ title: '确认批量删除', content: `确认删除选中的 ${checkedKeys.value.length} 个CDK？已使用的CDK不会被删除。`, positiveText: '确认删除', negativeText: '取消', onPositiveClick: async () => {
    try { await apiAdminCommerceCdkBatchDelete({ codes: checkedKeys.value }); message.success('批量删除成功'); await load() } catch (_) {}
  }})
}

async function submitRemark() {
  if (!checkedKeys.value.length) return
  try {
    await apiAdminCommerceCdkUpdateRemark({ codes: checkedKeys.value, remark: remarkInput.value })
    message.success('备注已更新')
    remarkInput.value = ''
    showRemarkModal.value = false
    await load()
  } catch (_) {}
}

onMounted(load)
</script>
<style scoped>
.commerce-page { display: flex; flex-direction: column; gap: 16px }
.pager { justify-content: flex-end; margin-top: 16px }
</style>
