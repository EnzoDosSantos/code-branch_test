package server

import (
	"net/http"

	"github.com/EnzoDosSantos/code-branch_test/internal/handlers"
	"github.com/EnzoDosSantos/code-branch_test/internal/repository"
	"github.com/EnzoDosSantos/code-branch_test/pkg/middleware"
)

func SetupRouter() http.Handler {
	mux := http.NewServeMux()
	
	repo := repository.NewInMemoryTaskRepository()
	taskHandler := handlers.NewTaskHandler(repo)

	mux.HandleFunc("GET /tasks", taskHandler.GetAllTasks)
	mux.HandleFunc("POST /tasks", taskHandler.CreateTask)
	mux.HandleFunc("GET /tasks/{id}", taskHandler.GetTaskByID)
	mux.HandleFunc("PUT /tasks/{id}", taskHandler.UpdateTask)
	mux.HandleFunc("DELETE /tasks/{id}", taskHandler.DeleteTask)

	wHandler := middleware.LoggingMiddleware(middleware.RecoveryMiddleware(mux))

	return wHandler
}