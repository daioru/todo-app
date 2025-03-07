package handlers_test

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/daioru/todo-app/internal/handlers"
	"github.com/daioru/todo-app/internal/helpers"
	"github.com/daioru/todo-app/internal/models"
	"github.com/daioru/todo-app/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTaskService struct {
	mock.Mock
}

func (m *MockTaskService) CreateTask(task *models.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockTaskService) GetTasksByUserID(userID int) ([]models.Task, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Task), args.Error(1)
}

func (m *MockTaskService) UpdateTask(updates map[string]interface{}) error {
	args := m.Called(updates)
	return args.Error(0)
}

func (m *MockTaskService) DeleteTask(taskID, userID int) error {
	args := m.Called(taskID, userID)
	return args.Error(0)
}

func TestCreateTask(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	t.Run("Successful task creation", func(t *testing.T) {
		t.Parallel()
		mockService := new(MockTaskService)
		handler := handlers.NewTaskHandler(mockService)

		task := models.Task{Title: "Test Task", Description: "Test Description", Status: "Pending", UserID: 1}
		mockService.On("CreateTask", &task).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		reqBody := `{"title": "Test Task", "description": "Test Description", "status": "Pending"}`
		c.Request = httptest.NewRequest(http.MethodPost, "/tasks/", bytes.NewBufferString(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("user_id", 1)

		handler.CreateTask(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		t.Parallel()
		mockService := new(MockTaskService)
		handler := handlers.NewTaskHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/tasks/", bytes.NewBufferString(`{"title":}`))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateTask(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid request")
	})

	t.Run("Field validation error", func(t *testing.T) {
		t.Parallel()
		mockService := new(MockTaskService)
		handler := handlers.NewTaskHandler(mockService)

		task := models.Task{
			Title:       "Test Task",
			Description: "Test Description",
			Status:      "sssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssss",
			UserID:      1,
		}
		mockService.On("CreateTask", &task).Return(helpers.NewValidationError("validation error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		reqBody := `{"title": "Test Task", "description": "Test Description", "status": "sssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssss"}`
		c.Request = httptest.NewRequest(http.MethodPost, "/tasks/", bytes.NewBufferString(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("user_id", 1)

		handler.CreateTask(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "validation error")
		mockService.AssertExpectations(t)
	})

	t.Run("Server side error", func(t *testing.T) {
		t.Parallel()
		mockService := new(MockTaskService)
		handler := handlers.NewTaskHandler(mockService)

		task := models.Task{
			Title:       "Test Task",
			Description: "Test Description",
			Status:      "Test status",
			UserID:      1,
		}
		mockService.On("CreateTask", &task).Return(errors.New("DB error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		reqBody := `{"title": "Test Task", "description": "Test Description", "status": "Test status"}`
		c.Request = httptest.NewRequest(http.MethodPost, "/tasks/", bytes.NewBufferString(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("user_id", 1)

		handler.CreateTask(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestGetTasks(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	t.Run("Successful get tasks", func(t *testing.T) {
		t.Parallel()
		mockService := new(MockTaskService)
		handler := handlers.NewTaskHandler(mockService)

		tasks := []models.Task{{Title: "Task 1"}, {Title: "Task 2"}}
		mockService.On("GetTasksByUserID", 1).Return(tasks, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/tasks/", nil)
		c.Set("user_id", 1)

		handler.GetTasks(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Internal server error", func(t *testing.T) {
		t.Parallel()
		mockService := new(MockTaskService)
		handler := handlers.NewTaskHandler(mockService)

		mockService.On("GetTasksByUserID", 1).Return(([]models.Task)(nil), errors.New("DB error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/tasks/", nil)
		c.Set("user_id", 1)

		handler.GetTasks(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestUpdateTask(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	t.Run("Successful update", func(t *testing.T) {
		t.Parallel()
		mockService := new(MockTaskService)
		handler := handlers.NewTaskHandler(mockService)

		updates := map[string]interface{}{"id": 1, "title": "Updated Task", "user_id": 1}
		mockService.On("UpdateTask", updates).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		reqBody := `{"title": "Updated Task"}`
		c.Request = httptest.NewRequest(http.MethodPut, "/tasks/1", bytes.NewBufferString(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Set("user_id", 1)

		handler.UpdateTask(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		t.Parallel()
		mockService := new(MockTaskService)
		handler := handlers.NewTaskHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		reqBody := `{"title":}`
		c.Request = httptest.NewRequest(http.MethodPut, "/tasks/1", bytes.NewBufferString(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateTask(c)

		fmt.Println(w.Body.String())

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid request")
	})

	t.Run("Invalid task ID", func(t *testing.T) {
		t.Parallel()
		mockService := new(MockTaskService)
		handler := handlers.NewTaskHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		reqBody := `{"title": "Updated Task"}`
		c.Request = httptest.NewRequest(http.MethodPut, "/tasks/abc", bytes.NewBufferString(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "id", Value: "abc"}}
		c.Set("user_id", 1)

		handler.UpdateTask(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid task ID")
	})

	t.Run("Invalid task or user ID (DB)", func(t *testing.T) {
		t.Parallel()
		mockService := new(MockTaskService)
		handler := handlers.NewTaskHandler(mockService)

		updates := map[string]interface{}{"id": 1, "title": "Updated Task", "user_id": 1}
		mockService.On("UpdateTask", updates).Return(repository.ErrNoRowsUpdated)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		reqBody := `{"title": "Updated Task"}`
		c.Request = httptest.NewRequest(http.MethodPut, "/tasks/1", bytes.NewBufferString(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Set("user_id", 1)

		handler.UpdateTask(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid task or user ID")
		mockService.AssertExpectations(t)
	})

	t.Run("Updates validation failed", func(t *testing.T) {
		t.Parallel()
		mockService := new(MockTaskService)
		handler := handlers.NewTaskHandler(mockService)

		updates := map[string]interface{}{"id": 1, "title": "Updated Task", "user_id": 1}
		mockService.On("UpdateTask", updates).Return(helpers.NewValidationError("validation error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		reqBody := `{"title": "Updated Task"}`
		c.Request = httptest.NewRequest(http.MethodPut, "/tasks/1", bytes.NewBufferString(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Set("user_id", 1)

		handler.UpdateTask(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "validation error")
		mockService.AssertExpectations(t)
	})

	t.Run("Server side error", func(t *testing.T) {
		t.Parallel()
		mockService := new(MockTaskService)
		handler := handlers.NewTaskHandler(mockService)

		updates := map[string]interface{}{"id": 1, "title": "Updated Task", "user_id": 1}
		mockService.On("UpdateTask", updates).Return(errors.New("DB error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		reqBody := `{"title": "Updated Task"}`
		c.Request = httptest.NewRequest(http.MethodPut, "/tasks/1", bytes.NewBufferString(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Set("user_id", 1)

		handler.UpdateTask(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "server side error")
		mockService.AssertExpectations(t)
	})
}

func TestDeleteTask(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	t.Run("Successful delete", func(t *testing.T) {
		t.Parallel()
		mockService := new(MockTaskService)
		handler := handlers.NewTaskHandler(mockService)

		mockService.On("DeleteTask", 1, 1).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodDelete, "/tasks/1", nil)
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Set("user_id", 1)

		handler.DeleteTask(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Task deleted")
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid task ID format", func(t *testing.T) {
		t.Parallel()
		mockService := new(MockTaskService)
		handler := handlers.NewTaskHandler(mockService)

		mockService.On("DeleteTask", 1, 1).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodDelete, "/tasks/abc", nil)
		c.Params = gin.Params{{Key: "id", Value: "abc"}}

		handler.DeleteTask(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid task ID format")
	})

	t.Run("Invalid/protected task ID", func(t *testing.T) {
		t.Parallel()
		mockService := new(MockTaskService)
		handler := handlers.NewTaskHandler(mockService)

		mockService.On("DeleteTask", 1, 1).Return(repository.ErrNoRowsUpdated)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodDelete, "/tasks/1", nil)
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Set("user_id", 1)

		handler.DeleteTask(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "doesn't exist or access denied")
		mockService.AssertExpectations(t)
	})

	t.Run("Server side error", func(t *testing.T) {
		t.Parallel()
		mockService := new(MockTaskService)
		handler := handlers.NewTaskHandler(mockService)

		mockService.On("DeleteTask", 1, 1).Return(errors.New("DB error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodDelete, "/tasks/1", nil)
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Set("user_id", 1)

		handler.DeleteTask(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "server side error")
		mockService.AssertExpectations(t)
	})
}
