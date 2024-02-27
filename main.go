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

	e.GET("/tags", endpoints.GetAllTags, infra.JWTAuth, infra.AuthorizeRoles("Admin", "Default"))
	e.GET("/tags/:id", endpoints.GetTagByID, infra.JWTAuth, infra.AuthorizeRoles("Admin", "Default"))
	e.POST("/tags", endpoints.CreateTag, infra.JWTAuth, infra.AuthorizeRoles("Admin"))
	e.PUT("/tags/:id", endpoints.UpdateTag, infra.JWTAuth, infra.AuthorizeRoles("Admin"))
	e.DELETE("/tags/:id", endpoints.DeleteTagByID, infra.JWTAuth, infra.AuthorizeRoles("Admin"))
	e.GET("/tags/search", endpoints.SearchTags, infra.JWTAuth, infra.AuthorizeRoles("Admin", "Default"))

	e.GET("/user_categories", endpoints.GetAllUserCategories, infra.JWTAuth, infra.AuthorizeRoles("Admin", "Default"))
	e.GET("/user_categories/:id", endpoints.GetUserCategoryByID, infra.JWTAuth, infra.AuthorizeRoles("Admin", "Default"))
	e.POST("/user_categories", endpoints.CreateUserCategory, infra.JWTAuth, infra.AuthorizeRoles("Admin"))
	e.PUT("/user_categories/:id", endpoints.UpdateUserCategory, infra.JWTAuth, infra.AuthorizeRoles("Admin"))
	e.DELETE("/user_categories/:id", endpoints.DeleteUserCategoryByID, infra.JWTAuth, infra.AuthorizeRoles("Admin"))
	e.GET("/user_categories/search", endpoints.SearchUserCategories, infra.JWTAuth, infra.AuthorizeRoles("Admin", "Default"))

	e.Logger.Fatal(e.Start(":8080"))
}
