package models

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Completed   bool   `json:"completed"`
}

type CreateTaskRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description,omitempty"`
}

type UpdateTaskRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Completed   *bool   `json:"completed"`
}