import request from "../../request/index.js";
const baseUrl = '/v1/normal/commerce/cdk'
export const apiCommerceCdkRedeem = (data) => request.request({url: `${baseUrl}/redeem`, method: 'POST', data})
