package routes

import (
	"queue-system-backend/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterViewTicketRoutes(router *gin.Engine) {
	router.GET("/myticket/:token", controllers.GetQueueTicketByTokenHandler)
	router.GET("/waiting-tickets", controllers.GetWaitingQueueTicketsHandler) // New route
}
