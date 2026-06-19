package routes

import (
	"github.com/Pedro-Wilker/api-eventos/controllers"
	"github.com/Pedro-Wilker/api-eventos/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "API online e pronta para o trabalho!"})
	})

	api := r.Group("/api")
	{
		api.POST("/register", controllers.Register)
		api.POST("/login", controllers.Login)

		protected := api.Group("/")
		protected.Use(middlewares.AuthMiddleware())
		{
			protected.GET("/me", func(c *gin.Context) {
				userID, _ := c.Get("userID")
				role, _ := c.Get("role")
				c.JSON(200, gin.H{
					"message": "Você está autenticado!",
					"user_id": userID,
					"role":    role,
				})
			})

			protected.POST("/convidados", controllers.CreateGuest)
			protected.GET("/convidados", controllers.ListGuests)

			protected.POST("/presentes", controllers.CreateGift)
			protected.GET("/presentes", controllers.ListGifts)
		}
	}
}
