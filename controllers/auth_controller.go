package controllers

import (
	"errors"
	"net/http"
	"queue-system-backend/models"
	"queue-system-backend/utils"
	"time"

	"github.com/gin-gonic/gin"
)

// Register new user
func Register(c *gin.Context) {
	var req struct {
		Username    string `json:"username" binding:"required"`
		Password    string `json:"password" binding:"required"`
		Email       string `json:"email" binding:"required"`
		CompanyName string `json:"company_name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	// Set RoleID to 0 (admin)
	adminRoleID := uint(0)

	// Create the user model
	user := models.User{
		Username:     req.Username,
		PasswordHash: req.Password,
		Email:        req.Email,
		CompanyName:  req.CompanyName,
		RoleID:       &adminRoleID,
	}

	// Save the user to the database
	if err := user.CreateUser(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// Login and generate JWT token
func Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Attempt to retrieve the user
	user, err := models.GetUserByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Verify password
	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Convert RoleID to string role name
	var roleName string
	if user.RoleID != nil {
		roleName, err = getRoleName(*user.RoleID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		roleName = "No Role Assigned"
	}

	// Generate token
	token, err := utils.GenerateToken(user.UserID, roleName, user.CompanyName, 24*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Successful response
	c.JSON(http.StatusOK, gin.H{
		"token":        token,
		"user_id":      user.UserID,
		"user_name":    user.Username,
		"email":        user.Email,
		"role":         roleName,
		"company_name": user.CompanyName,
		"owner_id":     user.OwnerID,
		"expires":      time.Now().Add(24 * time.Hour).Format(time.RFC3339),
	})
}

// Fetch role name dynamically from database
func getRoleName(roleID uint) (string, error) {
	role, err := models.GetRoleByID(roleID)
	if err != nil {
		return "", errors.New("role not found")
	}
	return role.RoleName, nil
}

// Logout and invalidate session
func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// Get current authenticated user's profile
func Me(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userClaims, ok := claims.(*utils.Claims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid token claims"})
		return
	}

	user, err := models.GetUserByID(userClaims.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Hide sensitive data
	user.PasswordHash = ""
	c.JSON(http.StatusOK, gin.H{"user": user})
}

// Forgot password - Trigger password recovery
func ForgotPassword(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	if err := utils.SendPasswordResetEmail(req.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send reset email", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset email sent"})
}
