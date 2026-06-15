package normal

import (
	"github.com/gin-gonic/gin"
	"server/controller/normal/wallet"
	"server/global"
	"server/router/middleware"
)

func InitWallet(group *gin.RouterGroup) {
	g := group.Group("wallet", middleware.Auth(global.Jwt))
	g.POST("summary", wallet.Summary)
	g.POST("ledger", wallet.Ledger)
}
