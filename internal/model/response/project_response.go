package response

type ProjectResponse struct {
	ID          uint     `json:"id"`
	Title       string   `json:"title"`
	Images      []string `json:"images"`
	Description string   `json:"description"`
	Tools       []string `json:"tools"`
}
