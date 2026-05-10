package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/c-learn/internal/auth"
	"github.com/c-learn/pkg/middleware"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "postgres://clearn:clearn@localhost:5432/clearn?sslmode=disable&search_path=auth"
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev-secret-change-in-production"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8001"
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatalf("ping db: %v", err)
	}

	repo := auth.NewRepository(db)
	jwtMgr := auth.NewJWTManager(jwtSecret)
	svc := auth.NewService(repo, jwtMgr)
	handler := auth.NewHandler(svc)
	internalHandler := auth.NewInternalHandler(svc)

	r := gin.Default()
	r.Use(middleware.CORS())

	// Public endpoints
	api := r.Group("/api/v1/auth")
	{
		api.POST("/register", handler.Register)
		api.POST("/login", handler.Login)
		api.POST("/refresh", handler.Refresh)
		api.GET("/me", handler.Me)
	}

	// Internal endpoints (called by other services)
	internal := r.Group("/internal")
	{
		internal.GET("/users/:id", internalHandler.GetUser)
	}

	log.Printf("Auth service starting on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("run server: %v", err)
	}
}
