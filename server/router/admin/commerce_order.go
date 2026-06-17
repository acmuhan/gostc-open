package admin

import (
	"github.com/gin-gonic/gin"
	"server/controller/admin/commerce_order"
	"server/global"
	"server/router/middleware"
)

func InitCommerceOrder(group *gin.RouterGroup) {
	g := group.Group("commerce/order", middleware.Auth(global.Jwt), middleware.AuthAdmin())
	g.POST("page", commerce_order.Page)
	g.POST("refund", commerce_order.Refund)
}
