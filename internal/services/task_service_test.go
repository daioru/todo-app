package services_test

import (
	"testing"

	"github.com/daioru/todo-app/internal/helpers"
	"github.com/daioru/todo-app/internal/models"
	"github.com/daioru/todo-app/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTaskRepo struct {
	mock.Mock
}

func (m *MockTaskRepo) CreateTask(task *models.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockTaskRepo) GetTaskByID(id int) (*models.Task, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Task), args.Error(1)
}

func (m *MockTaskRepo) GetTasksByUserID(userID int) ([]models.Task, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Task), args.Error(1)
}

func (m *MockTaskRepo) DeleteTask(taskID, userID int) error {
	args := m.Called(taskID, userID)
	return args.Error(0)
}

func (m *MockTaskRepo) UpdateTask(updates map[string]interface{}) error {
	args := m.Called(updates)
	return args.Error(0)
}

func TestCreateTask(t *testing.T) {
	t.Run("Successful creation", func(t *testing.T) {
		mockRepo := new(MockTaskRepo)
		service := services.NewTaskService(mockRepo)

		task := &models.Task{
			UserID:      1,
			Title:       "Test Task",
			Description: "Description",
			Status:      "pending",
		}

		mockRepo.On("CreateTask", task).Return(nil)

		err := service.CreateTask(task)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Blank title", func(t *testing.T) {
		mockRepo := new(MockTaskRepo)
		service := services.NewTaskService(mockRepo)

		task := &models.Task{
			UserID:      1,
			Title:       "",
			Description: "Description",
			Status:      "pending",
		}

		err := service.CreateTask(task)
		var baseErr *helpers.BaseValidationError
		assert.ErrorAs(t, err, &baseErr)
	})

	t.Run("Title too long", func(t *testing.T) {
		mockRepo := new(MockTaskRepo)
		service := services.NewTaskService(mockRepo)

		task := &models.Task{
			UserID:      1,
			Title:       "sssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssss",
			Description: "Description",
			Status:      "pending",
		}

		err := service.CreateTask(task)
		var baseErr *helpers.BaseValidationError
		assert.ErrorAs(t, err, &baseErr)
	})

	t.Run("Status empty", func(t *testing.T) {
		mockRepo := new(MockTaskRepo)
		service := services.NewTaskService(mockRepo)

		task := &models.Task{
			UserID:      1,
			Title:       "Title",
			Description: "Description",
			Status:      "",
		}

		err := service.CreateTask(task)
		var baseErr *helpers.BaseValidationError
		assert.ErrorAs(t, err, &baseErr)
	})

	t.Run("Status too long", func(t *testing.T) {
		mockRepo := new(MockTaskRepo)
		service := services.NewTaskService(mockRepo)

		task := &models.Task{
			UserID:      1,
			Title:       "Title",
			Description: "Description",
			Status:      "sssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssss",
		}

		err := service.CreateTask(task)
		var baseErr *helpers.BaseValidationError
		assert.ErrorAs(t, err, &baseErr)
	})
}
