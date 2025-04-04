package middlewares

import (
	"log"
	"net/http"
	"queue-system-backend/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware verifies JWT token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		// Support for tokens with the "Bearer " prefix
		if strings.HasPrefix(tokenString, "Bearer ") {
			tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		}

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(tokenString)

		if err != nil {
			log.Printf("🔴 Token error: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Set claims in the context
		c.Set("claims", claims)
		c.Set("user_id", claims.UserID)           // User ID from token
		c.Set("role", claims.Role)                // Role from token
		c.Set("company_name", claims.CompanyName) // Company name from token

		log.Printf("🟢 Token successfully verified: UserID=%d, Role=%s, CompanyName=%s",
			claims.UserID, claims.Role, claims.CompanyName)

		c.Next()
	}
}
