package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Middleware for validating and extracting JWT tokens
// If validation is successful, it sets "user_id" and "email" in the context
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string
		// Check for token in Authorization header
		authToken := c.GetHeader("Authorization")
		if authToken != "" && strings.HasPrefix(authToken, "Bearer ") {
			token = strings.TrimPrefix(authToken, "Bearer ")
		} else {
			// Fallback to HTTP-only cookie
			cookie, err := c.Cookie("token")
			if err == nil {
				token = cookie
			}
		}

		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
		}

		// Validate token and get claims
		claims, err := VerifyIDToken(c, token)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Set user_id and email in context
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)

		c.Next()
	}
}
