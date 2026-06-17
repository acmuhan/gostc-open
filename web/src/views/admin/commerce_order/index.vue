<template>
  <div class="commerce-page">
    <n-card title="订单管理">
      <template #header-extra>
        <n-space>
          <n-select v-model:value="filter.status" :options="statusOptions" placeholder="状态筛选" clearable style="width:130px" @update:value="resetAndLoad" />
          <n-select v-model:value="filter.bizType" :options="bizTypeOptions" placeholder="业务类型" clearable style="width:130px" @update:value="resetAndLoad" />
          <n-input v-model:value="filter.orderNo" placeholder="订单号" clearable style="width:180px" @clear="resetAndLoad" @keyup.enter="resetAndLoad" />
        </n-space>
      </template>
      <n-data-table :columns="columns" :data="list" :loading="loading" :row-key="r => r.order_no" />
      <n-pagination class="pager" v-model:page="page" :page-count="pageCount" @update:page="load" />
    </n-card>
  </div>
</template>
<script setup>
import { computed, h, onMounted, ref } from 'vue'
import { NButton, NTag } from 'naive-ui'
import { apiAdminCommerceOrderPage, apiAdminCommerceOrderRefund } from '../../../api/admin/commerce_order.js'

const message = window.$message
const dialog = window.$dialog

const statusOptions = [
  { label: '待支付', value: 1 },
  { label: '已支付', value: 2 },
  { label: '已关闭', value: 3 },
  { label: '已退款', value: 4 },
]
const bizTypeOptions = [
  { label: '隧道创建', value: 'tunnel_create' },
  { label: '隧道续费', value: 'tunnel_renew' },
  { label: '域名创建', value: 'host_create' },
  { label: '域名续费', value: 'host_renew' },
  { label: '转发创建', value: 'forward_create' },
  { label: '转发续费', value: 'forward_renew' },
  { label: '代理创建', value: 'proxy_create' },
  { label: '代理续费', value: 'proxy_renew' },
  { label: 'P2P创建', value: 'p2p_create' },
  { label: 'P2P续费', value: 'p2p_renew' },
]
const statusMap = { 1: '待支付', 2: '已支付', 3: '已关闭', 4: '已退款' }
const statusType = { 1: 'warning', 2: 'success', 3: 'default', 4: 'error' }

const filter = ref({ status: null, bizType: null, orderNo: '' })
const list = ref([])
const page = ref(1)
const size = 15
const total = ref(0)
const loading = ref(false)
const pageCount = computed(() => Math.ceil(total.value / size) || 1)

const columns = [
  { title: '订单号', key: 'order_no', width: 220 },
  { title: '用户', key: 'user_code', width: 120, ellipsis: true },
  { title: '业务类型', key: 'biz_type', width: 100 },
  { title: '金额', key: 'amount', width: 80 },
  { title: '支付方式', key: 'pay_type', width: 80 },
  {
    title: '状态', key: 'status', width: 80,
    render: (row) => h(NTag, { type: statusType[row.status] || 'default', size: 'small' }, () => statusMap[row.status] || row.status)
  },
  { title: '备注', key: 'remark', width: 120, ellipsis: true },
  {
    title: '支付时间', key: 'paid_at', width: 170,
    render: (row) => row.paid_at ? new Date(row.paid_at * 1000).toLocaleString() : '-'
  },
  {
    title: '创建时间', key: 'created_at', width: 170,
    render: (row) => row.created_at ? new Date(row.created_at).toLocaleString() : '-'
  },
  {
    title: '操作', key: 'action', width: 100,
    render: (row) => {
      if (row.status !== 2 || row.pay_type === 'free') return null
      return h(NButton, { size: 'small', type: 'error', onClick: () => refund(row) }, () => '退款')
    }
  },
]

async function load() {
  loading.value = true
  try {
    const res = await apiAdminCommerceOrderPage({
      page: page.value, size,
      status: filter.value.status || 0,
      bizType: filter.value.bizType || '',
      orderNo: filter.value.orderNo || '',
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

function refund(row) {
  dialog.warning({
    title: '确认退款',
    content: `确认退款订单 ${row.order_no}？退款金额 ${row.amount} 积分将返还给用户。`,
    positiveText: '确认退款',
    negativeText: '取消',
    onPositiveClick: async () => {
      const res = await apiAdminCommerceOrderRefund({ orderNo: row.order_no })
      if (res.code === 0) {
        message.success('退款成功')
        load()
      } else {
        message.error(res.msg || '退款失败')
      }
    }
  })
}

onMounted(load)
</script>
<style scoped>
.pager { margin-top: 16px; display: flex; justify-content: flex-end; }
</style>
