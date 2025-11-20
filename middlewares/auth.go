package middlewares

import (
	"mygo/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"err": "请求没有携带token，请先登录",
			})
			c.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "令牌格式错误"})
			c.Abort()
			return
		}
		claims := &utils.Claims{}
		token, err := jwt.ParseWithClaims(parts[1], claims, func(token *jwt.Token) (interface{}, error) {
			return utils.JWTSecret, nil
		})
		if err != nil || !token.Valid { /*如果token无效*/
			c.JSON(http.StatusUnauthorized, gin.H{
				"err": "无效的token，请重新登录",
			})
			c.Abort()
			return
		}
		c.Set("userID", claims.UserID)
		c.Next()
	}

}
