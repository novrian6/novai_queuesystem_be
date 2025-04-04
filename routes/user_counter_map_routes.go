package routes

import (
	"queue-system-backend/controllers"
	"queue-system-backend/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterUserCounterMapRoutes(router *gin.Engine) {
	userCounterMapRoutes := router.Group("/user-counter-map").Use(middlewares.AuthMiddleware())
	{
		userCounterMapRoutes.GET("/", controllers.GetAllUserCounterMaps)
		userCounterMapRoutes.GET("/:id", controllers.GetUserCounterMapByID)
		userCounterMapRoutes.POST("/", controllers.CreateUserCounterMap)
		userCounterMapRoutes.DELETE("/:id", controllers.DeleteUserCounterMap)
		userCounterMapRoutes.GET("/user-by-counter/:counter_id", controllers.GetUserIDByCounterID)
		userCounterMapRoutes.GET("/counter-by-user/:user_id", controllers.GetCounterIDByUserID)
	}
}
