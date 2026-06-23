package main

import (
	"log"
	"net/http"

	"biubiu/config"
	"biubiu/database"
	"biubiu/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	database.Init(config.DBPath)
	defer database.Close()
	database.SeedData(config.DataDir)

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	r.Static("/static", "./web")
	r.StaticFile("/", "./web/index.html")

	api := r.Group("/api")
	{
		api.GET("/weapons", handlers.GetWeapons)
		api.GET("/weapons/top10", handlers.GetTop10)
		api.GET("/weapons/:id", handlers.GetWeaponByID)
		api.POST("/weapons/compare", handlers.CompareWeapons)
		api.GET("/attachments", handlers.GetAttachments)
		api.GET("/mod-codes", handlers.GetModCodes)
		api.POST("/recommend", handlers.Recommend)
		api.GET("/knowledge", handlers.GetKnowledge)
	}

	log.Printf("Server starting on %s", config.ServerPort)
	log.Printf("Open http://localhost%s in browser", config.ServerPort)
	r.Run(config.ServerPort)
}
