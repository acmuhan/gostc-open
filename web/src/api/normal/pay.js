import request from "../../request/index.js";
const baseUrl = '/v1/normal/pay'
export const apiNormalPayRecharge = (data) => request.request({url: `${baseUrl}/recharge`, method: 'POST', data})
export const apiNormalPayClose = (data) => request.request({url: `${baseUrl}/close`, method: 'POST', data})
export const apiNormalPayDetail = (data) => request.request({url: `${baseUrl}/detail`, method: 'POST', data})
export const apiNormalPayRepay = (data) => request.request({url: `${baseUrl}/repay`, method: 'POST', data})
