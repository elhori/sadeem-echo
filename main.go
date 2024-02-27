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
	authGroup.POST("/tags", endpoints.CreateTag, infra.AuthorizationMiddleware("Admin"))
	authGroup.PUT("/tags/:id", endpoints.UpdateTag, infra.AuthorizationMiddleware("Admin"))
	authGroup.DELETE("/tags/:id", endpoints.DeleteTagByID, infra.AuthorizationMiddleware("Admin"))
	authGroup.GET("/tags/search", endpoints.SearchTags)

	authGroup.GET("/user_categories", endpoints.GetAllUserCategories)
	authGroup.GET("/user_categories/:id", endpoints.GetUserCategoryByID)
	authGroup.POST("/user_categories", endpoints.CreateUserCategory, infra.AuthorizationMiddleware("Admin"))
	authGroup.PUT("/user_categories/:id", endpoints.UpdateUserCategory, infra.AuthorizationMiddleware("Admin"))
	authGroup.DELETE("/user_categories/:id", endpoints.DeleteUserCategoryByID, infra.AuthorizationMiddleware("Admin"))
	authGroup.GET("/user_categories/search", endpoints.SearchUserCategories)

	authGroup.GET("/users", endpoints.GetAllUsers, infra.AuthorizationMiddleware("Admin"))
	authGroup.GET("/users/get", endpoints.GetUserByID)
	authGroup.POST("/users", endpoints.CreateUser, infra.AuthorizationMiddleware("Admin"))
	authGroup.PUT("/users/:id", endpoints.UpdateUser, infra.AuthorizationMiddleware("Admin"))
	authGroup.DELETE("/users/:id", endpoints.DeleteUserByID, infra.AuthorizationMiddleware("Admin"))
	authGroup.GET("/users/search", endpoints.SearchUsers, infra.AuthorizationMiddleware("Admin"))

	e.Logger.Fatal(e.Start(":16053"))
}
