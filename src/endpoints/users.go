package endpoints

import (
	"net/http"
	"strconv"

	"sadeem-echo/src/infra"
	"sadeem-echo/src/models"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func GetUserByID(c echo.Context) error {
	// Extract user ID from query parameters
	requestedUserIDStr := c.QueryParam("id")
	requestedUserID, err := strconv.Atoi(requestedUserIDStr)
	if err != nil || requestedUserID < 1 {
		// If no ID provided in query or invalid, use ID from JWT token
		requestedUserID = c.Get("user").(int)
	}

	// Extract current user's ID and role from JWT token
	currentUserID := c.Get("user").(int)
	role := c.Get("role").(string)

	// Check if the requester is an admin or if the requested user ID matches the current user ID
	if role != "Admin" && requestedUserID != currentUserID {
		return echo.ErrForbidden
	}

	// Find user by ID
	var user models.User
	if err := infra.DB().First(&user, requestedUserID).Error; err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user)
}

// CreateUser creates a new user (admin only)
func CreateUser(c echo.Context) error {
	// Parse request body
	req := new(models.User)
	if err := c.Bind(req); err != nil {
		return err
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	req.Password = string(hashedPassword)

	// Save user to database
	if err := infra.DB().Create(req).Error; err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, req)
}

// UpdateUser updates a user (admin or self)
func UpdateUser(c echo.Context) error {
	// Extract user ID from JWT token
	userID := c.Get("user").(int)

	// Parse request body
	req := new(models.User)
	if err := c.Bind(req); err != nil {
		return err
	}

	// Only allow admin to update other users
	if req.Id != userID {
		role := c.Get("role").(string)
		if role != "Admin" {
			return echo.ErrForbidden
		}
	}

	// Hash password if provided
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		req.Password = string(hashedPassword)
	}

	// Update user in database
	if err := infra.DB().Save(req).Error; err != nil {
		return err
	}

	return c.JSON(http.StatusOK, req)
}

// DeleteUserByID deletes a user by ID (admin only)
func DeleteUserByID(c echo.Context) error {
	// Extract user ID from path parameter
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil || userID < 1 {
		return echo.ErrBadRequest
	}

	// Extract user ID from JWT token
	currentUserID := c.Get("user").(uint)

	// Only allow admin to delete users
	role := c.Get("role").(string)
	if role != "Admin" {
		return echo.ErrForbidden
	}

	// Do not allow admin to delete itself
	if uint(userID) == currentUserID {
		return echo.ErrForbidden
	}

	// Delete user from database
	if err := infra.DB().Delete(&models.User{}, userID).Error; err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "User deleted successfully"})
}

// SearchUsers searches users based on query parameters (admin only)
func SearchUsers(c echo.Context) error {
	// Only allow admin to search users
	role := c.Get("role").(string)
	if role != "Admin" {
		return echo.ErrForbidden
	}

	// Parse query parameters
	page, err := strconv.Atoi(c.QueryParam("currentPage"))
	if err != nil || page < 1 {
		return echo.ErrBadRequest
	}
	pageSize, err := strconv.Atoi(c.QueryParam("pageSize"))
	if err != nil || pageSize < 1 {
		return echo.ErrBadRequest
	}
	query := c.QueryParam("q")

	// Search users
	var users []models.User
	result := infra.DB().Where("name LIKE ?", "%"+query+"%").Offset((page - 1) * pageSize).Limit(pageSize).Find(&users)
	if result.Error != nil {
		return result.Error
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":        users,
		"currentPage": page,
		"pageSize":    pageSize,
	})
}

func GetAllUsers(c echo.Context) error {
	// Only allow admin to access all users
	role := c.Get("role").(string)
	if role != "Admin" {
		return echo.ErrForbidden
	}

	// Parse query parameters
	page, err := strconv.Atoi(c.QueryParam("currentPage"))
	if err != nil || page < 1 {
		return echo.ErrBadRequest
	}
	pageSize, err := strconv.Atoi(c.QueryParam("pageSize"))
	if err != nil || pageSize < 1 {
		return echo.ErrBadRequest
	}

	// Retrieve all users
	var users []models.User
	result := infra.DB().Offset((page - 1) * pageSize).Limit(pageSize).Find(&users)
	if result.Error != nil {
		return result.Error
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":        users,
		"currentPage": page,
		"pageSize":    pageSize,
	})
}
