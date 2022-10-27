package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func main() {
	DB = connectToDb()
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	  }))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("/upload", uploadCSV)
	e.GET("/search", searchToJSON)
	e.GET("/search-csv", searchToCSV)
	e.Logger.Fatal(e.Start(":8080"))
}

func connectToDb() *gorm.DB {
	dsn := "host=localhost user=evo password=evo dbname=evo port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Connection to database failed.")
	}

	db.AutoMigrate(&Transaction{})
	return db
}