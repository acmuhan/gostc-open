package admin

import (
	"github.com/gin-gonic/gin"
	"server/controller/admin/wallet"
	"server/global"
	"server/router/middleware"
)

func InitWallet(group *gin.RouterGroup) {
	g := group.Group("wallet", middleware.Auth(global.Jwt), middleware.AuthAdmin())
	g.POST("adjust", wallet.Adjust)
}
