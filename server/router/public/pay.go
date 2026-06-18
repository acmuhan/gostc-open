package public

import (
	"github.com/gin-gonic/gin"
	"server/controller/public/pay"
)

func InitPay(group *gin.RouterGroup) {
	g := group.Group("pay")
	g.GET("notify", pay.Notify)
	g.GET("return", pay.Return)
}
