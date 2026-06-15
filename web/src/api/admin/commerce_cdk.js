import request from "../../request/index.js";
const baseUrl = '/v1/admin/commerce/cdk'
export const apiAdminCommerceCdkCreate = (data) => request.request({url: `${baseUrl}/create`, method: 'POST', data})
export const apiAdminCommerceCdkPage = (data) => request.request({url: `${baseUrl}/page`, method: 'POST', data})
export const apiAdminCommerceCdkDisable = (data) => request.request({url: `${baseUrl}/disable`, method: 'POST', data})
