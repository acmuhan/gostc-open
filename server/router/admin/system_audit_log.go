package admin

import (
	"github.com/gin-gonic/gin"
	"server/controller/admin/system_audit_log"
	"server/global"
	"server/router/middleware"
)

func InitSystemAuditLog(group *gin.RouterGroup) {
	g := group.Group("system/audit-log", middleware.Auth(global.Jwt), middleware.AuthAdmin())
	g.POST("page", system_audit_log.Page)
}
