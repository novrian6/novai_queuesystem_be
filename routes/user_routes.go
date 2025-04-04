package routes

import (
	"queue-system-backend/controllers"
	"queue-system-backend/middlewares"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	users := router.Group("/users").Use(middlewares.AuthMiddleware())
	{
		// List all users (Admin only)
		users.GET("", middlewares.AdminMiddleware(), controllers.ListUsers)

		// Create a new user (Admin only)
		users.POST("/", middlewares.AdminMiddleware(), controllers.CreateUser)

		// Get a specific user by ID (Accessible by the user themselves or Admin)
		users.GET("/:id", controllers.GetUser)

		// Update a user by ID (Admin only)
		users.PUT("/:id", middlewares.AdminMiddleware(), controllers.UpdateUser)

		// Delete a user by ID (Admin only)
		users.DELETE("/:id", middlewares.AdminMiddleware(), controllers.DeleteUser)

	}
}
