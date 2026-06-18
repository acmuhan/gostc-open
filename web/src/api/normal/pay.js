import request from "../../request/index.js";
const baseUrl = '/v1/normal/pay'
export const apiNormalPayRecharge = (data) => request.request({url: `${baseUrl}/recharge`, method: 'POST', data})
