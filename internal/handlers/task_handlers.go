package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/daioru/todo-app/internal/helpers"
	"github.com/daioru/todo-app/internal/models"
	"github.com/daioru/todo-app/internal/repository"
	"github.com/daioru/todo-app/internal/services"
	"github.com/gin-gonic/gin"
)

var baseErr *helpers.BaseValidationError

type TaskHandler struct {
	service *services.TaskService
}

func NewTaskHandler(taskService *services.TaskService) *TaskHandler {
	return &TaskHandler{service: taskService}
}

// @Summary CreateTask
// @Description create new task
// @Security Auth
// @Accept  json
// @Produce  json
// @Tags tasks
// @Param input body CreateTaskData true "user info"
// @Success 201 {object} models.Task
// @Failure 400 {object} ErrorResponse
// @Failure 401
// @Failure 500 {object} ErrorResponse
// @Router /tasks/ [post]
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	task.UserID = c.GetInt("user_id")
	if err := h.service.CreateTask(&task); err != nil {
		if errors.As(err, &baseErr) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "server side error"})
		return
	}

	c.JSON(http.StatusCreated, task)
}

// @Summary GetTasks
// @Description get all user tasks
// @Security Auth
// @Accept  json
// @Produce  json
// @Tags tasks
// @Success 200 {object} []models.Task
// @Failure 401
// @Failure 500 {object} ErrorResponse
// @Router /tasks/ [get]
func (h *TaskHandler) GetTasks(c *gin.Context) {
	userID := c.GetInt("user_id")
	tasks, err := h.service.GetTasksByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server side error"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// @Summary UpdateTask
// @Description update stated field in task with {id}
// @Security Auth
// @Accept  json
// @Produce  json
// @Tags tasks
// @Param id path int true "Task ID"
// @Param input body UpdateTaskData true "user info"
// @Success 200
// @Failure 400 {object} ErrorResponse
// @Failure 401
// @Failure 500 {object} ErrorResponse
// @Router /tasks/{id} [put]
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	updates["id"] = taskID

	updates["user_id"] = c.GetInt("user_id")

	if err := h.service.UpdateTask(updates); err != nil {
		if err == repository.ErrNoRowsUpdated {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task or user ID"})
			return
		}

		if errors.As(err, &baseErr) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "server side error"})
		return
	}

	c.Status(http.StatusOK)
}

// @Summary DeleteTask
// @Description delete task with {id}
// @Security Auth
// @Accept  json
// @Produce  json
// @Tags tasks
// @Param id path int true "Task ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401
// @Failure 500 {object} ErrorResponse
// @Router /tasks/{id} [delete]
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	userID := c.GetInt("user_id")
	if err := h.service.DeleteTask(taskID, userID); err != nil {
		if err == repository.ErrNoRowsUpdated {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("task with id: %d doesn't exist or access denied", taskID),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server side error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}
