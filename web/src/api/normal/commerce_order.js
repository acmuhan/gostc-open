import request from "../../request/index.js";
const baseUrl = '/v1/normal/commerce/order'
export const apiCommerceOrderPage = (data) => request.request({url: `${baseUrl}/page`, method: 'POST', data})
