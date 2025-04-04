package models

import (
	"errors"
	"queue-system-backend/database"

	"gorm.io/gorm"
)

type Venue struct {
	VenueID    uint   `json:"venue_id" gorm:"primaryKey;autoIncrement"`
	UserID     uint   `json:"user_id" gorm:"not null"` // Foreign key to Users table
	VenueName  string `json:"venue_name" gorm:"size:255;not null"`
	Address    string `json:"address" gorm:"size:255;not null"`
	City       string `json:"city" gorm:"size:100;not null"`
	Province   string `json:"province" gorm:"size:100;not null"`
	PostalCode string `json:"postal_code" gorm:"size:10;not null"`
	Phone      string `json:"phone" gorm:"size:20"`
	Email      string `json:"email" gorm:"size:100"`
	OpenTime   string `json:"open_time" gorm:"type:time"`
	CloseTime  string `json:"close_time" gorm:"type:time"`
}

// TableName ensures GORM uses the correct table name
func (Venue) TableName() string {
	return "Venues"
}

// GetVenueByID retrieves a venue by ID from the database
func GetVenueByID(id uint) (*Venue, error) {
	var venue Venue
	err := database.DB.First(&venue, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("venue not found")
		}
		return nil, errors.New("database error: " + err.Error())
	}
	return &venue, nil
}

// CreateVenue adds a new venue to the database
func (v *Venue) CreateVenue() error {
	if database.DB == nil {
		return errors.New("database connection is not initialized")
	}

	if v.UserID == 0 || v.VenueName == "" {
		return errors.New("user_id and venue_name are required")
	}

	if err := database.DB.Create(v).Error; err != nil {
		return err
	}
	return nil
}

// UpdateVenue updates an existing venue in the database
func (v *Venue) UpdateVenue() error {
	if database.DB == nil {
		return errors.New("database connection is not initialized")
	}

	if err := database.DB.Save(v).Error; err != nil {
		return err
	}
	return nil
}

// DeleteVenue deletes a venue from the database
func DeleteVenue(id uint) error {
	if database.DB == nil {
		return errors.New("database connection is not initialized")
	}

	if err := database.DB.Delete(&Venue{}, id).Error; err != nil {
		return err
	}
	return nil
}

// GetAllVenues retrieves all venues (Admin only)
func GetAllVenues() ([]Venue, error) {
	var venues []Venue
	if err := database.DB.Find(&venues).Error; err != nil {
		return nil, errors.New("failed to fetch venues: " + err.Error())
	}
	return venues, nil
}

// GetVenuesByUser retrieves venues filtered by user_id
func GetVenuesByUser(userID uint) ([]Venue, error) {
	var venues []Venue
	if err := database.DB.Where("user_id = ?", userID).Find(&venues).Error; err != nil {
		return nil, errors.New("failed to fetch venues: " + err.Error())
	}
	return venues, nil
}

func GetVenueNameByID(venueID uint) (string, error) {
	var venue Venue
	if err := database.DB.Select("venue_name").First(&venue, venueID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("venue not found")
		}
		return "", err
	}
	return venue.VenueName, nil
}
