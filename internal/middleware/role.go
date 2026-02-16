package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleValue, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "role not found in token",
			})
			return
		}

		userRole, ok := roleValue.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "invalid role type",
			})
			return
		}

		for _, role := range roles {
			if userRole == role {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "insufficient permissions",
		})
	}
}
