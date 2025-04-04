package controllers

import (
	"log"
	"net/http"
	"queue-system-backend/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateQueueDisplay creates a new QueueDisplay
func CreateQueueDisplay(c *gin.Context) {
	var display models.QueueDisplay
	if err := c.ShouldBindJSON(&display); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retrieve the logged-in user's ID from the context
	userID, exists := c.Get("user_id") //userID := c.GetUint("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Set the user_id of the display to the logged-in user's ID
	display.UserID = userID.(uint)

	if err := models.CreateQueueDisplay(&display); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create display"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Display created successfully", "data": display})
}

// GetQueueDisplay retrieves QueueDisplay by CounterID
func GetQueueDisplay(c *gin.Context) {
	counterID, err := strconv.Atoi(c.Param("counter_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid counter ID"})
		return
	}

	display, err := models.GetQueueDisplayByCounterID(uint(counterID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Display not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": display})
}

// UpdateQueueDisplay updates a QueueDisplay entry
func UpdateQueueDisplay(c *gin.Context) {
	var display models.QueueDisplay
	if err := c.ShouldBindJSON(&display); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.UpdateQueueDisplay(&display); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update display"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Display updated successfully"})
}

// AssignNextTicket assigns the next ticket to `CurrentTicket`
func AssignNextTicket(c *gin.Context) {
	counterID, err := strconv.Atoi(c.Param("counter_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid counter ID"})
		return
	}

	display, err := models.GetQueueDisplayByCounterID(uint(counterID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Queue display not found"})
		return
	}

	if err := display.AutoAssignNextTicket(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Next ticket assigned", "data": display})
}

// GetAllQueueDisplays retrieves all QueueDisplays with optional filters
func GetAllQueueDisplays(c *gin.Context) {
	var venueID, serviceID *uint

	if vID, err := strconv.Atoi(c.Query("venue_id")); err == nil {
		vIDUint := uint(vID)
		venueID = &vIDUint
	}
	userIDValue, exists := c.Get("user_id")
	var userID *uint
	if exists {
		uIDUint := userIDValue.(uint)
		userID = &uIDUint
	}
	if sID, err := strconv.Atoi(c.Query("service_id")); err == nil {
		sIDUint := uint(sID)
		serviceID = &sIDUint
	}

	displays, err := models.GetQueueDisplays(venueID, userID, serviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, displays)
}

// GetNextCounterController handles fetching the next available counter
func GetNextCounterController(c *gin.Context) {
	venueID, _ := strconv.Atoi(c.Query("venue_id"))
	serviceID, _ := strconv.Atoi(c.Query("service_id"))

	displays, err := models.GetNextCounter(uint(venueID), uint(serviceID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": displays})
}

// ResetQueueDisplayHandler resets all QueueDisplays
func ResetQueueDisplayHandler(c *gin.Context) {
	err := models.ResetQueueDisplay()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset display"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Display successfully reset"})
}

// GetDisplayAnalytics retrieves analytics for a specific venue and service
func GetDisplayAnalytics(c *gin.Context) {
	venueID, _ := strconv.Atoi(c.Query("venue_id"))
	serviceID, _ := strconv.Atoi(c.Query("service_id"))
	date := c.Query("date")

	analytics, err := models.GetDisplayAnalytics(uint(venueID), uint(serviceID), date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, analytics)
}

// GetCurrentTicket retrieves the current ticket for a specific counter
func GetCurrentTicket(c *gin.Context) {
	// Parse query parameters
	venueID := c.DefaultQuery("venue_id", "0")
	serviceID := c.DefaultQuery("service_id", "0")
	counterID := c.DefaultQuery("counter_id", "")

	// Convert venue_id, service_id to uint and counter_id to int
	var venueIDUint, serviceIDUint uint
	var counterIDInt int

	// Convert to integers and handle conversion errors
	if venueID != "0" {
		venueIDUint = parseUint(venueID)
	}
	if serviceID != "0" {
		serviceIDUint = parseUint(serviceID)
	}
	if counterID != "" {
		counterIDInt = parseInt(counterID)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "counter_id is required"})
		return
	}

	// Fetch current ticket from the model
	queueDisplay, err := models.GetCurrentTicket(venueIDUint, serviceIDUint, uint(counterIDInt))
	if err != nil {
		log.Println("Error fetching current ticket:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get current ticket"})
		return
	}

	// Return the result
	c.JSON(http.StatusOK, gin.H{
		"current_ticket": queueDisplay.CurrentTicket,
		"counter_id":     queueDisplay.CounterID,
		"status":         "success",
		"service_id":     queueDisplay.ServiceID,
		"venue_id":       queueDisplay.VenueID,
	})
}

// Helper function to parse uint
func parseUint(val string) uint {
	uintVal, err := strconv.ParseUint(val, 10, 32)
	if err != nil {
		return 0
	}
	return uint(uintVal)
}

// Helper function to parse int
func parseInt(val string) int {
	intVal, err := strconv.Atoi(val)
	if err != nil {
		return 0
	}
	return intVal
}
