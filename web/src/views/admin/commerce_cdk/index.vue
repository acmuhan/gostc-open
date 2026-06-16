<template>
  <div class="commerce-page">
    <n-card title="生成 CDK">
      <n-form inline :model="form">
        <n-form-item label="类型">
          <n-select v-model:value="form.type" :options="typeOptions" style="width:140px" />
        </n-form-item>
        <n-form-item label="面值">
          <n-input-number v-model:value="form.value" />
        </n-form-item>
        <n-form-item label="数量">
          <n-input-number v-model:value="form.count" />
        </n-form-item>
        <n-form-item>
          <n-button type="primary" :loading="creating" @click="create">生成</n-button>
        </n-form-item>
      </n-form>
      <n-input v-if="codes.length" type="textarea" :value="codes.join('\n')" :autosize="{minRows:3,maxRows:8}" />
    </n-card>
    <n-card title="CDK 列表">
      <n-data-table :columns="columns" :data="list" :loading="loading" />
      <n-pagination class="pager" v-model:page="page" :page-count="pageCount" @update:page="load" />
    </n-card>
  </div>
</template>
<script setup>
import { computed, h, onMounted, ref } from 'vue'
import { NButton, NTag, useMessage } from 'naive-ui'
import { apiAdminCommerceCdkCreate, apiAdminCommerceCdkDisable, apiAdminCommerceCdkPage } from '../../../api/admin/commerce_cdk.js'

const message = useMessage()
const form = ref({ type: 'balance', value: 10, count: 1 })
const typeOptions = [{ label: '余额', value: 'balance' }, { label: '积分', value: 'points' }]
const codes = ref([])
const list = ref([])
const page = ref(1)
const total = ref(0)
const loading = ref(false)
const creating = ref(false)
const pageCount = computed(() => Math.max(1, Math.ceil(total.value / 10)))

const columns = [
  { title: '兑换码', key: 'cdk_code' },
  { title: '类型', key: 'type' },
  { title: '面值', key: 'value' },
  { title: '批次', key: 'batch_no' },
  { title: '状态', key: 'status', render: r => h(NTag, { type: r.status === 1 ? 'info' : r.status === 2 ? 'success' : 'default' }, () => r.status === 1 ? '未使用' : r.status === 2 ? '已使用' : '已禁用') },
  { title: '使用用户', key: 'user_code' },
  { title: '操作', key: 'actions', render: r => r.status === 1 ? h(NButton, { size: 'small', onClick: () => disable(r.cdk_code) }, () => '禁用') : null },
]

async function create() {
  creating.value = true
  try {
    const res = await apiAdminCommerceCdkCreate(form.value)
    codes.value = res.data || []
    message.success('生成成功')
    await load()
  } catch (_) { /* 请求拦截器已处理错误提示 */ }
  creating.value = false
}

async function load() {
  loading.value = true
  try {
    const res = await apiAdminCommerceCdkPage({ page: page.value, size: 10 })
    list.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch (_) { /* 请求拦截器已处理错误提示 */ }
  loading.value = false
}

async function disable(code) {
  try {
    await apiAdminCommerceCdkDisable({ code })
    message.success('已禁用')
    await load()
  } catch (_) { /* 请求拦截器已处理错误提示 */ }
}

onMounted(load)
</script>
<style scoped>
.commerce-page { display: flex; flex-direction: column; gap: 16px }
.pager { justify-content: flex-end; margin-top: 16px }
</style>
