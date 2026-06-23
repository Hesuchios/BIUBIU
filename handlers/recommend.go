package handlers

import (
	"net/http"

	"biubiu/models"
	"biubiu/services"

	"github.com/gin-gonic/gin"
)

func Recommend(c *gin.Context) {
	var req models.RecommendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	results, err := services.GetRecommendations(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}
