package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/c-learn/internal/gamification"
	"github.com/c-learn/pkg/middleware"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "postgres://clearn:clearn@localhost:5432/clearn?sslmode=disable&search_path=gamification"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8004"
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatalf("ping db: %v", err)
	}

	repo := gamification.NewRepository(db)
	svc := gamification.NewService(repo)
	handler := gamification.NewHandler(svc)
	internalHandler := gamification.NewInternalHandler(svc)

	r := gin.Default()
	r.Use(middleware.CORS())

	// Public endpoints
	api := r.Group("/api/v1", middleware.RequireUser())
	{
		api.GET("/profile", handler.GetProfile)
		api.GET("/leaderboard", handler.GetLeaderboard)
	}

	// Internal endpoints
	internal := r.Group("/internal")
	{
		internal.POST("/xp/award", internalHandler.AwardXP)
	}

	log.Printf("Gamification service starting on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("run server: %v", err)
	}
}
