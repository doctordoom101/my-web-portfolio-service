package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"project-portfolio-api/internal/model/request"
	"project-portfolio-api/internal/repository"
	"project-portfolio-api/internal/service"
	"project-portfolio-api/pkg/custom_error"
	"strconv"
	"time"

	// "strings"

	"github.com/gin-gonic/gin"
)

type ProjectHandler struct {
	service service.ProjectService
	Repo    *repository.ProjectRepository
}

func NewProjectHandler(service service.ProjectService) *ProjectHandler {
	return &ProjectHandler{service: service}
}

func (h *ProjectHandler) Create(c *gin.Context) {
	// Ambil data form biasa
	title := c.PostForm("title")
	description := c.PostForm("description")
	toolsStr := c.PostForm("tools")

	// Parse tools dari string JSON ke slice
	var tools []string
	if err := json.Unmarshal([]byte(toolsStr), &tools); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tools format"})
		return
	}

	// Handle file uploads
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form"})
		return
	}

	files := form.File["images"]
	var imagePaths []string

	// Buat direktori uploads jika belum ada
	if err := os.MkdirAll("uploads", 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	// Simpan file dan collect path
	for _, file := range files {
		// Generate unique filename
		filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(file.Filename))
		filepath := filepath.Join("uploads", filename)

		if err := c.SaveUploadedFile(file, filepath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}
		imagePaths = append(imagePaths, filepath)
	}

	// Buat project request
	req := &request.CreateProjectRequest{
		Title:       title,
		Description: description,
		Tools:       tools,
	}

	// Panggil service dengan request dan image paths
	if err := h.service.Create(req, imagePaths); err != nil {
		if appErr, ok := err.(*custom_error.AppError); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Project created successfully",
	})
}

func (h *ProjectHandler) GetAll(c *gin.Context) {
	projects, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, projects)
}

func (h *ProjectHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	project, err := h.service.GetByID(uint(id))
	if err != nil {
		if appErr, ok := err.(*custom_error.AppError); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, project)
}

func (h *ProjectHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req request.UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Update(uint(id), &req); err != nil {
		if appErr, ok := err.(*custom_error.AppError); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project updated successfully"})
}

func (h *ProjectHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		if appErr, ok := err.(*custom_error.AppError); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project deleted successfully"})
}
