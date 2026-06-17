import request from "../../request/index.js";
const baseUrl = '/v1/admin/commerce/order'
export const apiAdminCommerceOrderPage = (data) => request.request({url: `${baseUrl}/page`, method: 'POST', data})
export const apiAdminCommerceOrderRefund = (data) => request.request({url: `${baseUrl}/refund`, method: 'POST', data})
