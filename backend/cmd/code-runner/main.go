package main

import (
	"log"
	"os"

	"github.com/c-learn/internal/runner"
	"github.com/c-learn/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8003"
	}
	sandboxImage := os.Getenv("SANDBOX_IMAGE")
	if sandboxImage == "" {
		sandboxImage = "gcc:13"
	}

	log.Printf("Initializing Docker runner (image: %s)...", sandboxImage)

	dockerRunner, err := runner.NewDockerRunner(sandboxImage)
	if err != nil {
		log.Fatalf("init docker runner: %v", err)
	}

	svc := runner.NewService(dockerRunner)
	handler := runner.NewHandler(svc)

	r := gin.Default()
	r.Use(middleware.CORS())

	// Internal-only endpoint
	r.POST("/internal/run", handler.Run)

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "sandbox_image": sandboxImage})
	})

	log.Printf("Code Runner service starting on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("run server: %v", err)
	}
}
