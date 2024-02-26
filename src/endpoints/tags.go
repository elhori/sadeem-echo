package endpoints

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"sadeem-echo/src/infra"
	"sadeem-echo/src/models"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func GetAllTags(c echo.Context) error {

	page, err := strconv.Atoi(c.QueryParam("currentPage"))
	if err != nil || page < 1 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid or missing currentPage parameter"})
	}

	pageSize, err := strconv.Atoi(c.QueryParam("pageSize"))
	if err != nil || pageSize < 1 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid or missing pageSize parameter"})
	}

	var tags []models.Tag
	result := infra.DB().Offset((page - 1) * pageSize).Limit(pageSize).Find(&tags)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch tags"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":        tags,
		"currentPage": page,
		"pageSize":    pageSize,
	})
}

func GetTagByID(c echo.Context) error {
	tagID, err := strconv.Atoi(c.Param("id"))
	if err != nil || tagID < 1 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid tag ID"})
	}

	var tag models.Tag
	result := infra.DB().First(&tag, tagID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Tag not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch tag"})
	}

	return c.JSON(http.StatusOK, tag)
}

func CreateTag(c echo.Context) error {
	tag := models.Tag{}
	form, err := c.MultipartForm()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Failed to parse multipart form"})
	}

	name := form.Value["name"][0]
	if name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Tag name is required"})
	}
	tag.Name = name

	pictureFiles := form.File["picture"]
	if len(pictureFiles) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Picture is required"})
	}

	picture := pictureFiles[0]
	file, err := picture.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to open picture"})
	}
	defer file.Close()

	// Create a directory if it doesn't exist
	err = os.MkdirAll("uploads", os.ModePerm)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create directory"})
	}

	filePath := fmt.Sprintf("uploads/%s", picture.Filename)
	dst, err := os.Create(filePath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create picture file"})
	}
	defer dst.Close()

	if _, err = io.Copy(dst, file); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save picture"})
	}

	tag.PictureUrl = filePath
	tag.IsActive = true

	result := infra.DB().Create(&tag)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create tag"})
	}

	return c.JSON(http.StatusCreated, tag)
}

func UpdateTag(c echo.Context) error {
	tagID, err := strconv.Atoi(c.Param("id"))
	if err != nil || tagID < 1 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid tag ID"})
	}

	var tag models.Tag
	result := infra.DB().First(&tag, tagID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Tag not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch tag"})
	}

	form, err := c.MultipartForm()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Failed to parse multipart form"})
	}

	name := form.Value["name"][0]
	if name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Tag name is required"})
	}
	tag.Name = name

	pictureFiles := form.File["picture"]
	if len(pictureFiles) > 0 {
		picture := pictureFiles[0]
		file, err := picture.Open()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to open picture"})
		}
		defer file.Close()

		// Create a directory if it doesn't exist
		err = os.MkdirAll("uploads", os.ModePerm)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create directory"})
		}

		filePath := fmt.Sprintf("uploads/%s", picture.Filename)
		dst, err := os.Create(filePath)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create picture file"})
		}
		defer dst.Close()

		if _, err = io.Copy(dst, file); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save picture"})
		}

		tag.PictureUrl = filePath
	}

	result = infra.DB().Save(&tag)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update tag"})
	}

	return c.JSON(http.StatusOK, tag)
}

func DeleteTagByID(c echo.Context) error {
	tagID, err := strconv.Atoi(c.Param("id"))
	if err != nil || tagID < 1 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid tag ID"})
	}

	result := infra.DB().Delete(&models.Tag{}, tagID)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete tag"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Tag deleted successfully"})
}

func SearchTags(c echo.Context) error {

	page, err := strconv.Atoi(c.QueryParam("currentPage"))
	if err != nil || page < 1 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid or missing currentPage parameter"})
	}

	pageSize, err := strconv.Atoi(c.QueryParam("pageSize"))
	if err != nil || pageSize < 1 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid or missing pageSize parameter"})
	}

	query := c.QueryParam("q")

	var tags []models.Tag
	result := infra.DB().Where("name LIKE ?", "%"+query+"%").Offset((page - 1) * pageSize).Limit(pageSize).Find(&tags)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch tags"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":        tags,
		"currentPage": page,
		"pageSize":    pageSize,
	})
}
