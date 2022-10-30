package main

import (
	"fmt"
	"os"

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

func main() {
	if err := initConfig(); err != nil {
		panic("Error initializing configs")
	}

	db := &DBManager{}
	db.ConnectToDb(getConfig())

	var h HandlerInterface = &Handler{db: db}

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n"}))
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))

	e.POST("/upload", h.UploadCSV)

	e.GET("/search", h.SearchQueryToJSON)
	e.POST("/search", h.SearchJSONToJSON)

	e.GET("/search-csv", h.SearchQueryToCSV)
	e.POST("/search-csv", h.SearchJSONToCSV)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Logger.Fatal(e.Start(viper.GetString("port")))
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func getConfig() *DBConfig {
	host, isHostPresent := os.LookupEnv("DB_HOST")
	if !isHostPresent {
		host = viper.GetString("db.host")
	}
	fmt.Println(host)
	return &DBConfig{
		Host:     host,
		Username: viper.GetString("db.username"),
		Password: viper.GetString("db.password"),
		DBName:   viper.GetString("db.dbname"),
		Port:     viper.GetString("db.port"),
		SSLMode:  viper.GetString("db.sslmode"),
	}
}