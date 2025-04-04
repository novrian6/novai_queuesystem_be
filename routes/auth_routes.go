package routes

import (
	"queue-system-backend/controllers"
	"queue-system-backend/middlewares"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
		auth.POST("/logout", controllers.Logout)
		auth.GET("/me", middlewares.AuthMiddleware(), controllers.Me)
		auth.POST("/forgot-password", controllers.ForgotPassword)
	}
}
