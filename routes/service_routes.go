package routes

import (
	"queue-system-backend/controllers"
	"queue-system-backend/middlewares"

	"github.com/gin-gonic/gin"
)

func ServiceRoutes(router *gin.Engine) {
	services := router.Group("/services").Use(middlewares.AuthMiddleware())
	{
		services.GET("", controllers.ListServices)
		services.GET("/:id", controllers.GetService)
		services.GET("/venue/:venue_id", controllers.ListServicesByVenueAndUser)
		services.POST("", middlewares.AdminMiddleware(), controllers.CreateService)
		services.PUT("/:id", middlewares.AdminMiddleware(), controllers.UpdateService)
		services.DELETE("/:id", middlewares.AdminMiddleware(), controllers.DeleteService)
	}
}
