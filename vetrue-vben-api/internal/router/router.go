package router

import (
	"net/http"
	"vetrue-vben-api/internal/config"
	"vetrue-vben-api/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	setMode()

	r := gin.Default()

	// 优化静态文件路径配置
	const staticPath = "/static"
	r.StaticFS(staticPath, http.Dir("./"+staticPath))

	r.Use(middleware.Cors())

	// 公共路由组
	publicGroup := r.Group("api/v1")
	{
		// 健康监测
		publicGroup.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, "ok")
		})

		InitSystemAuthPublicRouter(publicGroup)
	}

	// 需要认证的路由组
	privateAuthGroup := r.Group("api/v1")
	privateAuthGroup.Use(middleware.JWTAuth(), middleware.CasbinAuth())
	{
		InitApiAuthRouter(privateAuthGroup)
		InitDictAuthRouter(privateAuthGroup)
		InitMenuAuthRouter(privateAuthGroup)
		InitRecordAuthRouter(privateAuthGroup)
		InitRoleAuthRouter(privateAuthGroup)
		InitUserAuthRouter(privateAuthGroup)
	}

	return r
}

func setMode() {
	env := config.AppConfig.App.Env
	switch env {
	case "debug":
		gin.SetMode(gin.DebugMode)
	case "release":
		gin.SetMode(gin.ReleaseMode)
	default:
		gin.SetMode(gin.TestMode)
	}
}
