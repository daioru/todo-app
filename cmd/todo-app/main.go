package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/daioru/todo-app/internal/config"
	"github.com/daioru/todo-app/internal/handlers"
	"github.com/daioru/todo-app/internal/logger"
	"github.com/daioru/todo-app/internal/middlewares"
	"github.com/daioru/todo-app/internal/pkg/db"
	"github.com/daioru/todo-app/internal/repository"
	"github.com/daioru/todo-app/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	logger.InitLogger()
	log := logger.GetLogger()

	if err := config.ReadConfigYML("config.yml"); err != nil {
		log.Fatal().Msg("Failed init configuration")
	}
	cfg := config.GetConfigInstance()

	db, err := db.ConnectDB(&cfg.DB)
	if err != nil {
		log.Fatal().Msgf("sqlx_Open error: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal().Msgf("Error testing db connection: %v", err)
	}

	userRepo := repository.NewUserRepository(db)

	err = godotenv.Load()
	if err != nil {
		log.Fatal().Msg("Error loading .env file")
	}

	jwtSecret := os.Getenv("JWTSECRET")
	if jwtSecret == "" {
		log.Fatal().Msg("No jwtSecret in .env")
	}

	authService := services.NewAuthService(userRepo, []byte(jwtSecret))
	authHandler := handlers.NewAuthHandler(authService)

	taskRepo := repository.NewTaskRepository(db)
	taskService := services.NewTaskService(taskRepo)
	taskHandler := handlers.NewTaskHandler(taskService)

	r := gin.Default()
	r.SetTrustedProxies(nil)

	// Authentication routes
	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)

	// Task routes
	taskRoutes := r.Group("/tasks")
	taskRoutes.Use(middlewares.AuthMiddleware([]byte(jwtSecret))) // Защищаем маршруты
	{
		taskRoutes.POST("", taskHandler.CreateTask)
		taskRoutes.GET("", taskHandler.GetTasks)
		taskRoutes.PUT("", taskHandler.UpdateTask)
		taskRoutes.DELETE("/:id", taskHandler.DeleteTask)
	}

	// Testing route
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	port := ":8080"
	fmt.Println("Server is running on port" + port)
	if err := r.Run(port); err != nil {
		log.Fatal().Msgf("Failed to start server: %v", err)
	}
}
