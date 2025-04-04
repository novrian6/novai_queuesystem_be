package routes

import (
	"queue-system-backend/controllers"
	"queue-system-backend/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupStatisticsRoutes(r *gin.Engine) {
	statsController := controllers.NewStatisticsController()
	statistics := r.Group("/statistics").Use(middlewares.AuthMiddleware())
	{
		statistics.POST("/active-queues", statsController.GetActiveQueues)
		statistics.POST("/average-wait-time", statsController.GetAverageWaitTime)
		statistics.POST("/total-served", statsController.GetTotalServed)
	}
}
