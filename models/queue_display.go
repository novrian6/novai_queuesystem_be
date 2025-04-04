package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"queue-system-backend/database"
	"time"
)

// QueueDisplay Model
type QueueDisplay struct {
	DisplayID     uint      `json:"display_id" gorm:"primaryKey;autoIncrement"`
	VenueID       uint      `json:"venue_id" gorm:"not null"`      // Foreign key to Venues table
	UserID        uint      `json:"user_id" gorm:"not null"`       // Foreign key to Users table
	ServiceID     uint      `json:"service_id" gorm:"not null"`    // Foreign key to Services table
	CounterID     uint      `json:"counter_id" gorm:"not null"`    // Foreign key to Counters table
	CurrentTicket string    `json:"current_ticket" gorm:"size:10"` // Current ticket being served
	NextTickets   string    `json:"next_tickets" gorm:"type:json"` // Stored as JSON string
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName ensures GORM uses the correct table name
func (QueueDisplay) TableName() string {
	return "QueueDisplay"
}

// CreateQueueDisplay creates a new display entry
func CreateQueueDisplay(display *QueueDisplay) error {
	return database.DB.Create(display).Error
}

// GetQueueDisplayByCounterID retrieves display data by CounterID
func GetQueueDisplayByCounterID(counterID uint) (*QueueDisplay, error) {
	var display QueueDisplay
	if err := database.DB.Where("counter_id = ?", counterID).First(&display).Error; err != nil {
		return nil, errors.New("queue display not found")
	}
	return &display, nil
}

// UpdateQueueDisplay updates display details
func UpdateQueueDisplay(display *QueueDisplay) error {
	return database.DB.Save(display).Error
}

// AutoAssignNextTicket updates `CurrentTicket` and removes it from `NextTickets`
func (qd *QueueDisplay) AutoAssignNextTicket() error {
	var nextTickets []string
	if err := json.Unmarshal([]byte(qd.NextTickets), &nextTickets); err != nil {
		return errors.New("failed to parse next_tickets")
	}

	if len(nextTickets) == 0 {
		return errors.New("no next tickets available")
	}

	// Move the first ticket to CurrentTicket
	qd.CurrentTicket = nextTickets[0]

	// Remove the first ticket from the list
	nextTickets = nextTickets[1:]

	// Convert back to JSON string
	updatedTickets, err := json.Marshal(nextTickets)
	if err != nil {
		return errors.New("failed to update next_tickets")
	}

	qd.NextTickets = string(updatedTickets)

	return database.DB.Save(qd).Error
}

// GetQueueDisplays retrieves QueueDisplay entries with optional filters
func GetQueueDisplays(venueID, userID, serviceID *uint) ([]QueueDisplay, error) {
	var displays []QueueDisplay
	query := database.DB.Model(&QueueDisplay{})

	if venueID != nil {
		query = query.Where("venue_id = ?", *venueID)
	}
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	if serviceID != nil {
		query = query.Where("service_id = ?", *serviceID)
	}

	err := query.Find(&displays).Error
	return displays, err
}

// GetNextCounter finds the next available counter ready for serving the next ticket
func GetNextCounter(venueID uint, serviceID uint) ([]QueueDisplay, error) {
	var displays []QueueDisplay

	query := database.DB.Model(&QueueDisplay{})
	if venueID != 0 {
		query = query.Where("venue_id = ?", venueID)
	}
	if serviceID != 0 {
		query = query.Where("service_id = ?", serviceID)
	}

	err := query.Find(&displays).Error
	if err != nil {
		return nil, errors.New("no available counters found")
	}
	return displays, nil
}

// ResetQueueDisplay clears all tickets on the display
func ResetQueueDisplay() error {
	return database.DB.Model(&QueueDisplay{}).
		Where("1 = 1").
		Updates(map[string]interface{}{
			"current_ticket": "",
			"next_tickets":   "[]",
		}).Error
}

// GetDisplayAnalytics retrieves analytics for a specific venue and service
type DisplayAnalytics struct {
	TotalCalled  int    `json:"total_called"`
	TotalInQueue int    `json:"total_in_queue"`
	TopCounter   string `json:"top_counter"`
}

func GetDisplayAnalytics(venueID, serviceID uint, date string) (*DisplayAnalytics, error) {
	var analytics DisplayAnalytics
	var dateFilter string

	if date != "" {
		parsedDate, _ := time.Parse("2006-01-02", date)
		dateFilter = parsedDate.Format("2006-01-02")
	} else {
		dateFilter = time.Now().Format("2006-01-02")
	}

	query := database.DB.Table("QueueTickets").Select(
		"COUNT(CASE WHEN status = 'completed' THEN 1 ELSE NULL END) AS total_called",
		"COUNT(CASE WHEN status = 'waiting' THEN 1 ELSE NULL END) AS total_in_queue",
		"MAX(counter_id) AS top_counter").
		Where("DATE(created_at) = ?", dateFilter)

	if venueID != 0 {
		query = query.Where("venue_id = ?", venueID)
	}
	if serviceID != 0 {
		query = query.Where("service_id = ?", serviceID)
	}

	err := query.Scan(&analytics).Error
	if err != nil {
		return nil, err
	}

	return &analytics, nil
}

// GetCurrentTicket retrieves the current ticket for a specific counter
func GetCurrentTicket(venueID, serviceID, counterID uint) (*QueueDisplay, error) {
	var queueDisplay QueueDisplay

	// Build the query
	query := database.DB.Table("QueueDisplay").
		Select("display_id, venue_id, user_id, service_id, counter_id, current_ticket, next_tickets, updated_at").
		Where("counter_id = ?", counterID)

	// Optional filters for venue_id and service_id
	if venueID != 0 {
		query = query.Where("venue_id = ?", venueID)
	}
	if serviceID != 0 {
		query = query.Where("service_id = ?", serviceID)
	}

	// Execute the query and check for errors
	err := query.First(&queueDisplay).Error
	if err != nil {
		return nil, fmt.Errorf("could not fetch current ticket: %w", err)
	}

	return &queueDisplay, nil
}
