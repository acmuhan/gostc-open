package wallet

import (
	"github.com/gin-gonic/gin"
	"server/pkg/bean"
	"server/service/admin/wallet"
)

var svr = service.Service

func Adjust(c *gin.Context) {
	var req service.AdjustReq
	if err := c.ShouldBindJSON(&req); err != nil {
		bean.Response.Param(c, err)
		return
	}
	if err := svr.Adjust(req); err != nil {
		bean.Response.Fail(c, err.Error())
		return
	}
	bean.Response.Ok(c)
}
