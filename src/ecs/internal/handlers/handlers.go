package handlers

import (
	"github.com/drag0dev/argon/src/ecs/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/hello", HelloHandler)
}

func HelloHandler(c *gin.Context) {
	message := services.GetHelloMessage()
	c.JSON(http.StatusOK, gin.H{"message": message})
}
