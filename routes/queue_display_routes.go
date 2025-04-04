package routes

import (
	"queue-system-backend/controllers"
	"queue-system-backend/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterQueueDisplayRoutes(router *gin.Engine) {
	displayRoutes := router.Group("/display").Use(middlewares.AuthMiddleware())

	{
		displayRoutes.GET("/:counter_id", controllers.GetQueueDisplay)
		displayRoutes.POST("/create", controllers.CreateQueueDisplay)
		displayRoutes.PUT("/update", controllers.UpdateQueueDisplay)
		displayRoutes.PUT("/:counter_id/next", controllers.AssignNextTicket)
		displayRoutes.GET("/all", controllers.GetAllQueueDisplays)
		displayRoutes.GET("/next-counter", controllers.GetNextCounterController)
		displayRoutes.POST("/reset", controllers.ResetQueueDisplayHandler)
		displayRoutes.GET("/analytics", controllers.GetDisplayAnalytics)
		displayRoutes.GET("/current-ticket", controllers.GetCurrentTicket)

	}
}
