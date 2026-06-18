package pay

import (
	"github.com/gin-gonic/gin"
	"server/pkg/bean"
	"server/router/middleware"
	service "server/service/normal/pay"
)

var svr = service.Service

func Recharge(c *gin.Context) {
	var req service.RechargeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		bean.Response.Param(c, err)
		return
	}
	resp, err := svr.Recharge(middleware.GetClaims(c), req, c.ClientIP())
	if err != nil {
		bean.Response.Fail(c, err.Error())
		return
	}
	bean.Response.OkData(c, resp)
}
