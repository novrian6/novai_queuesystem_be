package middlewares

import (
	"log"
	"net/http"
	"strings" // Untuk strings.ToLower

	"queue-system-backend/utils"

	"github.com/gin-gonic/gin"
)

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("claims")

		// Log untuk memudahkan debugging
		log.Printf("🟠 Retrieved claims: %+v, exists: %v", user, exists)

		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
			c.Abort()
			return
		}

		// Validasi claims lebih aman
		claims, ok := user.(*utils.Claims)
		if !ok || !strings.EqualFold(claims.Role, "admin") {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}

		c.Next()
	}
}
