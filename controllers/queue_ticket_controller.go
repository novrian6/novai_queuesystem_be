// controllers/queue_ticket_controller.go
package controllers

import (
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"strconv"
	"strings"
	"time"

	"queue-system-backend/models"

	"github.com/gin-gonic/gin"
)

func GetQueueTicketByIDHandler(c *gin.Context) {
	ticketID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	userID := c.GetUint("user_id")
	isAdmin := c.GetString("role") == "admin"

	ticket, err := models.GetQueueTicketByID(uint(ticketID), userID, isAdmin)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ticket)
}

func CreateQueueTicketHandler(c *gin.Context) {
	userID := c.GetUint("user_id")

	// Parse POST input
	var input struct {
		VenueID       uint   `json:"venue_id" binding:"required"`
		ServiceID     uint   `json:"service_id" binding:"required"`
		CustomerName  string `json:"customer_name" binding:"required"`
		CustomerEmail string `json:"customer_email" binding:"required,email"`
		CustomerPhone string `json:"customer_phone" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Validate service ownership
	service, err := models.GetServiceByID(input.ServiceID)
	if err != nil || service.UserID == nil || *service.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized: Service does not belong to your account"})
		return
	}

	// Generate QueueNumber
	lastTicket, err := models.GetLastQueueTicketByServiceID(input.ServiceID)
	if err != nil && err.Error() != "record not found" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate queue number", "details": err.Error()})
		return
	}
	newQueueNumber := 1
	if lastTicket != nil {
		lastNumber, _ := strconv.Atoi(lastTicket.QueueNumber)
		newQueueNumber = lastNumber + 1
	}

	// Generate a unique token (example using a random string)
	token := generateToken()

	// Create ticket
	ticket := models.QueueTicket{
		UserID:        userID,
		ServiceID:     &input.ServiceID,
		VenueID:       &input.VenueID,
		CustomerName:  input.CustomerName,
		CustomerEmail: input.CustomerEmail,
		CustomerPhone: input.CustomerPhone,
		QueueNumber:   strconv.Itoa(newQueueNumber),
		Status:        "waiting", // Default status
		Token:         token,     // Set the generated token
	}

	if err := models.CreateQueueTicket(&ticket); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create ticket", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, ticket)
}

// Helper function to generate a unique token using a hash and Base62 encoding
func generateToken() string {
	// Use a combination of timestamp and random data for uniqueness
	data := strconv.FormatInt(time.Now().UnixNano(), 10)

	// Generate a SHA256 hash of the data
	hash := sha256.Sum256([]byte(data))

	// Encode the hash in Base64
	base64Encoded := base64.StdEncoding.EncodeToString(hash[:])

	// Convert Base64 to Base62 by removing non-alphanumeric characters
	base62Encoded := strings.Map(func(r rune) rune {
		if (r >= '0' && r <= '9') || (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
			return r
		}
		return -1
	}, base64Encoded)

	// Return a shortened Base62-encoded string (e.g., first 8 characters)
	return base62Encoded[:8]
}

func UpdateQueueTicketHandler(c *gin.Context) {
	ticketID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	var ticket models.QueueTicket
	if err := c.ShouldBindJSON(&ticket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("user_id")
	isAdmin := c.GetString("role") == "admin"

	ticket.TicketID = uint(ticketID)

	if err := models.UpdateQueueTicket(&ticket, userID, isAdmin); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update ticket", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ticket updated successfully"})
}

func DeleteQueueTicketHandler(c *gin.Context) {
	ticketID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	userID := c.GetUint("user_id")
	isAdmin := c.GetString("role") == "admin"

	if err := models.DeleteQueueTicket(uint(ticketID), userID, isAdmin); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete ticket", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ticket deleted successfully"})
}

func GetQueueTicketsByStatusHandler(c *gin.Context) {
	status := c.Query("status")
	userID := c.GetUint("user_id")
	isAdmin := c.GetString("role") == "admin"

	tickets, err := models.GetQueueTicketsByStatus(userID, status, isAdmin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tickets", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tickets)
}

func GetAllQueueTicketsHandler(c *gin.Context) {
	userID := c.GetUint("user_id")
	isAdmin := c.GetString("role") == "admin"

	tickets, err := models.GetAllQueueTickets(userID, isAdmin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tickets", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tickets)
}

func UpdateQueueTicketStatusHandler(c *gin.Context) {
	ticketID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	var input struct {
		Status     string `json:"status" binding:"required"`
		OperatorID *uint  `json:"operator_id"` // Optional operator ID
		CounterID  *uint  `json:"counter_id"`  // Optional counter ID
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	if err := models.UpdateQueueTicketStatus(uint(ticketID), input.Status, input.OperatorID, input.CounterID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update ticket status", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ticket status updated successfully"})
}
