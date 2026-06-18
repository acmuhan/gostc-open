package pay

import (
	"github.com/gin-gonic/gin"
	"server/model"
	"server/pkg/epay"
	"server/repository/cache"
	service "server/service/public/pay"
)

var svr = service.Service

func Notify(c *gin.Context) {
	params := epay.ParseNotify(c.Request)
	if err := svr.Notify(params); err != nil {
		c.String(200, "fail")
		return
	}
	c.String(200, "success")
}

func Return(c *gin.Context) {
	var cfgBase model.SystemConfigBase
	cache.GetSystemConfigBase(&cfgBase)
	c.Redirect(302, cfgBase.BaseUrl+"/#/normal/wallet")
}
