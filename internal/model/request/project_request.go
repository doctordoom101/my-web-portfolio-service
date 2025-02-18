package request

type CreateProjectRequest struct {
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description" binding:"required"`
	Tools       []string `json:"tools" binding:"required"`
}

type UpdateProjectRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tools       []string `json:"tools"`
}
