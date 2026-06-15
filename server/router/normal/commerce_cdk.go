package normal

import (
	"github.com/gin-gonic/gin"
	"server/controller/normal/commerce_cdk"
	"server/global"
	"server/router/middleware"
)

func InitCommerceCdk(group *gin.RouterGroup) {
	g := group.Group("commerce/cdk", middleware.Auth(global.Jwt))
	g.POST("redeem", commerce_cdk.Redeem)
}
