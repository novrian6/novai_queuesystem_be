package models

import (
	"errors"
	"queue-system-backend/database"

	"gorm.io/gorm"
)

// Role struct mapped to database structure
type Role struct {
	RoleID     uint   `gorm:"primaryKey;autoIncrement" json:"role_id"`
	RoleName   string `gorm:"not null" json:"role_name"`
	Permission string `gorm:"not null" json:"permission"`
}

// TableName ensures GORM uses the correct table name
func (Role) TableName() string {
	return "Roles"
}

// GetRoleByID retrieves a role by ID
func GetRoleByID(id uint) (*Role, error) {
	var role Role
	err := database.DB.Where("role_id = ?", id).First(&role).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("role not found")
	}
	return &role, err
}

// CreateRole creates a new role
func CreateRole(role *Role) error {
	if err := database.DB.Create(role).Error; err != nil {
		return errors.New("failed to create role: " + err.Error())
	}
	return nil
}

// UpdateRole updates an existing role
func UpdateRole(role *Role) error {
	if err := database.DB.Save(role).Error; err != nil {
		return errors.New("failed to update role: " + err.Error())
	}
	return nil
}

// DeleteRole deletes a role by ID
func DeleteRole(id uint) error {
	if err := database.DB.Delete(&Role{}, id).Error; err != nil {
		return errors.New("failed to delete role: " + err.Error())
	}
	return nil
}

// ListRoles retrieves all roles
func ListRoles() ([]Role, error) {
	var roles []Role
	if err := database.DB.Find(&roles).Error; err != nil {
		return nil, errors.New("failed to fetch roles: " + err.Error())
	}
	return roles, nil
}
