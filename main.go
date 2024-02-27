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

	e.GET("/tags", endpoints.GetAllTags /*infra.JWTAuth*/)
	e.GET("/tags/:id", endpoints.GetTagByID /*infra.JWTAuth*/)
	e.POST("/tags", endpoints.CreateTag /*infra.JWTAuth*/)
	e.PUT("/tags/:id", endpoints.UpdateTag /*infra.JWTAuth*/)
	e.DELETE("/tags/:id", endpoints.DeleteTagByID /*infra.JWTAuth*/)
	e.GET("/tags/search", endpoints.SearchTags /*infra.JWTAuth*/)

	e.GET("/user_categories", endpoints.GetAllUserCategories /*infra.JWTAuth*/)
	e.GET("/user_categories/:id", endpoints.GetUserCategoryByID /*infra.JWTAuth*/)
	e.POST("/user_categories", endpoints.CreateUserCategory /*infra.JWTAuth*/)
	e.PUT("/user_categories/:id", endpoints.UpdateUserCategory /*infra.JWTAuth*/)
	e.DELETE("/user_categories/:id", endpoints.DeleteUserCategoryByID /*infra.JWTAuth*/)
	e.GET("/user_categories/search", endpoints.SearchUserCategories /*infra.JWTAuth*/)

	e.GET("/users", endpoints.GetAllUsers /*infra.JWTAuth*/)
	e.GET("/users/:id", endpoints.GetUserByID /*infra.JWTAuth*/)
	e.POST("/users", endpoints.CreateUser /*infra.JWTAuth*/)
	e.PUT("/users/:id", endpoints.UpdateUser /*infra.JWTAuth*/)
	e.DELETE("/users/:id", endpoints.DeleteUserByID /*infra.JWTAuth*/)
	e.GET("/users/search", endpoints.SearchUsers /*infra.JWTAuth*/)

	e.Logger.Fatal(e.Start(":16053"))
}
