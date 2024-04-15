package middleware

import (
	"avito-banners/internal/config"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, valid := checkToken(c)
		if !valid {
			return
		}

		role, valid := getRoleByToken(token)
		if !valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		c.Set("role", role)
		c.Next()
	}
}

func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, valid := checkToken(c)
		if !valid {
			return
		}

		if token != config.GetAdminToken() {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Admin token required"})
			c.Abort()
			return
		}
		c.Next()
	}
}
func getRoleByToken(token string) (string, bool) {
	switch token {
	case config.GetAdminToken():
		return "admin", true
	case config.GetUserToken():
		return "user", true
	default:
		return "", false
	}
}

func checkToken(c *gin.Context) (string, bool) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
		c.Abort()
		return "", false
	}
	return token, true
}
