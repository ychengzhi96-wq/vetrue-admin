package middleware

import (
	"strings"
	"vetrue-vben-api/pkg/jwt"
	"vetrue-vben-api/pkg/response"

	"github.com/gin-gonic/gin"
)

// JWTAuth JWT认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.FailWithCode(c, 401, "请求未携带token，无权限访问")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.FailWithCode(c, 401, "token格式错误")
			c.Abort()
			return
		}

		claims, err := jwt.ParseToken(parts[1])
		if err != nil {
			response.FailWithCode(c, 401, "token无效")
			c.Abort()
			return
		}

		// 将用户信息存储到上下文
		c.Set("userId", claims.ID)
		c.Set("username", claims.Username)
		c.Set("roleId", claims.RoleID)

		c.Next()
	}
}
