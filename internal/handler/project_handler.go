package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"project-portfolio-api/internal/model/request"
	"project-portfolio-api/internal/repository"
	"project-portfolio-api/internal/service"
	"project-portfolio-api/pkg/custom_error"
	"strconv"
	"strings"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tools format", "details": err.Error()})
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

	// Batas ukuran maksimal file (contoh: 2MB)
	const maxFileSize = 20 << 20 // 2MB

	allowedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}

	// Simpan file dan collect path
	for _, file := range files {
		ext := strings.ToLower(filepath.Ext(file.Filename))
		fmt.Println("File extension:", ext) // Debugging

		if !allowedExtensions[ext] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Only JPG, JPEG, and PNG are allowed"})
			return
		}

		if file.Size > maxFileSize {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("File %s exceeds the maximum size of 2MB", file.Filename)})
			return
		}

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

	// Dapatkan form data
	title := c.PostForm("title")
	description := c.PostForm("description")
	toolsStr := c.PostForm("tools")

	// Parse tools array
	var tools []string
	if toolsStr != "" {
		toolsStr = strings.Trim(toolsStr, "[]")
		tools = strings.Split(toolsStr, ",")
		// Bersihkan spasi dan quotes
		for i, tool := range tools {
			tools[i] = strings.Trim(strings.TrimSpace(tool), "\"")
		}
	}

	// Handle file uploads
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form"})
		return
	}

	var images []string
	if files := form.File["images"]; len(files) > 0 {
		images = make([]string, len(files))
		for i, file := range files {
			// Simpan file ke direktori tertentu
			filename := filepath.Base(file.Filename)
			images[i] = filename

			// Opsional: jika perlu menyimpan file
			dst := "uploads/" + filename
			if err := c.SaveUploadedFile(file, dst); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
				return
			}
		}
	}

	// Buat request object
	req := &request.UpdateProjectRequest{
		Title:       title,
		Description: description,
		Tools:       tools,
		Images:      images,
	}

	// Panggil service
	if err := h.service.Update(uint(id), req); err != nil {
		if appErr, ok := err.(*custom_error.AppError); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Project with updated successfully",
		"data":    req,
	})
}

func (h *ProjectHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Log untuk debugging
	log.Printf("Attempting to delete project with ID: %d", id)

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
