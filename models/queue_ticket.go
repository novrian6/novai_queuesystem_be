package models

import (
	"errors"
	"queue-system-backend/database"
	"time"

	"gorm.io/gorm"
)

type QueueTicket struct {
	TicketID      uint       `json:"ticket_id" gorm:"primaryKey;autoIncrement"`
	UserID        uint       `json:"user_id" gorm:"not null"`        // Foreign key to Users table
	ServiceID     *uint      `json:"service_id" gorm:"default:null"` // Foreign key to Services table
	CounterID     *uint      `json:"counter_id" gorm:"default:null"` // Foreign key to Counters table
	VenueID       *uint      `json:"venue_id" gorm:"default:null"`   // Foreign key to Venues table
	CustomerName  string     `json:"customer_name" gorm:"size:255"`
	CustomerEmail string     `json:"customer_email" gorm:"size:255"`
	CustomerPhone string     `json:"customer_phone" gorm:"size:20"`
	PhotoURL      string     `json:"photo_url" gorm:"size:255"`
	QueueNumber   string     `json:"queue_number" gorm:"size:10;not null"`
	Token         string     `json:"token" gorm:"size:255"` // New field for token
	Status        string     `json:"status" gorm:"type:enum('waiting','called','skipped','completed');default:waiting"`
	CreatedAt     time.Time  `json:"created_at" gorm:"autoCreateTime"`
	CalledAt      *time.Time `json:"called_at"`
	CompletedAt   *time.Time `json:"completed_at"`
	SkippedAt     *time.Time `json:"skipped_at"`  // New field for skipped timestamp
	OperatorID    *uint      `json:"operator_id"` // New field for operator ID
}

// TableName ensures GORM uses the correct table name
func (QueueTicket) TableName() string {
	return "QueueTickets"
}

// CreateQueueTicket inserts a new ticket
func CreateQueueTicket(ticket *QueueTicket) error {
	var service Service
	if err := database.DB.First(&service, ticket.ServiceID).Error; err != nil {
		return errors.New("service not found")
	}

	// Validate access to the service
	if service.UserID == nil || *service.UserID != ticket.UserID {
		return errors.New("unauthorized: cannot create ticket for this service")
	}

	return database.DB.Create(ticket).Error
}

// GetQueueTicketByID retrieves a ticket by ID and user ID
func GetQueueTicketByID(ticketID uint, userID uint, isAdmin bool) (*QueueTicket, error) {
	var ticket QueueTicket
	query := database.DB.Where("ticket_id = ?", ticketID)
	if !isAdmin {
		query = query.Where("user_id = ?", userID)
	}

	err := query.First(&ticket).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("ticket not found")
	}
	return &ticket, err
}

// GetQueueTicketByToken retrieves a ticket by its token
func GetQueueTicketByToken(token string) (*QueueTicket, error) {
	var ticket QueueTicket
	err := database.DB.Where("token = ?", token).First(&ticket).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("ticket not found")
	}
	return &ticket, err
}

// UpdateQueueTicket updates a ticket's details
func UpdateQueueTicket(ticket *QueueTicket, userID uint, isAdmin bool) error {
	if !isAdmin && ticket.UserID != userID {
		return errors.New("unauthorized: cannot update this ticket")
	}

	// Only update specific fields without overwriting `created_at`
	return database.DB.Model(&QueueTicket{}).
		Where("ticket_id = ? AND user_id = ?", ticket.TicketID, userID).
		Updates(map[string]interface{}{
			"status":       ticket.Status,
			"called_at":    ticket.CalledAt,
			"completed_at": ticket.CompletedAt,
		}).Error
}

// UpdateQueueTicketStatus updates the status of a ticket and sets the appropriate timestamp
func UpdateQueueTicketStatus(ticketID uint, status string, operatorID *uint, counterID *uint) error {
	var updates = map[string]interface{}{
		"status": status,
	}

	// Set the appropriate timestamp, operator_id, and counter_id based on the status
	switch status {
	case "called":
		updates["called_at"] = time.Now()
		updates["operator_id"] = operatorID
		updates["counter_id"] = counterID
	case "completed":
		updates["completed_at"] = time.Now()
		updates["operator_id"] = operatorID
		updates["counter_id"] = counterID
	case "skipped":
		updates["skipped_at"] = time.Now()
		updates["operator_id"] = operatorID
		updates["counter_id"] = counterID
	case "waiting":
		// No additional timestamp for waiting
	default:
		return errors.New("invalid status")
	}

	// Update the ticket in the database with operator_id in the WHERE clause for statuses other than "called"
	query := database.DB.Model(&QueueTicket{}).Where("ticket_id = ?", ticketID)
	//if status != "called" && operatorID != nil {
	//	query = query.Where("operator_id = ?", *operatorID)
	//}

	return query.Updates(updates).Error
}

