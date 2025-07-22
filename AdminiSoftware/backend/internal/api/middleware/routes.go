
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// RequireRole middleware to check user roles
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
			c.Abort()
			return
		}

		role := userRole.(string)
		for _, allowedRole := range roles {
			if role == allowedRole {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		c.Abort()
	}
}

// CheckPermission middleware for granular permissions
func CheckPermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		permissions, exists := c.Get("user_permissions")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User permissions not found"})
			c.Abort()
			return
		}

		userPermissions := permissions.([]string)
		for _, userPerm := range userPermissions {
			if userPerm == permission || strings.HasSuffix(userPerm, "*") {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		c.Abort()
	}
}
