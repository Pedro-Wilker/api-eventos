package main

import (
	"strings"
	"time"

	"github.com/Pedro-Wilker/api-eventos/config"
	"github.com/Pedro-Wilker/api-eventos/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDB()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
		AllowOriginFunc: func(origin string) bool {
			if origin == "https://casamentopedroelara.com.br" || strings.HasSuffix(origin, ".casamentopedroelara.com.br") {
				return true
			}

			if origin == "https://artonbyte.com.br" || strings.HasSuffix(origin, ".artonbyte.com.br") {
				return true
			}

			if origin == "http://localhost:3000" || origin == "http://localhost:5173" || origin == "http://127.0.0.1:3000" {
				return true
			}

			return false
		},
	}))

	routes.SetupRoutes(r)

	r.Run(":8080")
}
