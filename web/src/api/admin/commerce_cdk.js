import request from "../../request/index.js";
const baseUrl = '/v1/admin/commerce/cdk'
export const apiAdminCommerceCdkCreate = (data) => request.request({url: `${baseUrl}/create`, method: 'POST', data})
export const apiAdminCommerceCdkPage = (data) => request.request({url: `${baseUrl}/page`, method: 'POST', data})
export const apiAdminCommerceCdkDisable = (data) => request.request({url: `${baseUrl}/disable`, method: 'POST', data})
export const apiAdminCommerceCdkBatchDisable = (data) => request.request({url: `${baseUrl}/batch-disable`, method: 'POST', data})
export const apiAdminCommerceCdkBatchEnable = (data) => request.request({url: `${baseUrl}/batch-enable`, method: 'POST', data})
export const apiAdminCommerceCdkBatchDelete = (data) => request.request({url: `${baseUrl}/batch-delete`, method: 'POST', data})
export const apiAdminCommerceCdkUpdateRemark = (data) => request.request({url: `${baseUrl}/update-remark`, method: 'POST', data})
