package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/daioru/todo-app/internal/config"
	"github.com/daioru/todo-app/internal/handlers"
	"github.com/daioru/todo-app/internal/pkg/db"
	"github.com/daioru/todo-app/internal/repository"
	"github.com/daioru/todo-app/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := config.ReadConfigYML("config.yml"); err != nil {
		log.Fatal("Failed init configuration")
	}
	cfg := config.GetConfigInstance()

	db, err := db.ConnectDB(&cfg.DB)
	if err != nil {
		log.Fatalf("sqlx_Open error: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error testing db connection: %v", err)
	}

	userRepo := repository.NewUserRepository(db)

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	jwtSecret := os.Getenv("JWTSECRET")
	if jwtSecret == "" {
		log.Fatal("No jwtSecret in .env")
	}

	authService := services.NewAuthService(userRepo, jwtSecret)

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
	r.POST("/tasks", taskHandler.CreateTask)
	r.PUT("/tasks", taskHandler.UdpateTask)
	r.GET("/tasks", taskHandler.GetTasks)
	r.DELETE("/tasks", taskHandler.DeleteTask)

	// Testing route
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	port := ":8080"
	fmt.Println("Server is running on port" + port)
	if err := r.Run(port); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
