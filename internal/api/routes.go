package api

import (
	"project-portfolio-api/internal/handler"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, projectHandler *handler.ProjectHandler) {
	api := r.Group("/api")
	api.Use(cors.Default())
	{
		projects := api.Group("/projects")
		{
			projects.POST("/", projectHandler.Create)
			projects.GET("/", projectHandler.GetAll)
			projects.GET("/:id", projectHandler.GetByID)
			projects.PUT("/:id", projectHandler.Update)
			projects.DELETE("/:id", projectHandler.Delete)
		}
	}
}
