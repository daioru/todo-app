package handlers

import (
	"github.com/daioru/todo-app/internal/middlewares"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	authHandler    *AuthHandler
	taskHandler    *TaskHandler
}

func NewHandlers(authHandler *AuthHandler, taskHandler *TaskHandler) *Handlers {
	return &Handlers{
		authHandler:    authHandler,
		taskHandler:    taskHandler,
	}
}

func (h *Handlers) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", h.authHandler.Register)
			auth.POST("/login", h.authHandler.Login)
		}

		tasks := api.Group("/tasks", middlewares.AuthMiddleware())
		{
			tasks.POST("/", h.taskHandler.CreateTask)
			tasks.GET("/", h.taskHandler.GetTasks)
			tasks.PUT("/:id", h.taskHandler.UpdateTask)
			tasks.DELETE("/:id", h.taskHandler.DeleteTask)
		}
	}
}
