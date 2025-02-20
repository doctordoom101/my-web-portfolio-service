package main

import (
	"database/sql"
	"log"
	"os"
	"project-portfolio-api/infrastructure"
	"project-portfolio-api/internal/api"
	"project-portfolio-api/internal/handler"
	"project-portfolio-api/internal/repository"
	"project-portfolio-api/internal/service"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	db := infrastructure.NewPostgresConnection()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}(db)

	// Setup repositories, services, and handlers
	projectRepo := repository.NewProjectRepository(db)
	projectService := service.NewProjectService(projectRepo)
	projectHandler := handler.NewProjectHandler(projectService)

	// Setup Gin router
	r := gin.Default()
	r.MaxMultipartMemory = 20 << 20

	// Enable CORS for all origins
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Create uploads directory if it doesn't exist
	if err := os.MkdirAll("uploads", 0755); err != nil {
		log.Fatal("Failed to create uploads directory:", err)
	}

	// Serve static files
	r.Static("/uploads", "./uploads")

	// Setup routes
	api.SetupRoutes(r, projectHandler)

	// Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
