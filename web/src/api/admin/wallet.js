import request from "../../request/index.js";
const baseUrl = '/v1/admin/wallet'
export const apiAdminWalletAdjust = (data) => request.request({url: `${baseUrl}/adjust`, method: 'POST', data})
