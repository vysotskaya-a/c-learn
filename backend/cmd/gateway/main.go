package main

import (
	"log"
	"os"

	"github.com/c-learn/internal/gateway"
)

func main() {
	cfg := gateway.Config{
		AuthServiceURL:         getEnv("AUTH_SERVICE_URL", "http://localhost:8001"),
		LMSServiceURL:          getEnv("LMS_SERVICE_URL", "http://localhost:8002"),
		GamificationServiceURL: getEnv("GAMIFICATION_SERVICE_URL", "http://localhost:8004"),
		JWTSecret:              getEnv("JWT_SECRET", "dev-secret-change-in-production"),
	}

	port := getEnv("PORT", "8080")

	r := gateway.NewRouter(cfg)

	log.Printf("API Gateway starting on :%s", port)
	log.Printf("  Auth Service:         %s", cfg.AuthServiceURL)
	log.Printf("  LMS Service:          %s", cfg.LMSServiceURL)
	log.Printf("  Gamification Service: %s", cfg.GamificationServiceURL)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("run server: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
