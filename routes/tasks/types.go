package tasks

// Save Tasks Handlers custom Structures
type Task struct {
	Title string `json:"title"`
	Description string `json:"description"`
	Points float64 `json:"points"`
	Status string `json:"status"`
}

type ErrorResponse struct {
	Code int `json:"code"`
	Message string `json:"message"`
}

type ErrorsResponse struct {
	Code int `json:"code"`
	Messages []string `json:"messages"`
}