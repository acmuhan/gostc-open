package commerce_cdk

import (
	"github.com/gin-gonic/gin"
	"server/pkg/bean"
	"server/router/middleware"
	"server/service/normal/commerce_cdk"
)

var svr = service.Service

func Redeem(c *gin.Context) {
	var req service.RedeemReq
	if err := c.ShouldBindJSON(&req); err != nil {
		bean.Response.Param(c, err)
		return
	}
	if err := svr.Redeem(middleware.GetClaims(c), req); err != nil {
		bean.Response.Fail(c, err.Error())
		return
	}
	bean.Response.Ok(c)
}
