package wallet

import (
	"github.com/gin-gonic/gin"
	"server/pkg/bean"
	"server/router/middleware"
	"server/service/normal/wallet"
)

var svr = service.Service

func Summary(c *gin.Context) { bean.Response.OkData(c, svr.Summary(middleware.GetClaims(c))) }
func Ledger(c *gin.Context) {
	var req service.LedgerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		bean.Response.Param(c, err)
		return
	}
	list, total := svr.Ledger(middleware.GetClaims(c), req)
	bean.Response.OkData(c, bean.NewPage(list, total))
}
