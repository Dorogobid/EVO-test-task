package main

import (
	_ "github.com/Dorogobid/EVO-test-task/docs"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title EVO test application API
// @version 1.0.0
// @host localhost:8080
// @BasePath /
var DB *DBManager

func main() {
	DB = &DBManager{}

	if err := initConfig(); err != nil {
		panic("Error initializing configs")
	}

	DB.connectToDb(DBConfig{
		Host: viper.GetString("db.host"),
		Username: viper.GetString("db.username"),
		Password: viper.GetString("db.password"),
		DBName: viper.GetString("db.dbname"),
		Port: viper.GetString("db.port"),
		SSLMode: viper.GetString("db.sslmode"),
	})

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",}))
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))

	e.POST("/upload", uploadCSV)

	e.GET("/search", searchQueryToJSON)
	e.POST("/search", searchJSONToJSON)

	e.GET("/search-csv", searchQueryToCSV)
	e.POST("/search-csv", searchJSONToCSV)

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	
	e.Logger.Fatal(e.Start(viper.GetString("port")))
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}