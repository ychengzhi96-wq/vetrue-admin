package router

import (
	"vetrue-vben-api/internal/config"
	"vetrue-vben-api/pkg/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	setMode()

	r := gin.Default()
	fs := "/static"
	r.StaticFS(fs, http.Dir("./"+fs))

	r.Use(middleware.Cors())
	publicGroup := r.Group("api/v1")
	{
		// 健康监测
		publicGroup.GET("/health", func(c *gin.Context) {
			c.JSON(200, "ok")
		})

		InitSystemAuthPublicRouter(publicGroup)
	}

	privateAuthGroup := r.Group("api/v1")
	privateAuthGroup.Use(middleware.JWTAuth()).Use(middleware.CasbinAuth())
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
