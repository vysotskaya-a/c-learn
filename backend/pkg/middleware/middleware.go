package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/c-learn/pkg/errs"
	"github.com/c-learn/pkg/response"
	"github.com/gin-gonic/gin"
)

const (
	HeaderUserID   = "X-User-ID"
	HeaderUserRole = "X-User-Role"
)

// RequireUser extracts X-User-ID from headers set by gateway.
func RequireUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader(HeaderUserID)
		if userID == "" {
			response.Error(c, errs.NewUnauthorized("missing user identity"))
			c.Abort()
			return
		}
		c.Set("user_id", userID)
		c.Set("user_role", c.GetHeader(HeaderUserRole))
		c.Next()
	}
}

// RequireAdmin checks that X-User-Role == "admin".
func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetHeader(HeaderUserRole)
		if role != "admin" {
			response.Error(c, errs.NewForbidden("Admin role required"))
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequestID adds a request ID header for tracing.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader("X-Request-ID")
		if id == "" {
			id = fmt.Sprintf("%d", time.Now().UnixNano())
		}
		c.Set("request_id", id)
		c.Header("X-Request-ID", id)
		c.Next()
	}
}

// CORS middleware.
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type,Authorization,X-Request-ID")
		c.Header("Access-Control-Expose-Headers", "X-Request-ID")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
