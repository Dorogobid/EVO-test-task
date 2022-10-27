package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

var DB *gorm.DB

func main() {
	DB = connectToDb()
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",}))
	e.Use(middleware.BodyLimit("90M"))

	e.POST("/upload", uploadCSV)
	e.GET("/search", searchToJSON)
	e.GET("/search-csv", searchToCSV)
	e.Logger.Fatal(e.Start(":8080"))
}