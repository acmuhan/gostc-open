<template>
  <div class="commerce-page">
    <n-grid :cols="1" :x-gap="16" :y-gap="16" responsive="screen">
      <n-grid-item><n-card><div class="metric-label">积分</div><div class="metric-value">{{ money(wallet?.amount) }}</div></n-card></n-grid-item>
    </n-grid>
    <n-card class="panel" title="CDK 兑换">
      <n-space>
        <n-input v-model:value="cdk" placeholder="输入兑换码" class="cdk-input" />
        <n-button type="primary" :loading="redeemLoading" @click="redeem">兑换</n-button>
      </n-space>
    </n-card>
    <n-card class="panel" title="积分充值">
      <n-space>
        <n-input-number v-model:value="rechargeAmount" :min="0.01" :precision="2" placeholder="充值金额" style="width:160px" />
        <n-select v-model:value="rechargePayType" :options="payTypeOptions" placeholder="支付方式" style="width:130px" />
        <n-button type="primary" :loading="rechargeLoading" @click="recharge">充值</n-button>
      </n-space>
    </n-card>
    <n-card class="panel" title="积分流水">
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
import { apiNormalPayRecharge } from '../../../api/normal/pay.js'

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
const rechargeAmount = ref(10)
const rechargePayType = ref('alipay')
const rechargeLoading = ref(false)
const payTypeOptions = [
  { label: '支付宝', value: 'alipay' },
  { label: '微信支付', value: 'wxpay' },
  { label: 'QQ钱包', value: 'qqpay' },
]
const pageCount = computed(() => Math.max(1, Math.ceil(total.value / 10)))
const money = v => Number(v || 0).toFixed(2)

const ledgerColumns = [
  { title: '时间', key: 'created_at' },
  { title: '类型', key: 'biz_type', render: r => ({ recharge: '充值', consume: '消费', refund: '退款', admin_adjust: '管理员调整', cdk_redeem: 'CDK兑换', order_pay: '订单支付', checkin: '签到' }[r.biz_type] || r.biz_type) },
  { title: '积分', key: 'amount', render: r => h('span', { class: r.direction === 2 ? 'out' : 'in' }, `${r.direction === 2 ? '-' : '+'}${money(r.amount)}`) },
  { title: '变化', key: 'balance_after', render: r => `${money(r.balance_before)} → ${money(r.balance_after)}` },
  { title: '备注', key: 'remark' },
]
const orderColumns = [
  { title: '订单号', key: 'order_no' },
  { title: '业务', key: 'biz_type', render: r => ({ recharge: '充值', tunnel_create: '隧道开通', tunnel_renew: '隧道续费', host_create: '域名开通', host_renew: '域名续费', forward_create: '转发开通', forward_renew: '转发续费', proxy_create: '代理开通', proxy_renew: '代理续费', p2p_create: 'P2P开通', p2p_renew: 'P2P续费', cdk_redeem: 'CDK兑换', admin_adjust: '管理员调整', auto_renew_tunnel: '自动续费', auto_renew_host: '自动续费', auto_renew_forward: '自动续费', auto_renew_proxy: '自动续费', auto_renew_p2p: '自动续费' }[r.biz_type] || r.biz_type) },
  { title: '支付', key: 'pay_type', render: r => ({ amount: '积分', free: '免费', admin: '管理员', alipay: '支付宝', wxpay: '微信', qqpay: 'QQ钱包' }[r.pay_type] || r.pay_type) },
  { title: '积分', key: 'amount', render: r => money(r.amount) },
  { title: '状态', key: 'status', render: r => h(NTag, { type: { 1: 'warning', 2: 'success', 3: 'default', 4: 'error' }[r.status] || 'default' }, () => ({ 1: '待支付', 2: '已支付', 3: '已关闭', 4: '已退款' }[r.status] || '未知')) },
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

async function recharge() {
  if (!rechargeAmount.value || rechargeAmount.value <= 0) { message.warning('请输入充值金额'); return }
  rechargeLoading.value = true
  try {
    const res = await apiNormalPayRecharge({ amount: String(rechargeAmount.value), payType: rechargePayType.value })
    if (res.data?.payUrl) {
      window.open(res.data.payUrl, '_blank')
    }
  } catch (_) {}
  rechargeLoading.value = false
}

onMounted(() => { loadWallet(); loadLedger(); loadOrders() })
</script>
<style scoped>
.commerce-page { display: flex; flex-direction: column; gap: 16px }
.panel { margin-top: 4px }
.metric-label { color: #777; font-size: 13px }
.metric-value { font-size: 30px; font-weight: 700; margin-top: 8px }
.pager { justify-content: flex-end; margin-top: 16px }
.cdk-input { width: 320px; max-width: 100% }
.in { color: #18a058 }
.out { color: #d03050 }
@media (max-width: 520px) {
  .commerce-page { gap: 8px }
  .metric-value { font-size: 22px }
  :deep(.n-data-table) { overflow-x: auto }
  :deep(.n-card__content) { padding: 12px !important }
  :deep(.n-input) { width: 100% !important }
}
</style>
