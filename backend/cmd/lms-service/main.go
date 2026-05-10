package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/c-learn/internal/lms"
	"github.com/c-learn/pkg/middleware"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "postgres://clearn:clearn@localhost:5432/clearn?sslmode=disable&search_path=lms"
	}
	runnerURL := os.Getenv("CODE_RUNNER_URL")
	if runnerURL == "" {
		runnerURL = "http://localhost:8003"
	}
	gamifURL := os.Getenv("GAMIFICATION_URL")
	if gamifURL == "" {
		gamifURL = "http://localhost:8004"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8002"
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatalf("ping db: %v", err)
	}

	repo := lms.NewRepository(db)
	runnerClient := lms.NewRunnerClient(runnerURL)
	gamifClient := lms.NewGamificationClient(gamifURL)
	svc := lms.NewService(repo, runnerClient, gamifClient)

	handler := lms.NewHandler(svc)
	submitHandler := lms.NewSubmitHandler(svc)
	adminHandler := lms.NewAdminHandler(repo)

	r := gin.Default()
	r.Use(middleware.CORS())

	// Public endpoints (require X-User-ID from gateway)
	api := r.Group("/api/v1", middleware.RequireUser())
	{
		api.GET("/courses/tree", handler.GetCourseTree)
		api.GET("/lessons/:id", handler.GetLesson)
		api.POST("/tasks/:id/run", submitHandler.Run)
		api.POST("/tasks/:id/submit", submitHandler.Submit)
		api.GET("/solutions", handler.ListSolutions)
	}

	// Admin CMS endpoints
	admin := r.Group("/api/v1/admin", middleware.RequireUser(), middleware.RequireAdmin())
	{
		admin.GET("/modules", adminHandler.ListModules)
		admin.GET("/modules/:id", adminHandler.GetModuleFull)
		admin.POST("/modules", adminHandler.CreateModule)
		admin.PUT("/modules/:id", adminHandler.UpdateModule)
		admin.DELETE("/modules/:id", adminHandler.DeleteModule)

		admin.POST("/lessons", adminHandler.CreateLesson)
		admin.PUT("/lessons/:id", adminHandler.UpdateLesson)
		admin.DELETE("/lessons/:id", adminHandler.DeleteLesson)

		admin.POST("/tasks", adminHandler.CreateTask)
		admin.PUT("/tasks/:id", adminHandler.UpdateTask)
		admin.DELETE("/tasks/:id", adminHandler.DeleteTask)
		admin.PUT("/tasks/:id/test-cases", adminHandler.UpdateTestCases)
	}

	log.Printf("LMS service starting on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("run server: %v", err)
	}
}
