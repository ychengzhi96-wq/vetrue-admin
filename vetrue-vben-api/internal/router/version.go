package router

import (
	"vetrue-vben-api/pkg/middleware"

	"github.com/gin-gonic/gin"
)

// VersionGroup 版本路由组配置
type VersionGroup struct {
	Version    string                   // 版本号,如 "v1", "v2"
	Middleware []gin.HandlerFunc        // 该版本专用中间件
	Routers    []func(*gin.RouterGroup) // 路由初始化函数列表
}

// RouterConfig 路由配置
type RouterConfig struct {
	Prefix     string                      // 全局前缀,如 "/api"
	RateLimit  *middleware.RateLimitConfig // 全局限流配置
	Middleware []gin.HandlerFunc           // 全局中间件
	Versions   []VersionGroup              // 版本路由组
}

// SetupVersionRoutes 设置版本化路由
func SetupVersionRoutes(engine *gin.Engine, config RouterConfig) {
	// 应用全局中间件
	if len(config.Middleware) > 0 {
		engine.Use(config.Middleware...)
	}

	// 应用全局限流
	if config.RateLimit != nil {
		engine.Use(middleware.RateLimit(*config.RateLimit))
	}

	// 为每个版本创建路由组
	for _, vg := range config.Versions {
		versionPath := config.Prefix
		if vg.Version != "" {
			versionPath = config.Prefix + "/" + vg.Version
		}

		group := engine.Group(versionPath)

		// 应用版本专用中间件
		if len(vg.Middleware) > 0 {
			group.Use(vg.Middleware...)
		}

		// 注册该版本的所有路由
		for _, routerFunc := range vg.Routers {
			routerFunc(group)
		}
	}
}
