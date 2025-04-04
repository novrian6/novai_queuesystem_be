package controllers

import (
	"net/http"
	"queue-system-backend/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetCompanyByIDHandler handles GET requests for a company by ID
func GetCompanyByIDHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	companyIDValue, exists := c.Get("company_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "company_id not found in context"})
		return
	}
	companyID, ok := companyIDValue.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid company ID format"})
		return
	}
	if companyID != uint(id) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only view your own company"})
		return
	}
	company, err := models.GetCompanyByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, company)
}

// CreateCompanyHandler handles POST requests to create a company
func CreateCompanyHandler(c *gin.Context) {
	var company models.Company
	if err := c.ShouldBindJSON(&company); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.CreateCompany(&company); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, company)
}

// UpdateCompanyHandler handles PUT requests to update a company
func UpdateCompanyHandler(c *gin.Context) {
	var company models.Company

	// Ensure the user can only update their own company
	userCompanyIDValue, exists := c.Get("company_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "company_id not found in context"})
		return
	}

	userCompanyID, ok := userCompanyIDValue.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid company ID format"})
		return
	}
	if company.CompanyID != userCompanyID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only update your own company"})
		return
	}

	if err := c.ShouldBindJSON(&company); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.UpdateCompany(&company); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, company)
}

// DeleteCompanyHandler handles DELETE requests to delete a company
func DeleteCompanyHandler(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	if err := models.DeleteCompany(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Company deleted successfully"})
}