// DeleteQueueTicket deletes a ticket by ID
func DeleteQueueTicket(ticketID uint, userID uint, isAdmin bool) error {
	query := database.DB.Where("ticket_id = ?", ticketID)
	if !isAdmin {
		query = query.Where("user_id = ?", userID)
	}

	return query.Delete(&QueueTicket{}).Error
}

// GetAllQueueTickets retrieves all tickets for a user or admin
func GetAllQueueTickets(userID uint, isAdmin bool) ([]QueueTicket, error) {
	var tickets []QueueTicket

	query := database.DB.Model(&QueueTicket{})
	if !isAdmin {
		query = query.Where("user_id = ?", userID)
	}

	err := query.Find(&tickets).Error
	return tickets, err
}

// GetQueueTicketsByStatus retrieves tickets filtered by status
func GetQueueTicketsByStatus(userID uint, status string, isAdmin bool) ([]QueueTicket, error) {
	var tickets []QueueTicket

	query := database.DB.Where("status = ?", status)
	if !isAdmin {
		query = query.Where("user_id = ?", userID)
	}

	err := query.Find(&tickets).Error
	return tickets, err
}

// GetLastQueueTicketByServiceID retrieves the last ticket for a specific service
func GetLastQueueTicketByServiceID(serviceID uint) (*QueueTicket, error) {
	var ticket QueueTicket
	err := database.DB.Where("service_id = ?", serviceID).
		Order("queue_number DESC").
		First(&ticket).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("record not found")
	}
	return &ticket, err
}

// IsTokenExists checks if a token already exists in the database
func IsTokenExists(token string) bool {
	var count int64
	database.DB.Model(&QueueTicket{}).Where("token = ?", token).Count(&count)
	return count > 0
}

// GetQueueTicketsSorted retrieves queue tickets sorted by queue_number in ascending order
func GetQueueTicketsSorted(status string, venueID uint, serviceID uint) ([]QueueTicket, error) {
	var tickets []QueueTicket
	query := database.DB.Where("venue_id = ? AND service_id = ?", venueID, serviceID)

	// If status is not "all", filter by status
	if status != "all" {
		query = query.Where("status = ?", status)
	}

	err := query.Order("queue_number ASC").Find(&tickets).Error
	if err != nil {
		return nil, err
	}
	return tickets, nil
}

func GetQueueTicketsByStatusVenueAndService(status string, venueID uint, serviceID uint, lastTopN *int) ([]QueueTicket, error) {
	var tickets []QueueTicket
	query := database.DB.Where("status = ? AND venue_id = ? AND service_id = ?", status, venueID, serviceID).
		Order("queue_number ASC")

	if lastTopN != nil && *lastTopN > 0 {
		query = query.Limit(*lastTopN)
	}

	err := query.Find(&tickets).Error
	if err != nil {
		return nil, err
	}
	return tickets, nil
}

// CalculateAverageQueuingTime calculates the average queuing time for completed tickets
func CalculateAverageQueuingTime(venueID uint, serviceID uint) (time.Duration, error) {
	var tickets []QueueTicket
	err := database.DB.Where("status = ? AND venue_id = ? AND service_id = ?", "completed", venueID, serviceID).Find(&tickets).Error
	if err != nil {
		return 0, err
	}

	var totalDuration time.Duration
	var count int
	for _, ticket := range tickets {
		if ticket.CompletedAt != nil {
			totalDuration += ticket.CompletedAt.Sub(ticket.CreatedAt)
			count++
		}
	}

	if count == 0 {
		return 0, errors.New("no completed tickets found")
	}

	averageDuration := totalDuration / time.Duration(count)
	return time.Duration(averageDuration.Minutes()), nil
}
