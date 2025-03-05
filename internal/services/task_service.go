package services

import (
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

func (s *TaskService) CreateTask(task *models.Task) error {
	err := helpers.ValidateTaskFields(task)
	if err != nil {
		return err
	}

	return s.taskRepo.CreateTask(task)
}

func (s *TaskService) GetTasksByUserID(userID int) ([]models.Task, error) {
	return s.taskRepo.GetTasksByUserID(userID)
}

func (s *TaskService) UpdateTask(updates map[string]interface{}) error {
	updates, err := helpers.ValidateUpdates(updates)
	if err != nil {
		return err
	}

	return s.taskRepo.UpdateTask(updates)
}

func (s *TaskService) DeleteTask(taskID, userID int) error {
	return s.taskRepo.DeleteTask(taskID, userID)
}
