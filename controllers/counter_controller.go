package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"queue-system-backend/models"
	"queue-system-backend/utils"

	"github.com/gin-gonic/gin"
)

// CreateCounter adds a new counter
func CreateCounter(c *gin.Context) {
	var counter models.Counter
	if err := c.ShouldBindJSON(&counter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	userID := c.GetUint("user_id")

	// Verify service ownership before creation
	service, err := models.GetServiceByID(*counter.ServiceID)
	if err != nil || service.UserID == nil || *service.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized: Service does not belong to your account"})
		return
	}

	// Validate time format
	const timeFormat = "15:04:05"
	if _, err := time.Parse(timeFormat, counter.OpenTime); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid open_time format, expected HH:mm:ss"})
		return
	}
	if _, err := time.Parse(timeFormat, counter.CloseTime); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid close_time format, expected HH:mm:ss"})
		return
	}

	// Set the UserID from the token claims
	counter.UserID = userID

	if err := counter.CreateCounter(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create counter", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, counter)
}

// GetCounter retrieves a specific counter by ID
func GetCounter(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetUint("user_id")
	isAdmin := c.GetString("role") == "admin"

	counter, err := models.GetCounterByID(uint(id), userID, isAdmin)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Counter not found"})
		return
	}

	c.JSON(http.StatusOK, counter)
}

// GetCounters retrieves all counters for the user or admin
func GetCounters(c *gin.Context) {
	//userID := c.GetUint("user_id")
	isAdmin := c.GetString("role") == "admin"

	var counters []models.Counter
	var err error

	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		return
	}

	userClaims, ok := claims.(*utils.Claims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid token claims"})
		return
	}

	if strings.ToLower(userClaims.Role) == "operator" {
		user, err := models.GetUserByID(userClaims.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user information"})
			c.Abort()
			return
		}

		counters, err = models.GetAllCounters(*user.OwnerID, false)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch venues"})
			return
		}

	} else {

		counters, err = models.GetAllCounters(userClaims.UserID, isAdmin)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch venues"})
			return
		}

	}

	c.JSON(http.StatusOK, counters)
}

// UpdateCounter updates an existing counter
func UpdateCounter(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetUint("user_id")
	isAdmin := c.GetString("role") == "admin"

	counter, err := models.GetCounterByID(uint(id), userID, isAdmin)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Counter not found"})
		return
	}

	if err := c.ShouldBindJSON(&counter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate time format
	const timeFormat = "15:04:05"
	if _, err := time.Parse(timeFormat, counter.OpenTime); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid open_time format, expected HH:mm:ss"})
		return
	}
	if _, err := time.Parse(timeFormat, counter.CloseTime); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid close_time format, expected HH:mm:ss"})
		return
	}

	if err := counter.UpdateCounter(userID, isAdmin); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update counter", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, counter)
}

// DeleteCounter deletes a counter
func DeleteCounter(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetUint("user_id")
	isAdmin := c.GetString("role") == "admin"

	if err := models.DeleteCounter(uint(id), userID, isAdmin); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete counter", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Counter deleted successfully"})
}

// GetCountersByVenue retrieves all counters for a specific venue
func GetCountersByVenue(c *gin.Context) {
	venueID, _ := strconv.Atoi(c.Param("venue_id"))

	counters, err := models.GetCountersByVenue(uint(venueID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve counters", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, counters)
}
