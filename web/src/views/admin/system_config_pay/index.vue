<script setup>
import {onBeforeMount, ref} from "vue";
import {apiAdminSystemConfigPay, apiAdminSystemConfigQuery} from "../../../api/admin/system_config.js";

const $message = window.$message

const state = ref({
  data: {
    enable: '2',
    apiVersion: 'v1',
    apiUrl: '',
    pid: '',
    key: '',
    privateKey: '',
    publicKey: '',
  },
  submitLoading: false,
})

const submit = async () => {
  try {
    state.value.submitLoading = true
    await apiAdminSystemConfigPay(state.value.data)
    $message.success('保存成功')
    await queryFunc()
  } finally {
    state.value.submitLoading = false
  }
}

const queryFunc = async () => {
  try {
    let res = await apiAdminSystemConfigQuery({kind: 'SystemConfigPay'})
    if (res.data) {
      state.value.data = res.data
    }
  } finally {}
}

onBeforeMount(() => {
  queryFunc()
})
</script>

<template>
  <div style="padding: 20px">
    <n-card title="支付配置">
      <n-form>
        <n-form-item label="启用支付">
          <n-switch
            v-model:value="state.data.enable"
            checked-value="1"
            unchecked-value="2"
          >
            <template #checked>启用</template>
            <template #unchecked>禁用</template>
          </n-switch>
        </n-form-item>
        <n-form-item label="API版本">
          <n-radio-group v-model:value="state.data.apiVersion">
            <n-radio value="v1">V1（MD5签名）</n-radio>
            <n-radio value="v2">V2（RSA签名）</n-radio>
          </n-radio-group>
        </n-form-item>
        <n-form-item label="易支付网关地址">
          <n-input v-model:value="state.data.apiUrl" placeholder="如 https://pay.example.com"></n-input>
        </n-form-item>
        <n-form-item label="商户ID (pid)">
          <n-input v-model:value="state.data.pid" placeholder="商户ID"></n-input>
        </n-form-item>
        <n-form-item v-if="state.data.apiVersion==='v1'" label="商户密钥 (Key)">
          <n-input v-model:value="state.data.key" placeholder="V1 MD5密钥" type="password" show-password-on="click"></n-input>
        </n-form-item>
        <n-form-item v-if="state.data.apiVersion==='v2'" label="商户私钥">
          <n-input v-model:value="state.data.privateKey" type="textarea" placeholder="V2 RSA商户私钥" :autosize="{minRows:3,maxRows:8}"></n-input>
        </n-form-item>
        <n-form-item v-if="state.data.apiVersion==='v2'" label="平台公钥">
          <n-input v-model:value="state.data.publicKey" type="textarea" placeholder="V2 RSA平台公钥" :autosize="{minRows:3,maxRows:8}"></n-input>
        </n-form-item>
        <n-form-item>
          <n-button type="primary" :loading="state.submitLoading" @click="submit">保存</n-button>
        </n-form-item>
      </n-form>
    </n-card>
  </div>
</template>

<style scoped lang="scss">
@media (max-width: 520px) {
  :deep(.n-card__content) { padding: 12px !important }
}
</style>
