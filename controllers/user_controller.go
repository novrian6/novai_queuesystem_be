package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"queue-system-backend/database"
	"queue-system-backend/models"
	"queue-system-backend/utils"

	"github.com/gin-gonic/gin"
)

// List all users (Admin only)
func ListUsers(c *gin.Context) {
	user, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		return
	}

	claims, ok := user.(*utils.Claims)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid claims"})
		return
	}

	var users []models.User
	var err error

	// Admin can list users with id equal to their ID or owner_id equal to their ID
	if strings.EqualFold(claims.Role, "admin") {
		users, err = models.ListUsersByAdmin(claims.UserID)
	} else {
		// Non-admin can only list themselves
		users, err = models.ListUsersByID(claims.UserID)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	// Hide sensitive data
	for i := range users {
		users[i].PasswordHash = ""
	}

	c.JSON(http.StatusOK, users)
}

// Create new user
func CreateUser(c *gin.Context) {
	userClaims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		return
	}

	claims, ok := userClaims.(*utils.Claims)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized access"})
		return
	}

	var req struct {
		Username    string `json:"username" binding:"required"`
		Password    string `json:"password" binding:"required"`
		Email       string `json:"email" binding:"required"`
		CompanyName string `json:"company_name"`
		RoleID      uint   `json:"role_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	// Ensure the role is not 0 (super admin) or 1 (admin)
	if req.RoleID == 0 || req.RoleID == 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot create user with role super admin or admin"})
		return
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password", "details": err.Error()})
		return
	}

	// Create the user model
	user := models.User{
		Username:     req.Username,
		PasswordHash: hashedPassword,
		Email:        req.Email,
		CompanyName:  req.CompanyName,
		RoleID:       &req.RoleID,
		OwnerID:      &claims.UserID, // Set the owner_id to the requesting user's ID
	}

	// Save the user to the database
	if err := user.CreateUser(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user", "details": err.Error()})
		return
	}

	// Hide sensitive data
	user.PasswordHash = ""
	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "user": user})
}

// GetUser retrieves a specific user by ID
func GetUser(c *gin.Context) {
	requestingUserID, _ := c.Get("user_id")
	userClaims, _ := c.Get("claims")

	reqUserID, _ := requestingUserID.(uint)
	claims, _ := userClaims.(*utils.Claims)

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if uint(id) != reqUserID && strings.ToLower(claims.Role) != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized access"})
		return
	}

	user, err := models.GetUserByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Hide sensitive data
	user.PasswordHash = ""
	c.JSON(http.StatusOK, user)
}

// Update user
func UpdateUser(c *gin.Context) {
	userClaims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		return
	}

	claims, ok := userClaims.(*utils.Claims)
	if !ok || strings.ToLower(claims.Role) != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized access"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req struct {
		Username    string `json:"username"`
		Password    string `json:"password"`
		Email       string `json:"email"`
		CompanyName string `json:"company_name"`
		RoleID      uint   `json:"role_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.First(&user, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Ensure the admin can only update users with owner_id equal to their user_id
	if *user.OwnerID != uint(0) && *user.OwnerID != claims.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized access"})
		return
	}

	// Update fields if provided
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Password != "" {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password", "details": err.Error()})
			return
		}
		user.PasswordHash = hashedPassword
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.CompanyName != "" {
		user.CompanyName = req.CompanyName
	}
	if req.RoleID != 0 {
		user.RoleID = &req.RoleID
	}

	// Save the updated user to the database
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user", "details": err.Error()})
		return
	}

	// Hide sensitive data
	user.PasswordHash = ""
	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully", "user": user})
}

// Delete user (Admin only)
func DeleteUser(c *gin.Context) {
	userClaims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		return
	}

	claims, ok := userClaims.(*utils.Claims)
	if !ok || strings.ToLower(claims.Role) != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized access"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Ensure the admin can only delete users with owner_id equal to their user_id
	if user.OwnerID == nil || *user.OwnerID != claims.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized access"})
		return
	}

	if err := models.DeleteUser(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
