package routes

import (
	"queue-system-backend/controllers"
	"queue-system-backend/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterCounterRoutes(router *gin.Engine) {
	counters := router.Group("/counters").Use(middlewares.AuthMiddleware())
	{
		counters.POST("/", controllers.CreateCounter)
		counters.GET("/:id", controllers.GetCounter)
		counters.GET("/", controllers.GetCounters)
		counters.PUT("/:id", controllers.UpdateCounter)
		counters.DELETE("/:id", controllers.DeleteCounter)
		//counters.GET("/company/:company_id", controllers.GetCountersByCompany)
		counters.GET("/venue/:venue_id", controllers.GetCountersByVenue)
	}
}
