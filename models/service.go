package models

import (
	"errors"
	"queue-system-backend/database"

	"gorm.io/gorm"
)

type Service struct {
	ServiceID   uint   `json:"service_id" gorm:"primaryKey;autoIncrement"`
	UserID      *uint  `json:"user_id" gorm:"default:null"`  // Foreign key to Users table
	VenueID     *uint  `json:"venue_id" gorm:"default:null"` // Foreign key to Venues table
	ServiceName string `json:"service_name" gorm:"size:255;not null"`
	Description string `json:"description" gorm:"size:100;default:null"`
}

// TableName ensures GORM uses the correct table name
func (Service) TableName() string {
	return "Services"
}

// GetServicesByUser retrieves services filtered by user_id
func GetServicesByUser(userID uint) ([]Service, error) {
	var services []Service
	if err := database.DB.Where("user_id = ?", userID).Find(&services).Error; err != nil {
		return nil, err
	}
	return services, nil
}

// GetServicesByVenue retrieves services filtered by venue_id
func GetServicesByVenue(venueID uint) ([]Service, error) {
	var services []Service
	if err := database.DB.Where("venue_id = ?", venueID).Find(&services).Error; err != nil {
		return nil, err
	}
	return services, nil
}

// CreateService adds a new service to the database
func CreateService(service *Service) error {
	if service.ServiceName == "" {
		return errors.New("service_name is required")
	}
	// No validation for description
	return database.DB.Create(service).Error
}

// GetServiceByID retrieves a specific service by ID
func GetServiceByID(serviceID uint) (*Service, error) {
	var service Service
	if err := database.DB.First(&service, serviceID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("service not found")
		}
		return nil, err
	}
	return &service, nil
}

// UpdateService updates an existing service in the database
func UpdateService(service *Service) error {
	if service.ServiceName == "" {
		return errors.New("service_name is required")
	}
	// No validation for description
	return database.DB.Save(service).Error
}

// DeleteService deletes a service from the database
func DeleteService(serviceID uint) error {
	return database.DB.Delete(&Service{}, serviceID).Error
}

// GetAllServices retrieves all services from the database
func GetAllServices() ([]Service, error) {
	var services []Service
	if err := database.DB.Find(&services).Error; err != nil {
		return nil, err
	}
	return services, nil
}

func GetServicesByVenueAndUser(venueID string, userID uint) ([]Service, error) {
	var services []Service
	if err := database.DB.Where("venue_id = ? AND user_id = ?", venueID, userID).Find(&services).Error; err != nil {
		return nil, err
	}
	return services, nil
}

func GetServiceNameByID(serviceID uint) (string, error) {
	var service Service
	if err := database.DB.Select("service_name").First(&service, serviceID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("service not found")
		}
		return "", err
	}
	return service.ServiceName, nil
}
