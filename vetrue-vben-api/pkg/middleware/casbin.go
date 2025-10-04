package middleware

import (
	"fmt"
	"vetrue-vben-api/internal/config"
	"vetrue-vben-api/internal/database"
	"vetrue-vben-api/pkg/response"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"
)

var enforcer *casbin.Enforcer

// InitCasbin 初始化Casbin
func InitCasbin() error {
	adapter, err := gormadapter.NewAdapterByDB(database.GetDB())
	if err != nil {
		return fmt.Errorf("创建casbin adapter失败: %w", err)
	}

	modelPath := config.AppConfig.Casbin.ModePath
	if modelPath == "" {
		modelPath = "./configs/rbac_model.conf"
	}

	enforcer, err = casbin.NewEnforcer(modelPath, adapter)
	if err != nil {
		return fmt.Errorf("创建casbin enforcer失败: %w", err)
	}

	// 加载策略
	if err := enforcer.LoadPolicy(); err != nil {
		return fmt.Errorf("加载casbin策略失败: %w", err)
	}

	return nil
}

// CasbinAuth Casbin权限中间件
func CasbinAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if enforcer == nil {
			response.FailWithCode(c, 500, "权限系统未初始化")
			c.Abort()
			return
		}

		// 获取用户角色ID
		roleId, exists := c.Get("roleId")
		if !exists {
			response.FailWithCode(c, 403, "无法获取用户角色")
			c.Abort()
			return
		}

		// 超级管理员（roleId = 1）拥有所有权限
		if roleId == int64(1) {
			c.Next()
			return
		}

		// 获取请求路径和方法
		obj := c.Request.URL.Path
		act := c.Request.Method

		// 检查权限
		sub := fmt.Sprintf("%v", roleId)
		ok, err := enforcer.Enforce(sub, obj, act)
		if err != nil {
			response.FailWithCode(c, 500, "权限检查失败")
			c.Abort()
			return
		}

		if !ok {
			response.FailWithCode(c, 403, "无权限访问")
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetEnforcer 获取Enforcer实例
func GetEnforcer() *casbin.Enforcer {
	return enforcer
}
