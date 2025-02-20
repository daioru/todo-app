package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/daioru/todo-app/internal/config"
	"github.com/daioru/todo-app/internal/pkg/db"
	"github.com/gin-gonic/gin"
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

	r := gin.Default()
	r.SetTrustedProxies(nil)

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
