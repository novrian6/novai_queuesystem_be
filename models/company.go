package models

import (
	"errors"
	"queue-system-backend/database"

	"gorm.io/gorm"
)

// Company struct mapped to database structure
type Company struct {
	CompanyID    uint   `gorm:"primaryKey;autoIncrement" json:"company_id"`
	CompanyName  string `gorm:"not null" json:"company_name"`
	ContactEmail string `json:"contact_email"`
	ContactPhone string `json:"contact_phone"`
	CreatedAt    string `json:"created_at"`
}

// TableName ensures GORM uses the correct table name
func (Company) TableName() string {
	return "Companies"
}

// GetCompanyByID retrieves a company by ID
func GetCompanyByID(id uint) (*Company, error) {
	var company Company
	err := database.DB.Where("company_id = ?", id).First(&company).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("company not found")
		}
		return nil, errors.New("database error: " + err.Error())
	}
	return &company, nil
}

// CreateCompany creates a new company
func CreateCompany(company *Company) error {
	if err := database.DB.Create(company).Error; err != nil {
		return errors.New("failed to create company: " + err.Error())
	}
	return nil
}

// UpdateCompany updates an existing company
func UpdateCompany(company *Company) error {
	if err := database.DB.Save(company).Error; err != nil {
		return errors.New("failed to update company: " + err.Error())
	}
	return nil
}

// DeleteCompany deletes a company by ID
func DeleteCompany(id uint) error {
	if err := database.DB.Delete(&Company{}, id).Error; err != nil {
		return errors.New("failed to delete company: " + err.Error())
	}
	return nil
}
