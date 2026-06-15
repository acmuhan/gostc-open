package normal

import (
	"github.com/gin-gonic/gin"
	"server/controller/normal/commerce_order"
	"server/global"
	"server/router/middleware"
)

func InitCommerceOrder(group *gin.RouterGroup) {
	g := group.Group("commerce/order", middleware.Auth(global.Jwt))
	g.POST("page", commerce_order.Page)
}
