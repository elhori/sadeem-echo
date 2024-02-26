package main

import (
	"sadeem-echo/src/endpoints"
	"sadeem-echo/src/infra"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	infra.DatabaseInit()
	gorm := infra.DB()

	dbGorm, err := gorm.DB()
	if err != nil {
		panic(err)
	}

	dbGorm.Ping()

	e.GET("/tags", endpoints.GetAllTags)
	e.GET("/tags/:id", endpoints.GetTagByID)
	e.POST("/tags", endpoints.CreateTag)
	e.PUT("/tags/:id", endpoints.UpdateTag)
	e.DELETE("/tags/:id", endpoints.DeleteTagByID)
	e.GET("/tags/search", endpoints.SearchTags)

	e.GET("/user_categories", endpoints.GetAllUserCategories)
	e.GET("/user_categories/:id", endpoints.GetUserCategoryByID)
	e.POST("/user_categories", endpoints.CreateUserCategory)
	e.PUT("/user_categories/:id", endpoints.UpdateUserCategory)
	e.DELETE("/user_categories/:id", endpoints.DeleteUserCategoryByID)
	e.GET("/user_categories/search", endpoints.SearchUserCategories)

	e.Logger.Fatal(e.Start(":8080"))
}
