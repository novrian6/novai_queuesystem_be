package models

import (
	"errors"
	"queue-system-backend/database"
	"time"
)

type UserCounterMap struct {
	UserCounterMapID int       `gorm:"primaryKey;autoIncrement" json:"user_counter_map_id"` // Primary key
	UserID           int       `gorm:"not null" json:"user_id"`                             // Foreign key to Users table
	CounterID        int       `gorm:"not null" json:"counter_id"`                          // Foreign key to Counters table
	AssignedAt       time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"assigned_at"`        // Timestamp for assignment
	OwnerID          uint      `gorm:"not null" json:"owner_id"`
}

// TableName ensures GORM uses the correct table name
func (UserCounterMap) TableName() string {
	return "User_Counter_Map"
}

// ValidateOwnership checks if the counter belongs to the user
func (u *UserCounterMap) ValidateOwnership() error {
	var counter UserCounterMap
	if err := database.DB.First(&counter, u.CounterID).Error; err != nil {
		return errors.New("invalid counter_id")
	}
	if counter.OwnerID != u.OwnerID {
		return errors.New("user does not own this counter")
	}
	return nil
}

// FetchByID retrieves a UserCounterMap by its ID
func FetchByID(id uint) (*UserCounterMap, error) {
	var mapping UserCounterMap
	if err := database.DB.First(&mapping, id).Error; err != nil {
		return nil, err
	}
	return &mapping, nil
}

// FetchAll retrieves all UserCounterMap records
func FetchAll() ([]UserCounterMap, error) {
	var mappings []UserCounterMap
	if err := database.DB.Find(&mappings).Error; err != nil {
		return nil, err
	}
	return mappings, nil
}

// FetchByOwnerID retrieves all UserCounterMap records for a specific owner
func FetchByOwnerID(ownerID uint) ([]UserCounterMap, error) {
	var mappings []UserCounterMap
	if err := database.DB.Where("owner_id = ?", ownerID).Find(&mappings).Error; err != nil {
		return nil, err
	}
	return mappings, nil
}

// Create inserts a new UserCounterMap record into the database
func (u *UserCounterMap) Create() error {
	if err := database.DB.Create(u).Error; err != nil {
		return err
	}
	return nil
}

// Update saves changes to an existing UserCounterMap record
func (u *UserCounterMap) Update() error {
	if err := database.DB.Save(u).Error; err != nil {
		return err
	}
	return nil
}

// Delete removes a UserCounterMap record from the database
func (u *UserCounterMap) Delete() error {
	if err := database.DB.Delete(u).Error; err != nil {
		return err
	}
	return nil
}

// GetUserIDByCounterID retrieves the user_id associated with a given counter_id
func GetUserIDByCounterID(counterID int) (int, error) {
	var mapping UserCounterMap
	if err := database.DB.Where("counter_id = ?", counterID).First(&mapping).Error; err != nil {
		return 0, err
	}
	return mapping.UserID, nil
}

// GetCounterIDByUserID retrieves the counter_id associated with a given user_id
func GetCounterIDByUserID(userID int) (int, error) {
	var mapping UserCounterMap
	if err := database.DB.Where("user_id = ?", userID).First(&mapping).Error; err != nil {
		return 0, err
	}
	return mapping.CounterID, nil
}
