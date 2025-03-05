package services

import (
	"fmt"

	"github.com/daioru/todo-app/internal/helpers"
	"github.com/daioru/todo-app/internal/models"
)

type ITaskRepository interface {
	CreateTask(task *models.Task) error
	GetTaskByID(id int) (*models.Task, error)
	GetTasksByUserID(userID int) ([]models.Task, error)
	DeleteTask(taskID, userID int) error
	UpdateTask(updates map[string]interface{}) error
}

type TaskService struct {
	taskRepo ITaskRepository
}

func NewTaskService(taskRepo ITaskRepository) *TaskService {
	return &TaskService{taskRepo: taskRepo}
}

func ValidateTaskFields(task *models.Task) error {
	if task.Title == "" {
		return fmt.Errorf("validation failed: %w", NewSpecificValidationError("title", "cannot be blank"))
	}

	if len(task.Title) > 100 {
		return fmt.Errorf("validation failed: %w", NewSpecificValidationError("title", "field too long"))
	}

	if task.Status == "" {
		return fmt.Errorf("validation failed: %w", NewSpecificValidationError("status", "cannot be blank"))
	}

	if len(task.Status) > 100 {
		return fmt.Errorf("validation failed: %w", NewSpecificValidationError("title", "status too long"))
	}

	return nil
}

func (s *TaskService) CreateTask(task *models.Task) error {
	err := ValidateTaskFields(task)
	if err != nil {
		return err
	}

	return s.taskRepo.CreateTask(task)
}

func (s *TaskService) GetTasksByUser(userID int) ([]models.Task, error) {
	return s.taskRepo.GetTasksByUserID(userID)
}

func (s *TaskService) UpdateTask(updates map[string]interface{}) error {
	updates, err := helpers.Validate(updates)
	if err != nil {
		return err
	}

	return s.taskRepo.UpdateTask(updates)
}

func (s *TaskService) DeleteTask(taskID, userID int) error {
	return s.taskRepo.DeleteTask(taskID, userID)
}
