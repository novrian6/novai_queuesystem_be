package routes

import (
	"queue-system-backend/controllers"
	"queue-system-backend/middlewares"

	"github.com/gin-gonic/gin"
)

func VenueRoutes(router *gin.Engine) {
	venues := router.Group("/venues").Use(middlewares.AuthMiddleware())
	{
		venues.GET("", controllers.ListVenues)
		venues.POST("", middlewares.AdminMiddleware(), controllers.CreateVenue)
		venues.GET("/:id", controllers.GetVenue)
		venues.PUT("/:id", middlewares.AdminMiddleware(), controllers.UpdateVenue)
		venues.DELETE("/:id", middlewares.AdminMiddleware(), controllers.DeleteVenue)
	}
}
