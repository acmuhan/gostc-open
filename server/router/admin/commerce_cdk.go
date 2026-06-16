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
	g.POST("batch-disable", commerce_cdk.BatchDisable)
	g.POST("batch-enable", commerce_cdk.BatchEnable)
	g.POST("batch-delete", commerce_cdk.BatchDelete)
	g.POST("update-remark", commerce_cdk.UpdateRemark)
}
