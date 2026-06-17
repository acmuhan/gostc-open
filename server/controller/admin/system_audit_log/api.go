package system_audit_log

import (
	"github.com/gin-gonic/gin"
	"server/pkg/bean"
	"server/service/admin/system_audit_log"
)

var svr = service.Service

func Page(c *gin.Context) {
	var req service.PageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		bean.Response.Param(c, err)
		return
	}
	list, total := svr.Page(req)
	bean.Response.OkData(c, bean.NewPage(list, total))
}
