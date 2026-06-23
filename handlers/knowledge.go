package handlers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"biubiu/config"
)

func GetKnowledge(c *gin.Context) {
	data, err := os.ReadFile(config.DataDir + "/knowledge.json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "knowledge file not found"})
		return
	}
	c.Data(http.StatusOK, "application/json; charset=utf-8", data)
}
