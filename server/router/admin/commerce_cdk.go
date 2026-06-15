package admin

import (
	"github.com/gin-gonic/gin"
	"server/controller/admin/commerce_cdk"
	"server/global"
	"server/router/middleware"
)

func InitCommerceCdk(group *gin.RouterGroup) {
	g := group.Group("commerce/cdk", middleware.Auth(global.Jwt), middleware.AuthAdmin())
	g.POST("create", commerce_cdk.Create)
	g.POST("page", commerce_cdk.Page)
	g.POST("disable", commerce_cdk.Disable)
}
