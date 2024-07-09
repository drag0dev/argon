package handlers

import (
	"github.com/drag0dev/argon/src/ecs/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/subscriptions", HandleSubscriptions)
	r.GET("/movies", HandleMovies)
	r.GET("/shows", HandleShows)
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

func HandleMovies(c *gin.Context) {
	movies, err := services.GetMovies()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"movies": movies})
	return
}

func HandleShows(c *gin.Context) {
	shows, err := services.GetShows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"shows": shows})
	return
}
