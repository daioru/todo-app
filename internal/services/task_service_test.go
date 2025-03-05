package services_test

import (
	"testing"

	"github.com/daioru/todo-app/internal/helpers"
	"github.com/daioru/todo-app/internal/models"
	"github.com/daioru/todo-app/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var baseErr *helpers.BaseValidationError

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
	t.Parallel()
	t.Run("Successful creation", func(t *testing.T) {
		t.Parallel()
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
		t.Parallel()
		mockRepo := new(MockTaskRepo)
		service := services.NewTaskService(mockRepo)

		task := &models.Task{
			UserID:      1,
			Title:       "",
			Description: "Description",
			Status:      "pending",
		}

		err := service.CreateTask(task)
		assert.ErrorAs(t, err, &baseErr)
		mockRepo.AssertNotCalled(t, "CreateTask")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Title too long", func(t *testing.T) {
		t.Parallel()
		mockRepo := new(MockTaskRepo)
		service := services.NewTaskService(mockRepo)

		task := &models.Task{
			UserID:      1,
			Title:       "sssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssss",
			Description: "Description",
			Status:      "pending",
		}

		err := service.CreateTask(task)
		assert.ErrorAs(t, err, &baseErr)
		mockRepo.AssertNotCalled(t, "CreateTask")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Status empty", func(t *testing.T) {
		t.Parallel()
		mockRepo := new(MockTaskRepo)
		service := services.NewTaskService(mockRepo)

		task := &models.Task{
			UserID:      1,
			Title:       "Title",
			Description: "Description",
			Status:      "",
		}

		err := service.CreateTask(task)
		assert.ErrorAs(t, err, &baseErr)
		mockRepo.AssertNotCalled(t, "CreateTask")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Status too long", func(t *testing.T) {
		t.Parallel()
		mockRepo := new(MockTaskRepo)
		service := services.NewTaskService(mockRepo)

		task := &models.Task{
			UserID:      1,
			Title:       "Title",
			Description: "Description",
			Status:      "sssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssss",
		}

		err := service.CreateTask(task)
		assert.ErrorAs(t, err, &baseErr)
		mockRepo.AssertNotCalled(t, "CreateTask")
		mockRepo.AssertExpectations(t)
	})
}

func TestGetTasksByUser(t *testing.T) {
	t.Parallel()
	mockRepo := new(MockTaskRepo)
	service := services.NewTaskService(mockRepo)

	tasks := []models.Task{
		{ID: 1, Title: "Task 1", UserID: 1},
		{ID: 2, Title: "Task 2", UserID: 1},
	}

	mockRepo.On("GetTasksByUserID", 1).Return(tasks, nil)

	result, err := service.GetTasksByUserID(1)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	mockRepo.AssertExpectations(t)
}

func TestUpdateTask(t *testing.T) {
	t.Parallel()
	t.Run("Successful update", func(t *testing.T) {
		t.Parallel()
		mockRepo := new(MockTaskRepo)
		service := services.NewTaskService(mockRepo)

		updates := map[string]interface{}{
			"id":          1,
			"user_id":     1,
			"title":       "Updated title",
			"description": "Updated description",
			"status":      "completed",
		}

		mockRepo.On("UpdateTask", updates).Return(nil)

		err := service.UpdateTask(updates)
		assert.NoError(t, err)
		mockRepo.AssertNotCalled(t, "UpdateTask")
		mockRepo.AssertExpectations(t)
	})

	t.Run("UserID not specified", func(t *testing.T) {
		t.Parallel()
		mockRepo := new(MockTaskRepo)
		service := services.NewTaskService(mockRepo)

		updates := map[string]interface{}{
			"id":          1,
			"title":       "Updated title",
			"description": "Updated description",
			"status":      "completed",
		}

		err := service.UpdateTask(updates)
		assert.ErrorAs(t, err, &baseErr)
		mockRepo.AssertNotCalled(t, "UpdateTask")
		mockRepo.AssertExpectations(t)
	})

	t.Run("TaskID not specified", func(t *testing.T) {
		t.Parallel()
		mockRepo := new(MockTaskRepo)
		service := services.NewTaskService(mockRepo)

		updates := map[string]interface{}{
			"user_id":     1,
			"title":       "Updated title",
			"description": "Updated description",
			"status":      "completed",
		}

		err := service.UpdateTask(updates)
		assert.ErrorAs(t, err, &baseErr)
		mockRepo.AssertNotCalled(t, "UpdateTask")
		mockRepo.AssertExpectations(t)
	})

	t.Run("No fields to update", func(t *testing.T) {
		t.Parallel()
		mockRepo := new(MockTaskRepo)
		service := services.NewTaskService(mockRepo)

		updates := map[string]interface{}{
			"id":      1,
			"user_id": 1,
		}

		err := service.UpdateTask(updates)
		assert.ErrorAs(t, err, &baseErr)
		mockRepo.AssertNotCalled(t, "UpdateTask")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Unexpected field", func(t *testing.T) {
		t.Parallel()
		mockRepo := new(MockTaskRepo)
		service := services.NewTaskService(mockRepo)

		updates := map[string]interface{}{
			"id":               1,
			"user_id":          1,
			"unexpected_field": "unexpected update",
		}

		err := service.UpdateTask(updates)
		assert.ErrorAs(t, err, &baseErr)
		mockRepo.AssertNotCalled(t, "UpdateTask")
		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteTask(t *testing.T) {
	t.Parallel()
	mockRepo := new(MockTaskRepo)
	service := services.NewTaskService(mockRepo)

	mockRepo.On("DeleteTask", 1, 1).Return(nil)

	err := service.DeleteTask(1, 1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
