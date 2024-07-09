package handlers

import (
	"github.com/drag0dev/argon/src/ecs/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/subscriptions", HandleSubscriptions)
}

func HandleSubscriptions(c *gin.Context) {
	userId := c.Query("userId")
	if len(userId) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Malformed input."})
		return
	}

	subscriptions, err := services.GetSubscriptions(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"subscriptions": subscriptions})
	return
}
