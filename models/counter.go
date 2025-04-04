package models

import (
	"errors"
	"fmt"
	"queue-system-backend/database"

	"gorm.io/gorm"
)

type Counter struct {
	CounterID    uint   `json:"counter_id" gorm:"primaryKey;column:counter_id"`
	VenueID      *uint  `json:"venue_id" gorm:"column:venue_id"`
	ServiceID    *uint  `json:"service_id" gorm:"column:service_id"`
	CounterName  string `json:"counter_name" gorm:"column:counter_name;size:100;not null"`
	OperatorName string `json:"operator_name" gorm:"column:operator_name;size:255"`
	OperatorNIK  string `json:"operator_nik" gorm:"column:operator_nik;size:20"`
	OpenTime     string `json:"open_time" gorm:"column:open_time;type:time"`
	CloseTime    string `json:"close_time" gorm:"column:close_time;type:time"`
	IsVIP        bool   `json:"is_vip" gorm:"column:is_vip;default:0"`
	UserID       uint   `json:"user_id" gorm:"column:user_id;not null"`
}

// TableName ensures GORM uses the correct table name
func (Counter) TableName() string {
	return "Counters"
}

// GetAllCounters retrieves counters based on the user's role
func GetAllCounters(userID uint, isAdmin bool) ([]Counter, error) {
	var counters []Counter

	if isAdmin {
		// Admin can view all counters
		if err := database.DB.Find(&counters).Error; err != nil {
			return nil, err
		}
	} else {
		// Non-admin can only view counters they own
		if err := database.DB.Where("user_id = ?", userID).Find(&counters).Error; err != nil {
			return nil, err
		}
	}

	return counters, nil
}

// GetCounters retrieves counters by user and ID
func GetCounters(userID uint, isAdmin bool) ([]Counter, error) {
	var counters []Counter

	query := database.DB.Model(&Counter{})
	if !isAdmin {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.Find(&counters).Error; err != nil {
		return nil, err
	}

	return counters, nil
}

// CreateCounter inserts a new counter for the given user
func (c *Counter) CreateCounter(userID uint) error {
	// Validate required fields
	if c.CounterName == "" {
		return errors.New("counter name is required")
	}

	// Validate venue if venue_id is provided
	if c.VenueID != nil {
		var venue Venue
		if err := database.DB.First(&venue, *c.VenueID).Error; err != nil {
			return errors.New("venue not found")
		}
	}

	// Validate service if service_id is provided
	if c.ServiceID != nil {
		var service Service
		if err := database.DB.First(&service, *c.ServiceID).Error; err != nil {
			return errors.New("service not found")
		}
	}

	// Set the user ID
	c.UserID = userID

	return database.DB.Create(c).Error
}

// GetCounterByID retrieves a counter by ID
func GetCounterByID(id uint, userID uint, isAdmin bool) (*Counter, error) {
	var counter Counter

	query := database.DB.Where("counter_id = ?", id)
	//if !isAdmin {
	//	query = query.Where("user_id = ?", userID)
	//}

	err := query.First(&counter).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("counter not found")
		}
		return nil, fmt.Errorf("error fetching counter: %v", err)
	}

	return &counter, nil
}

// UpdateCounter updates an existing counter
func (c *Counter) UpdateCounter(userID uint, isAdmin bool) error {
	// Ensure the counter belongs to the user's venue if they are not admin
	if !isAdmin && c.UserID != userID {
		return errors.New("unauthorized: cannot update this counter")
	}
	return database.DB.Save(c).Error
}

// DeleteCounter deletes a counter by ID
func DeleteCounter(id uint, userID uint, isAdmin bool) error {
	query := database.DB.Where("counter_id = ?", id)
	if !isAdmin {
		query = query.Where("user_id = ?", userID)
	}

	return query.Delete(&Counter{}).Error
}

// GetCountersByVenue retrieves all counters belonging to a specific venue
func GetCountersByVenue(venueID uint) ([]Counter, error) {
	var counters []Counter
	query := database.DB.Debug().Where("venue_id = ?", venueID)

	// Print the SQL query to the console
	fmt.Println(query.Statement.SQL.String())

	err := query.Find(&counters).Error
	return counters, err
}
