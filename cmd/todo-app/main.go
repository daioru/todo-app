package main

import (
	"fmt"
	"os"

	"github.com/daioru/todo-app/internal/config"
	"github.com/daioru/todo-app/internal/handlers"
	"github.com/daioru/todo-app/internal/logger"
	"github.com/daioru/todo-app/internal/pkg/db"
	"github.com/daioru/todo-app/internal/repository"
	"github.com/daioru/todo-app/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// @title TODO App API
// @version 1.0
// @description API Server for TODO Application
// @host localhost:8080
// @BasePath /api/
// @securityDefinitions.cookie Auth
// @in cookie
// @name Authorization
func main() {
	//Logger
	logger.InitLogger()
	log := logger.GetLogger()

	//Config
	if err := config.ReadConfigYML("config.yml"); err != nil {
		log.Fatal().Msg("Failed init configuration")
	}
	cfg := config.GetConfigInstance()

	//Database
	db, err := db.ConnectDB(&cfg.DB)
	if err != nil {
		log.Fatal().Msgf("sqlx_Open error: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal().Msgf("Error testing db connection: %v", err)
	}

	//Repositories
	userRepo := repository.NewUserRepository(db)
	taskRepo := repository.NewTaskRepository(db)

	//JWT
	err = godotenv.Load()
	if err != nil {
		log.Fatal().Msg("Error loading .env file")
	}

	jwtSecret := os.Getenv("JWTSECRET")
	if jwtSecret == "" {
		log.Fatal().Msg("No jwtSecret in .env")
	}

	//Services
	authService := services.NewAuthService(userRepo)
	taskService := services.NewTaskService(taskRepo)

	//Handlers
	authHandler := handlers.NewAuthHandler(authService)
	taskHandler := handlers.NewTaskHandler(taskService)

	handlers := handlers.NewHandlers(authHandler, taskHandler)

	//Server
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.SetTrustedProxies(nil)

	handlers.RegisterRoutes(r)

	port := ":8080"
	fmt.Println("Server is running on port" + port)
	if err := r.Run(port); err != nil {
		log.Fatal().Msgf("Failed to start server: %v", err)
	}
}
