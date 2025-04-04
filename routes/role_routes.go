package routes

import (
	"queue-system-backend/controllers"

	"github.com/gin-gonic/gin"
)

// RoleRoutes registers role-related routes
func RoleRoutes(router *gin.Engine) {
	roleGroup := router.Group("/roles")
	{
		roleGroup.GET("/", controllers.ListRolesHandler)
		roleGroup.GET("/:id", controllers.GetRoleByIDHandler)
		roleGroup.POST("/", controllers.CreateRoleHandler)
		roleGroup.PUT("/", controllers.UpdateRoleHandler)
		roleGroup.DELETE("/:id", controllers.DeleteRoleHandler)
	}
}
