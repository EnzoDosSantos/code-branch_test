package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/EnzoDosSantos/code-branch_test/internal/models"
	"github.com/EnzoDosSantos/code-branch_test/internal/repository"
	"github.com/EnzoDosSantos/code-branch_test/pkg/utils"
)

type TaskHandler struct {
    repo repository.TaskRepository
}

func NewTaskHandler(repo repository.TaskRepository) *TaskHandler {
    return &TaskHandler{repo: repo}
}

func (h *TaskHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
    done := make(chan struct{})

    go func() {
        tasks := h.repo.GetAll()

		if len(tasks) == 0 {
			utils.RespondError(w, http.StatusNotFound, "no tasks found")
			close(done)
			return
		} 

        utils.RespondJSON(w, http.StatusOK, tasks)
		
        close(done)
    }()

    <-done
}

func (h *TaskHandler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
    if err != nil {
        utils.RespondError(w, http.StatusBadRequest, "invalid task ID")
        return
    }

    task, err := h.repo.GetByID(id)
    if err != nil {
        utils.RespondError(w, http.StatusNotFound, err.Error())
        return
    }

    utils.RespondJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
    var req models.CreateTaskRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        utils.RespondError(w, http.StatusBadRequest, "invalid request payload")
        return
    }
    
    if req.Title == "" {
        utils.RespondError(w, http.StatusBadRequest, "title is required")
        return
    }

    task := models.Task{
        Title:       req.Title,
        Description: req.Description,
    }

    created := h.repo.Create(task)
    utils.RespondJSON(w, http.StatusCreated, created)
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(r.PathValue("id"))
    if err != nil {
        utils.RespondError(w, http.StatusBadRequest, "invalid task ID")
        return
    }

    var req models.UpdateTaskRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        utils.RespondError(w, http.StatusBadRequest, "invalid request payload")
        return
    }

	if req.Title == nil && req.Description == nil && req.Completed == nil {
		utils.RespondError(w, http.StatusBadRequest, "no fields to update")
		return
	}

    existing, err := h.repo.GetByID(id)
    if err != nil {
        utils.RespondError(w, http.StatusNotFound, err.Error())
        return
    }

    update := *existing
    if req.Title != nil {
        if *req.Title == "" {
            utils.RespondError(w, http.StatusBadRequest, "title cannot be empty")
            return
        }
        update.Title = *req.Title
    }
    if req.Description != nil {
        update.Description = *req.Description
    }
    if req.Completed != nil {
        update.Completed = *req.Completed
    }

    updated, err := h.repo.Update(id, update)
    if err != nil {
        utils.RespondError(w, http.StatusInternalServerError, "failed to update task")
        return
    }

    utils.RespondJSON(w, http.StatusOK, updated)
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(r.PathValue("id"))
    if err != nil {
        utils.RespondError(w, http.StatusBadRequest, "invalid task ID")
        return
    }

    if err := h.repo.Delete(id); err != nil {
        utils.RespondError(w, http.StatusNotFound, err.Error())
        return
    }
    
    w.WriteHeader(http.StatusNoContent)
}