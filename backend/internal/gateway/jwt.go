package gateway

import (
	"net/http"
	"strings"

	"github.com/c-learn/internal/auth"
	"github.com/gin-gonic/gin"
)

func JWTMiddleware(jwtMgr *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip JWT for public endpoints
		path := c.Request.URL.Path
		if isPublicPath(path) {
			c.Next()
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized", "message": "missing authorization header"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized", "message": "invalid authorization format"})
			c.Abort()
			return
		}

		claims, err := jwtMgr.ValidateAccessToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized", "message": "invalid or expired token"})
			c.Abort()
			return
		}

		// Set headers for downstream services
		c.Request.Header.Set("X-User-ID", claims.UserID)
		c.Request.Header.Set("X-User-Role", claims.Role)
		c.Next()
	}
}

func isPublicPath(path string) bool {
	publicPaths := []string{
		"/api/v1/auth/register",
		"/api/v1/auth/login",
		"/api/v1/auth/refresh",
	}
	for _, p := range publicPaths {
		if path == p {
			return true
		}
	}
	return false
}
