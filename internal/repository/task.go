package repository

import (
	"errors"
	"sync"

	"slices"

	"github.com/EnzoDosSantos/code-branch_test/internal/models"
)

var (
    ErrTaskNotFound = errors.New("task not found")
)

type TaskRepository interface {
    GetAll() []models.Task
    GetByID(id int) (*models.Task, error)
    Create(task models.Task) models.Task
    Update(id int, task models.Task) (*models.Task, error)
    Delete(id int) error
}

type InMemoryTaskRepository struct {
    tasks []models.Task
    mu    sync.Mutex
}

func NewInMemoryTaskRepository() *InMemoryTaskRepository {
    return &InMemoryTaskRepository{
        tasks: make([]models.Task, 0),
    }
}

func (r *InMemoryTaskRepository) GetAll() []models.Task {
    r.mu.Lock()
    defer r.mu.Unlock()
	
    return r.tasks
}

func (r *InMemoryTaskRepository) GetByID(id int) (*models.Task, error) {
    r.mu.Lock()
    defer r.mu.Unlock()

    for _, task := range r.tasks {
        if task.ID == id {
            return &task, nil
        }
    }

    return nil, ErrTaskNotFound
}

func (r *InMemoryTaskRepository) Create(task models.Task) models.Task {
    r.mu.Lock()
    defer r.mu.Unlock()

    task.ID = len(r.tasks) + 1
    r.tasks = append(r.tasks, task)

    return task
}

func (r *InMemoryTaskRepository) Update(id int, update models.Task) (*models.Task, error) {
    r.mu.Lock()
    defer r.mu.Unlock()

    for i, task := range r.tasks {
        if task.ID == id {
            if update.Title != "" {
                r.tasks[i].Title = update.Title
            }

            if update.Description != "" {
                r.tasks[i].Description = update.Description
            }

            r.tasks[i].Completed = update.Completed
            return &r.tasks[i], nil
        }
    }
    return nil, ErrTaskNotFound
}

func (r *InMemoryTaskRepository) Delete(id int) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    for i, task := range r.tasks {
        if task.ID == id {
            r.tasks = slices.Delete(r.tasks, i, i+1)
            return nil
        }
    }

    return ErrTaskNotFound
}