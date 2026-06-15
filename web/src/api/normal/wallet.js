import request from "../../request/index.js";
const baseUrl = '/v1/normal/wallet'
export const apiWalletSummary = () => request.request({url: `${baseUrl}/summary`, method: 'POST'})
export const apiWalletLedger = (data) => request.request({url: `${baseUrl}/ledger`, method: 'POST', data})
