package main

import (
	"github.com/Dorogobid/EVO-test-task/docs"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gorm.io/gorm"
)

// @title User API documentation
// @version 1.0.0
// @host localhost:8080
// @BasePath /

var DB *gorm.DB

func main() {
	docs.SwaggerInfo.Title = "EVO test application API"

	DB = connectToDb()
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",}))
	// e.Use(middleware.BodyLimit("90M"))

	e.POST("/upload", uploadCSV)
	e.GET("/search", searchQueryToJSON)
	e.POST("/search", searchJSONToJSON)
	e.GET("/search-csv", searchQueryToCSV)
	e.POST("/search-csv", searchJSONToCSV)

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	
	e.Logger.Fatal(e.Start(":8080"))
}