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

	// 应用CORS中间件
	r.Use(middleware.Cors())

	// 应用全局限流中间件
	if config.AppConfig.App.RateLimit.Enabled {
		r.Use(middleware.IPRateLimit(
			config.AppConfig.App.RateLimit.Rate,
			config.AppConfig.App.RateLimit.Capacity,
		))
	}

	// 设置版本化路由
	setupRoutes(r)

	return r
}

// setupRoutes 设置所有版本的路由
func setupRoutes(r *gin.Engine) {
	// API 路由前缀
	prefix := "api"
	if config.AppConfig.App.RouterPrefix != "" {
		prefix = config.AppConfig.App.RouterPrefix
	}

	// v1 版本路由
	setupV1Routes(r, prefix)

	// 未来可以添加 v2 版本
	// setupV2Routes(r, prefix)
}

// setupV1Routes 设置 v1 版本路由
func setupV1Routes(r *gin.Engine, prefix string) {
	v1 := r.Group(prefix + "/v1")

	// 健康检查(公共接口)
	v1.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"version": "v1",
		})
	})

	// 公共路由组
	publicGroup := v1.Group("")
	{
		InitSystemAuthPublicRouter(publicGroup)
	}

	// 需要认证的路由组
	privateGroup := v1.Group("")
	privateGroup.Use(middleware.JWTAuth(), middleware.CasbinAuth())
	{
		InitApiAuthRouter(privateGroup)
		InitDictAuthRouter(privateGroup)
		InitMenuAuthRouter(privateGroup)
		InitRecordAuthRouter(privateGroup)
		InitRoleAuthRouter(privateGroup)
		InitUserAuthRouter(privateGroup)
	}
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
