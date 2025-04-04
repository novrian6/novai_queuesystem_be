package routes

import (
	"queue-system-backend/controllers"
	"queue-system-backend/middlewares"

	"github.com/gin-gonic/gin"
)

func QueueTicketRoutes(router *gin.Engine) {
	tickets := router.Group("/queue-tickets").Use(middlewares.AuthMiddleware())
	{
		tickets.GET("/:id", controllers.GetQueueTicketByIDHandler)
		tickets.POST("/", controllers.CreateQueueTicketHandler)
		tickets.PUT("/:id", controllers.UpdateQueueTicketHandler)
		tickets.DELETE("/:id", controllers.DeleteQueueTicketHandler)
		tickets.PUT("/:id/status", controllers.UpdateQueueTicketStatusHandler)
	}
}
