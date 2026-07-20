package middleware

import (
	"blog_project/pkg/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWT鉴权中间件，验证Token，保护需要登录的接口
func JWT(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		appG := utils.Gin{C: c}

		//1、从header获取Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			appG.Error(utils.UNAUTHORIZED, "缺少认证令牌", nil)
			c.Abort()
			return
		}

		//解析Bearer Token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			appG.Error(utils.UNAUTHORIZED, "认证格式无效，请使用Bearer Token", nil)
			c.Abort()
			return
		}
		tokenString := parts[1]

		//验证token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			//确保签名方法正确
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})

		if err != nil || token.Valid {
			appG.Error(utils.UNAUTHORIZED, "无效或已过期令牌", nil)
			c.Abort()
		}

		//提取用户信息存入Context
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			//从claims中提取userID
			if userID, exists := claims["user_id"]; exists {
				c.Set("user_id", userID)
			}
			if username, exists := claims["username"]; exists {
				c.Set("username", username)
			}
		}

		//放行，进入Handler
		c.Next()
	}
}
