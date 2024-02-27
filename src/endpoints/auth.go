package endpoints

import (
	"net/http"

	"sadeem-echo/src/infra"
	"sadeem-echo/src/models"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func Login(c echo.Context) error {
	// Parse request body
	req := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	// Find user by email
	var user models.User
	if err := infra.DB().Where("email = ?", req.Email).First(&user).Error; err != nil {
		return echo.ErrUnauthorized
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return echo.ErrUnauthorized
	}

	// Generate JWT token
	token, err := infra.GenerateToken(user.Id)
	if err != nil {
		return err
	}

	// Set token in user object
	user.Token = token

	c.Set("role", user.Role)
	c.Set("token", user.Token)
	//c.Set("user", user)

	return c.JSON(http.StatusOK, user)
}

func Register(c echo.Context) error {
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

	// Set role to Default
	req.Role = "Default"

	// Save user to database
	if err := infra.DB().Create(req).Error; err != nil {
		return err
	}

	// Generate JWT token
	token, err := infra.GenerateToken(req.Id)
	if err != nil {
		return err
	}

	// Set token in user object
	req.Token = token

	return c.JSON(http.StatusCreated, req)
}
