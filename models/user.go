package models

import (
	"errors"
	"queue-system-backend/database"
	"queue-system-backend/utils"

	"gorm.io/gorm"
)

// User struct mapped to database structure
type User struct {
	UserID       uint   `gorm:"primaryKey;autoIncrement" json:"user_id"`
	CompanyName  string `gorm:"size:100" json:"company_name"`
	RoleID       *uint  `gorm:"default:null" json:"role_id"`
	Username     string `gorm:"unique;size:50;not null" json:"username"`
	PasswordHash string `gorm:"size:255;not null" json:"-"`
	Email        string `gorm:"unique;size:100;not null" json:"email"`
	OwnerID      *uint  `json:"owner_id"` // New attribute, optional
}

// TableName ensures GORM uses the correct table name "Users"
func (User) TableName() string {
	return "Users"
}

// LoginRequest for authentication payload
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// CreateUser creates a new user in the database
func (user *User) CreateUser() error {
	// Hash password before saving
	hashedPassword, err := utils.HashPassword(user.PasswordHash)
	if err != nil {
		return errors.New("failed to hash password: " + err.Error())
	}

	user.PasswordHash = hashedPassword

	// Create user in the database
	if err := database.DB.Create(user).Error; err != nil {
		return errors.New("failed to create user: " + err.Error())
	}
	return nil
}

// GetUserByUsername retrieves a user by username
func GetUserByUsername(username string) (*User, error) {
	var user User
	err := database.DB.Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}
	return &user, err
}

// GetUserByID retrieves a user by ID
func GetUserByID(id uint) (*User, error) {
	var user User
	err := database.DB.First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}
	return &user, err
}

// UpdateUser updates an existing user's details
func (user *User) UpdateUser(data map[string]interface{}) error {
	if err := database.DB.Model(user).Updates(data).Error; err != nil {
		return errors.New("failed to update user: " + err.Error())
	}
	return nil
}

// DeleteUser deletes a user by ID
func DeleteUser(id uint) error {
	if err := database.DB.Delete(&User{}, id).Error; err != nil {
		return errors.New("failed to delete user: " + err.Error())
	}
	return nil
}

// ListUsers retrieves all users (with optional filters)
func ListUsers() ([]User, error) {
	var users []User
	if err := database.DB.Find(&users).Error; err != nil {
		return nil, errors.New("failed to fetch users: " + err.Error())
	}
	return users, nil
}

// ListUsersByOwnerID retrieves users with owner_id equal to the given user ID
func ListUsersByOwnerID(ownerID uint) ([]User, error) {
	var users []User
	if err := database.DB.Where("owner_id = ?", ownerID).Find(&users).Error; err != nil {
		return nil, errors.New("failed to fetch users: " + err.Error())
	}
	return users, nil
}

// ListUsersByID retrieves a user by their ID
func ListUsersByID(userID uint) ([]User, error) {
	var user User
	if err := database.DB.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("failed to fetch user: " + err.Error())
	}
	return []User{user}, nil
}

// ListUsersByAdmin retrieves users with id equal to the admin's user_id or owner_id equal to the admin's user_id
func ListUsersByAdmin(adminID uint) ([]User, error) {
	var users []User
	if err := database.DB.Where("user_id = ? OR owner_id = ?", adminID, adminID).Find(&users).Error; err != nil {
		return nil, errors.New("failed to fetch users: " + err.Error())
	}
	return users, nil
}
