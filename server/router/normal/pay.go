package normal

import (
	"github.com/gin-gonic/gin"
	"server/controller/normal/pay"
	"server/global"
	"server/router/middleware"
)

func InitPay(group *gin.RouterGroup) {
	g := group.Group("pay", middleware.Auth(global.Jwt))
	g.POST("recharge", pay.Recharge)
}
