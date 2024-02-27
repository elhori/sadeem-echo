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

	e.POST("/login", endpoints.Login)
	e.POST("/register", endpoints.Register)

	authGroup := e.Group("")
	authGroup.Use(infra.AuthenticationMiddleware)

	authGroup.GET("/tags", endpoints.GetAllTags)
	authGroup.GET("/tags/:id", endpoints.GetTagByID)
	authGroup.POST("/tags", endpoints.CreateTag)
	authGroup.PUT("/tags/:id", endpoints.UpdateTag)
	authGroup.DELETE("/tags/:id", endpoints.DeleteTagByID)
	authGroup.GET("/tags/search", endpoints.SearchTags)

	authGroup.GET("/user_categories", endpoints.GetAllUserCategories)
	authGroup.GET("/user_categories/:id", endpoints.GetUserCategoryByID)
	authGroup.POST("/user_categories", endpoints.CreateUserCategory)
	authGroup.PUT("/user_categories/:id", endpoints.UpdateUserCategory)
	authGroup.DELETE("/user_categories/:id", endpoints.DeleteUserCategoryByID)
	authGroup.GET("/user_categories/search", endpoints.SearchUserCategories)

	authGroup.GET("/users", endpoints.GetAllUsers)
	authGroup.GET("/users/:id", endpoints.GetUserByID)
	authGroup.POST("/users", endpoints.CreateUser)
	authGroup.PUT("/users/:id", endpoints.UpdateUser)
	authGroup.DELETE("/users/:id", endpoints.DeleteUserByID)
	authGroup.GET("/users/search", endpoints.SearchUsers)

	e.Logger.Fatal(e.Start(":16053"))
}
