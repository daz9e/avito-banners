package middleware

import (
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

		if token != "admin_token" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Admin token required"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func getRoleByToken(token string) (string, bool) {
	switch token {
	case "admin_token":
		return "admin", true
	case "user_token":
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
