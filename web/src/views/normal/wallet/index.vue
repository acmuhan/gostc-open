<template>
  <div class="commerce-page">
    <n-grid :cols="1" :x-gap="16" :y-gap="16" responsive="screen">
      <n-grid-item><n-card><div class="metric-label">积分</div><div class="metric-value">{{ money(wallet?.amount) }}</div></n-card></n-grid-item>
    </n-grid>
    <n-card class="panel" title="CDK 兑换">
      <n-space>
        <n-input v-model:value="cdk" placeholder="输入兑换码" style="width: 320px" />
        <n-button type="primary" :loading="redeemLoading" @click="redeem">兑换</n-button>
      </n-space>
    </n-card>
    <n-card class="panel" title="资金流水">
      <n-data-table :columns="ledgerColumns" :data="ledgers" :loading="loading" />
      <n-pagination class="pager" v-model:page="page" :page-count="pageCount" @update:page="loadLedger" />
    </n-card>
    <n-card class="panel" title="订单记录">
      <n-data-table :columns="orderColumns" :data="orders" :loading="orderLoading" />
    </n-card>
  </div>
</template>
<script setup>
import { computed, h, onMounted, ref } from 'vue'
import { NTag } from 'naive-ui'
import { apiWalletSummary, apiWalletLedger } from '../../../api/normal/wallet.js'
import { apiCommerceOrderPage } from '../../../api/normal/commerce_order.js'
import { apiCommerceCdkRedeem } from '../../../api/normal/commerce_cdk.js'

const message = window.$message
const wallet = ref({})
const ledgers = ref([])
const orders = ref([])
const page = ref(1)
const total = ref(0)
const loading = ref(false)
const orderLoading = ref(false)
const redeemLoading = ref(false)
const cdk = ref('')
const pageCount = computed(() => Math.max(1, Math.ceil(total.value / 10)))
const money = v => Number(v || 0).toFixed(2)

const ledgerColumns = [
  { title: '时间', key: 'created_at' },
  { title: '账户', key: 'account_type' },
  { title: '类型', key: 'biz_type' },
  { title: '金额', key: 'amount', render: r => h('span', { class: r.direction === 2 ? 'out' : 'in' }, `${r.direction === 2 ? '-' : '+'}${money(r.amount)}`) },
  { title: '余额变化', key: 'balance_after', render: r => `${money(r.balance_before)} → ${money(r.balance_after)}` },
  { title: '备注', key: 'remark' },
]
const orderColumns = [
  { title: '订单号', key: 'order_no' },
  { title: '业务', key: 'biz_type' },
  { title: '支付', key: 'pay_type' },
  { title: '金额', key: 'amount', render: r => money(r.amount) },
  { title: '状态', key: 'status', render: r => h(NTag, { type: r.status === 2 ? 'success' : 'default' }, () => r.status === 2 ? '已支付' : '待处理') },
  { title: '时间', key: 'created_at' },
]

async function loadWallet() {
  try {
    const res = await apiWalletSummary()
    wallet.value = res.data || {}
  } catch (_) { /* 请求拦截器已处理错误提示 */ }
}

async function loadLedger() {
  loading.value = true
  try {
    const res = await apiWalletLedger({ page: page.value, size: 10 })
    ledgers.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch (_) { /* 请求拦截器已处理错误提示 */ }
  loading.value = false
}

async function loadOrders() {
  orderLoading.value = true
  try {
    const res = await apiCommerceOrderPage({ page: 1, size: 10 })
    orders.value = res.data?.list || []
  } catch (_) { /* 请求拦截器已处理错误提示 */ }
  orderLoading.value = false
}

async function redeem() {
  if (!cdk.value) { message.warning('请输入兑换码'); return }
  redeemLoading.value = true
  try {
    await apiCommerceCdkRedeem({ code: cdk.value })
    message.success('兑换成功')
    cdk.value = ''
    await loadWallet()
    await loadLedger()
  } catch (_) { /* 请求拦截器已处理错误提示 */ }
  redeemLoading.value = false
}

onMounted(() => { loadWallet(); loadLedger(); loadOrders() })
</script>
<style scoped>
.commerce-page { display: flex; flex-direction: column; gap: 16px }
.panel { margin-top: 4px }
.metric-label { color: #777; font-size: 13px }
.metric-value { font-size: 30px; font-weight: 700; margin-top: 8px }
.pager { justify-content: flex-end; margin-top: 16px }
.in { color: #18a058 }
.out { color: #d03050 }
</style>
