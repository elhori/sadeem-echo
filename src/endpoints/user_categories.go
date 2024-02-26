package endpoints

import (
	"errors"
	"net/http"
	"sadeem-echo/src/infra"
	"sadeem-echo/src/models"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func GetAllUserCategories(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("currentPage"))
	if err != nil || page < 1 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid or missing currentPage parameter"})
	}

	pageSize, err := strconv.Atoi(c.QueryParam("pageSize"))
	if err != nil || pageSize < 1 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid or missing pageSize parameter"})
	}

	var userCategories []models.UserCategory
	result := infra.DB().Offset((page - 1) * pageSize).Limit(pageSize).Find(&userCategories)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch user categories"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":        userCategories,
		"currentPage": page,
		"pageSize":    pageSize,
	})
}

func GetUserCategoryByID(c echo.Context) error {
	categoryID, err := strconv.Atoi(c.Param("id"))
	if err != nil || categoryID < 1 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid category ID"})
	}

	var category models.UserCategory
	result := infra.DB().First(&category, categoryID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Category not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch category"})
	}

	return c.JSON(http.StatusOK, category)
}

func CreateUserCategory(c echo.Context) error {
	var category models.UserCategory
	if err := c.Bind(&category); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
	}

	if category.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Category name is required"})
	}

	result := infra.DB().Create(&category)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create category"})
	}

	return c.JSON(http.StatusCreated, category)
}

func UpdateUserCategory(c echo.Context) error {
	categoryID, err := strconv.Atoi(c.Param("id"))
	if err != nil || categoryID < 1 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid category ID"})
	}

	var updatedCategory models.UserCategory
	if err := c.Bind(&updatedCategory); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
	}

	if updatedCategory.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Category name is required"})
	}

	result := infra.DB().Model(&models.UserCategory{}).Where("id = ?", categoryID).Updates(updatedCategory)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update category"})
	}

	return c.JSON(http.StatusOK, updatedCategory)
}

func DeleteUserCategoryByID(c echo.Context) error {
	categoryID, err := strconv.Atoi(c.Param("id"))
	if err != nil || categoryID < 1 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid category ID"})
	}

	result := infra.DB().Delete(&models.UserCategory{}, categoryID)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete category"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Category deleted successfully"})
}

func SearchUserCategories(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("currentPage"))
	if err != nil || page < 1 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid or missing currentPage parameter"})
	}

	pageSize, err := strconv.Atoi(c.QueryParam("pageSize"))
	if err != nil || pageSize < 1 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid or missing pageSize parameter"})
	}

	query := c.QueryParam("q")

	var categories []models.UserCategory
	result := infra.DB().Where("name LIKE ?", "%"+query+"%").Offset((page - 1) * pageSize).Limit(pageSize).Find(&categories)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch categories"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":        categories,
		"currentPage": page,
		"pageSize":    pageSize,
	})
}
