package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
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
