package main

import (
	"github.com/drag0dev/argon/src/ecs/internal/handlers"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	r := gin.Default()

	api := r.Group("/api")
	handlers.RegisterRoutes(api)

	if err := r.Run(":8081"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
