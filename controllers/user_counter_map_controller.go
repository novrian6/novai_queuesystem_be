package controllers

import (
	"net/http"
	"strconv"
	"time"

	"queue-system-backend/models"

	"github.com/gin-gonic/gin"
)

func GetAllUserCounterMaps(c *gin.Context) {
	// Extract user_id from the request context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userIDUInt, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	// Use FetchByOwnerID from the model
	mappings, err := models.FetchByOwnerID(userIDUInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch mappings", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, mappings)
}

func GetUserCounterMapByID(c *gin.Context) {
	id := c.Param("id")

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	mapping, err := models.FetchByID(uint(idUint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mapping not found"})
		return
	}

	c.JSON(http.StatusOK, mapping)
}

func CreateUserCounterMap(c *gin.Context) {
	var req struct {
		UserID    int `json:"user_id" binding:"required"`
		CounterID int `json:"counter_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retrieve userID from context and ensure it's a uint
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userIDUInt, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	mapping := models.UserCounterMap{
		UserID:     req.UserID,
		CounterID:  req.CounterID,
		OwnerID:    userIDUInt, // Use the properly casted userID
		AssignedAt: time.Now(), // Set the current timestamp for AssignedAt
	}

	if err := mapping.Create(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create mapping"})
		return
	}

	c.JSON(http.StatusCreated, mapping)
}

func UpdateUserCounterMap(c *gin.Context) {
	id := c.Param("id")

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// Extract user_id from the request context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userIDUInt, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	mapping, err := models.FetchByID(uint(idUint))
	if err != nil || mapping.OwnerID != userIDUInt {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mapping not found or unauthorized"})
		return
	}

	// Define a struct for binding the request payload
	var req struct {
		UserID    int `json:"user_id"`
		CounterID int `json:"counter_id"`
	}

	// Bind the updated data
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	// Update fields if provided
	if req.UserID != 0 {
		mapping.UserID = req.UserID
	}
	if req.CounterID != 0 {
		mapping.CounterID = req.CounterID
	}

	// Validate ownership using the model method
	if err := mapping.ValidateOwnership(); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	// Use Update from the model
	if err := mapping.Update(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update mapping", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Mapping updated successfully", "mapping": mapping})
}

func DeleteUserCounterMap(c *gin.Context) {
	id := c.Param("id")

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	mapping, err := models.FetchByID(uint(idUint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mapping not found"})
		return
	}

	if err := mapping.Delete(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete mapping"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Mapping deleted successfully"})
}

// GetUserIDByCounterID handles the request to fetch user_id by counter_id
func GetUserIDByCounterID(c *gin.Context) {
	counterID, err := strconv.Atoi(c.Param("counter_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid counter_id"})
		return
	}

	userID, err := models.GetUserIDByCounterID(counterID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user_id": userID})
}

// GetCounterIDByUserID handles the request to fetch counter_id by user_id
func GetCounterIDByUserID(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
		return
	}

	counterID, err := models.GetCounterIDByUserID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"counter_id": counterID})
}
