package services

import (
	"errors"

	"github.com/daioru/todo-app/internal/models"
	"github.com/daioru/todo-app/internal/repository"
)

type TaskService struct {
	taskRepo *repository.TaskRepository
}

func NewTaskService(taskRepo *repository.TaskRepository) *TaskService {
	return &TaskService{taskRepo: taskRepo}
}

func (s *TaskService) CreateTask(task *models.Task) error {
	if task.Title == "" {
		return errors.New("task title cannot be blank")
	}

	if len(task.Title) > 100 {
		return errors.New("title field too long")
	}

	if task.Status == "" {
		return errors.New("task status cannot be blank")
	}

	if len(task.Status) > 100 {
		return errors.New("status field too long")
	}

	return s.taskRepo.CreateTask(task)
}

func (s *TaskService) GetTasksByUser(userID int) ([]models.Task, error) {
	return s.taskRepo.GetTasksByUserID(userID)
}

func (s *TaskService) UpdateTask(updates map[string]interface{}) error {
	return s.taskRepo.UpdateTask(updates)
}

func (s *TaskService) DeleteTask(taskID, userID int) error {
	return s.taskRepo.DeleteTask(taskID, userID)
}
