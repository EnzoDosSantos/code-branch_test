package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/EnzoDosSantos/code-branch_test/internal/handlers"
	"github.com/EnzoDosSantos/code-branch_test/internal/models"
	"github.com/EnzoDosSantos/code-branch_test/internal/repository"
)

func TestTaskHandlers(t *testing.T) {
    testCases := []struct {
        name          string
        method        string
        path          string
        body          string
        existingTasks []models.Task
        wantStatus    int
        validate      func(t *testing.T, resp *http.Response, repo *repository.InMemoryTaskRepository)
    }{
        {
            name:       "GET tasks list",
            method:     "GET",
            path:       "/tasks",
            existingTasks: []models.Task{{Title: "Existing Task"}},
            wantStatus: http.StatusOK,
            validate: func(t *testing.T, resp *http.Response, repo *repository.InMemoryTaskRepository) {
                var tasks []models.Task
                json.NewDecoder(resp.Body).Decode(&tasks)
                if len(tasks) != 1 {
                    t.Errorf("Expected 1 tasks, got %d", len(tasks))
                }
            },
        },
        {
            name:       "GET empty tasks list",
            method:     "GET",
            path:       "/tasks",
            wantStatus: http.StatusNotFound,
            validate: func(t *testing.T, resp *http.Response, repo *repository.InMemoryTaskRepository) {
                var tasks []models.Task
                json.NewDecoder(resp.Body).Decode(&tasks)
                if len(tasks) != 0 {
                    t.Errorf("Expected 0 tasks, got %d", len(tasks))
                }
            },
        },
        {
            name:       "POST create valid task",
            method:     "POST",
            path:       "/tasks",
            body:       `{"title":"Test Task"}`,
            wantStatus: http.StatusCreated,
            validate: func(t *testing.T, resp *http.Response, repo *repository.InMemoryTaskRepository) {
                var task models.Task
                json.NewDecoder(resp.Body).Decode(&task)
                if task.ID != 1 || task.Title != "Test Task" {
                    t.Errorf("Unexpected task data: %+v", task)
                }
                if len(repo.GetAll()) != 1 {
                    t.Error("Task not stored in repository")
                }
            },
        },
        {
            name:       "POST create task with missing title",
            method:     "POST",
            path:       "/tasks",
            body:       `{"description":"No title"}`,
            wantStatus: http.StatusBadRequest,
            validate: func(t *testing.T, resp *http.Response, repo *repository.InMemoryTaskRepository) {
                var response map[string]string
                json.NewDecoder(resp.Body).Decode(&response)
                if response["error"] != "title is required" {
                    t.Errorf("Unexpected error message: %s", response["error"])
                }
            },
        },
        {
            name:          "GET existing task by ID",
            method:        "GET",
            path:          "/tasks/{{id}}",
            existingTasks: []models.Task{{Title: "Existing Task"}},
            wantStatus:    http.StatusOK,
            validate: func(t *testing.T, resp *http.Response, repo *repository.InMemoryTaskRepository) {
                var task models.Task
                json.NewDecoder(resp.Body).Decode(&task)
                if task.Title != "Existing Task" {
                    t.Errorf("Expected 'Existing Task', got '%s'", task.Title)
                }
            },
        },
        {
            name:       "GET non-existent task",
            method:     "GET",
            path:       "/tasks/999",
            wantStatus: http.StatusNotFound,
            validate: func(t *testing.T, resp *http.Response, repo *repository.InMemoryTaskRepository) {
                var response map[string]string
                json.NewDecoder(resp.Body).Decode(&response)
                if response["error"] != "task not found" {
                    t.Errorf("Expected 'task not found', got '%s'", response["error"])
                }
            },
        },
        {
            name:          "UPDATE existing task",
            method:        "PUT",
            path:          "/tasks/{{id}}",
            body:          `{"title":"Updated Title","completed":true}`,
            existingTasks: []models.Task{{Title: "Original Title"}},
            wantStatus:    http.StatusOK,
            validate: func(t *testing.T, resp *http.Response, repo *repository.InMemoryTaskRepository) {
                var task models.Task
                json.NewDecoder(resp.Body).Decode(&task)
                if task.Title != "Updated Title" || !task.Completed {
                    t.Errorf("Update failed: %+v", task)
                }
            },
        },
        {
            name:          "DELETE existing task",
            method:        "DELETE",
            path:          "/tasks/{{id}}",
            existingTasks: []models.Task{{Title: "To Delete"}},
            wantStatus:    http.StatusNoContent,
            validate: func(t *testing.T, resp *http.Response, repo *repository.InMemoryTaskRepository) {
                if len(repo.GetAll()) != 0 {
                    t.Error("Task not deleted from repository")
                }
            },
        },
        {
            name:       "DELETE non-existent task",
            method:     "DELETE",
            path:       "/tasks/999",
            wantStatus: http.StatusNotFound,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            repo := repository.NewInMemoryTaskRepository()
            for _, task := range tc.existingTasks {
                repo.Create(task)
            }

            handler := handlers.NewTaskHandler(repo)

            path := tc.path
            if strings.Contains(path, "{{id}}") {
                tasks := repo.GetAll()

                if len(tasks) == 0 {
                    t.Fatal("Test setup error: ID placeholder but no existing tasks")
                }

                path = strings.ReplaceAll(path, "{{id}}", strconv.Itoa(tasks[0].ID))
            }

            req := httptest.NewRequest(tc.method, path, bytes.NewBufferString(tc.body))
            req.Header.Set("Content-Type", "application/json")
            w := httptest.NewRecorder()

            mux := http.NewServeMux()
            mux.HandleFunc("GET /tasks", handler.GetAllTasks)
            mux.HandleFunc("POST /tasks", handler.CreateTask)
            mux.HandleFunc("GET /tasks/{id}", handler.GetTaskByID)
            mux.HandleFunc("PUT /tasks/{id}", handler.UpdateTask)
            mux.HandleFunc("DELETE /tasks/{id}", handler.DeleteTask)

            mux.ServeHTTP(w, req)

            if w.Code != tc.wantStatus {
                t.Errorf("Expected status %d, got %d", tc.wantStatus, w.Code)
            }

            if tc.validate != nil {
                resp := &http.Response{
                    StatusCode: w.Code,
                    Header:     w.Header(),
                    Body: w.Result().Body,
                }

				defer resp.Body.Close()

                tc.validate(t, resp, repo)
            }
        })
    }
}

func TestConcurrentTaskCreation(t *testing.T) {
    repo := repository.NewInMemoryTaskRepository()
    handler := handlers.NewTaskHandler(repo)

    const numRequests = 100
    results := make(chan int, numRequests)

    for i := range numRequests {
        go func(i int) {
            payload := `{"title":"Task ` + strconv.Itoa(i) + `"}`
            req := httptest.NewRequest("POST", "/tasks", bytes.NewBufferString(payload))
            w := httptest.NewRecorder()
            handler.CreateTask(w, req)
            results <- w.Code
        }(i)
    }

    successCount := 0
    for range numRequests {
        if <-results == http.StatusCreated {
            successCount++
        }
    }

    if successCount != numRequests {
        t.Errorf("Expected %d successful creates, got %d", numRequests, successCount)
    }

    tasks := repo.GetAll()
    if len(tasks) != numRequests {
        t.Errorf("Expected %d tasks, got %d", numRequests, len(tasks))
    }
}