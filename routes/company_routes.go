package routes

import (
	"queue-system-backend/controllers"

	"github.com/gin-gonic/gin"
)

// CompanyRoutes registers company-related routes
func CompanyRoutes(router *gin.Engine) {
	companyGroup := router.Group("/companies")
	{
		companyGroup.GET("/:id", controllers.GetCompanyByIDHandler)
		companyGroup.POST("/", controllers.CreateCompanyHandler)
		companyGroup.PUT("/", controllers.UpdateCompanyHandler)
		companyGroup.DELETE("/:id", controllers.DeleteCompanyHandler)
	}
}
