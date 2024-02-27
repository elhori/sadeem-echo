package endpoints

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sadeem-echo/src/infra"
	"sadeem-echo/src/models"
	"strconv"
	"time"

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

func CreateUser(c echo.Context) error {
	// Extract form data
	name := c.FormValue("name")
	email := c.FormValue("email")
	password := c.FormValue("password")
	role := c.FormValue("role")
	categoryID, err := strconv.Atoi(c.FormValue("category_id"))
	if err != nil {
		return err
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Upload picture
	file, err := c.FormFile("picture")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Failed to upload picture"})
	}

	// Create a directory if it doesn't exist
	err = os.MkdirAll("uploads", os.ModePerm)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create directory"})
	}

	// Generate file path
	ext := filepath.Ext(file.Filename)
	picturePath := "uploads/" + strconv.FormatInt(time.Now().Unix(), 10) + ext

	// Save uploaded file to the server
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to open uploaded file"})
	}
	defer src.Close()

	dst, err := os.Create(picturePath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create picture file"})
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save picture"})
	}

	// Create user object
	user := &models.User{
		Name:       name,
		Email:      email,
		Password:   string(hashedPassword),
		PictureUrl: picturePath,
		CategoryId: categoryID,
		Role:       role,
	}

	// Save user to database
	if err := infra.DB().Create(user).Error; err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, user)
}

func UpdateUser(c echo.Context) error {
	// Extract user ID from route parameter
	requestedUserIDStr := c.Param("id")
	requestedUserID, err := strconv.Atoi(requestedUserIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	// Extract form data
	name := c.FormValue("name")
	email := c.FormValue("email")
	password := c.FormValue("password")
	categoryID, err := strconv.Atoi(c.FormValue("category_id"))
	if err != nil {
		return err
	}

	// Hash password if provided
	if password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		password = string(hashedPassword)
	}

	// Upload picture
	file, err := c.FormFile("picture")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Failed to upload picture"})
	}

	// Create a directory if it doesn't exist
	err = os.MkdirAll("uploads", os.ModePerm)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create directory"})
	}

	// Generate file path
	ext := filepath.Ext(file.Filename)
	picturePath := "uploads/" + strconv.FormatInt(time.Now().Unix(), 10) + ext

	// Save uploaded file to the server
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to open uploaded file"})
	}
	defer src.Close()

	dst, err := os.Create(picturePath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create picture file"})
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save picture"})
	}

	// Find user in database
	var user models.User
	if err := infra.DB().First(&user, requestedUserID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	// Update user data
	user.Name = name
	user.Email = email
	user.Password = password
	user.PictureUrl = picturePath
	user.CategoryId = categoryID

	// Save updated user to database
	if err := infra.DB().Save(&user).Error; err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user)
}

// DeleteUserByID deletes a user by ID (admin only)
func DeleteUserByID(c echo.Context) error {
	// Extract user ID from path parameter
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil || userID < 1 {
		return echo.ErrBadRequest
	}

	// Extract user ID from JWT token
	currentUserID := c.Get("user").(int)

	// Only allow admin to delete users
	role := c.Get("role").(string)
	if role != "Admin" {
		return echo.ErrForbidden
	}

	// Do not allow admin to delete itself
	if int(userID) == currentUserID {
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
