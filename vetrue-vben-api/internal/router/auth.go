package router

import (
	"vetrue-vben-api/internal/handler"

	"github.com/gin-gonic/gin"
)

func InitSystemAuthPublicRouter(routerGroup *gin.RouterGroup) {
	systemAuthGroup := routerGroup.Group("/system/auth")
	{
		systemAuthGroup.POST("/login", handler.SystemAuthLogin)
		systemAuthGroup.POST("/logout", handler.SystemAuthLogout)
	}
}
