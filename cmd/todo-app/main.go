package main

import (
	"fmt"
	"os"

	"github.com/daioru/todo-app/internal/config"
	"github.com/daioru/todo-app/internal/handlers"
	"github.com/daioru/todo-app/internal/logger"
	"github.com/daioru/todo-app/internal/middlewares"
	"github.com/daioru/todo-app/internal/pkg/db"
	"github.com/daioru/todo-app/internal/repository"
	"github.com/daioru/todo-app/internal/services"
	"github.com/daioru/todo-app/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

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

	//JWT Service
	err = godotenv.Load()
	if err != nil {
		log.Fatal().Msg("Error loading .env file")
	}

	jwtSecret := os.Getenv("JWTSECRET")
	if jwtSecret == "" {
		log.Fatal().Msg("No jwtSecret in .env")
	}

	jwtService := utils.NewJwtService([]byte(jwtSecret))

	//Services
	authService := services.NewAuthService(userRepo, jwtService)
	taskService := services.NewTaskService(taskRepo)

	//Handlers
	authHandler := handlers.NewAuthHandler(authService)
	taskHandler := handlers.NewTaskHandler(taskService)
	authMiddleware := middlewares.NewAuthMiddleware(jwtService)

	handlers := handlers.NewHandlers(authHandler, taskHandler, authMiddleware)

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
