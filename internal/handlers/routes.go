package handlers

import (
	"github.com/daioru/todo-app/internal/middlewares"
	"github.com/gin-gonic/gin"
	files "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handlers struct {
	authHandler *AuthHandler
	taskHandler *TaskHandler
}

func NewHandlers(authHandler *AuthHandler, taskHandler *TaskHandler) *Handlers {
	return &Handlers{
		authHandler: authHandler,
		taskHandler: taskHandler,
	}
}

func (h *Handlers) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
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

	r.GET("/swagger/*any", ginSwagger.WrapHandler(files.Handler))
}
