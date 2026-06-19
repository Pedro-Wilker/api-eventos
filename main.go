package main

import (
	"github.com/gin-gonic/gin"
	"github.com/seunome/api-eventos/config"
)

func main() {
	config.ConnectDB()

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "API online!"})
	})

	r.Run(":8080")
}
