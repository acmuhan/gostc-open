import request from "../../request/index.js";
const baseUrl = '/v1/admin/system/audit-log'
export const apiAdminSystemAuditLogPage = (data) => request.request({url: `${baseUrl}/page`, method: 'POST', data})
