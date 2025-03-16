package storage

import (
	"errors"
	"sync"

	"github.com/EnzoDosSantos/code-branch_test/internal/models"
)

type TaskStorage struct {
	tasks  map[int]*models.Task
	mu     sync.RWMutex
}

func NewTaskStorage() *TaskStorage {
	return &TaskStorage{
		tasks:  make(map[int]*models.Task),
	}
}

func (s *TaskStorage) AddTask(task *models.Task) {
	s.mu.Lock()
	defer s.mu.Unlock()

	taskId := len(s.tasks) + 1

	task.ID = taskId
	s.tasks[taskId] = task
}

func (s *TaskStorage) GetLastTask() (*models.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.tasks) == 0 {
		return nil, errors.New("no tasks found")
	}

	return s.tasks[len(s.tasks)], nil
}

func (s *TaskStorage) GetTask(id int) (*models.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, exists := s.tasks[id]

	if !exists {
		return nil, errors.New("task not found")
	}

	return task, nil
}

func (s *TaskStorage) GetAllTasks() []*models.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]*models.Task, 0, len(s.tasks))

	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}

	return tasks
}

func (s *TaskStorage) UpdateTask(id int, updatedTask *models.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[id]; !exists {
		return errors.New("task not found")
	}

	s.tasks[id] = updatedTask
	
	return nil
}

func (s *TaskStorage) DeleteTask(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[id]; !exists {
		return errors.New("task not found")
	}

	delete(s.tasks, id)

	return nil
}