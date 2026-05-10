package gateway

import (
	"github.com/c-learn/internal/auth"
	"github.com/gin-gonic/gin"
)

type Config struct {
	AuthServiceURL         string
	LMSServiceURL          string
	GamificationServiceURL string
	JWTSecret              string
}

func NewRouter(cfg Config) *gin.Engine {
	r := gin.Default()

	jwtMgr := auth.NewJWTManager(cfg.JWTSecret)

	// Global middleware
	r.Use(CORSMiddleware())
	r.Use(JWTMiddleware(jwtMgr))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Auth Service routes
	authProxy := ProxyHandler(cfg.AuthServiceURL)
	r.Any("/api/v1/auth/*path", authProxy)

	// LMS Service routes
	lmsProxy := ProxyHandler(cfg.LMSServiceURL)
	r.Any("/api/v1/courses/*path", lmsProxy)
	r.Any("/api/v1/lessons/*path", lmsProxy)
	r.Any("/api/v1/tasks/*path", lmsProxy)
	r.Any("/api/v1/solutions", lmsProxy)
	r.Any("/api/v1/solutions/*path", lmsProxy)
	r.Any("/api/v1/admin/*path", lmsProxy)

	// Gamification Service routes
	gamifProxy := ProxyHandler(cfg.GamificationServiceURL)
	r.Any("/api/v1/profile", gamifProxy)
	r.Any("/api/v1/profile/*path", gamifProxy)
	r.Any("/api/v1/leaderboard", gamifProxy)
	r.Any("/api/v1/leaderboard/*path", gamifProxy)

	return r
}
