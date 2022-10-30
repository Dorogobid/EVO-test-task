package main

import (
	_ "github.com/Dorogobid/EVO-test-task/docs"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

// @title EVO test application API
// @version 1.0.0
// @host localhost:8080
// @BasePath /

func main() {
	var s ServerInterface
	db := &DBManager{}

	if err := initConfig(); err != nil {
		panic("Error initializing configs")
	}

	db.ConnectToDb(DBConfig{
		Host:     viper.GetString("db.host"),
		Username: viper.GetString("db.username"),
		Password: viper.GetString("db.password"),
		DBName:   viper.GetString("db.dbname"),
		Port:     viper.GetString("db.port"),
		SSLMode:  viper.GetString("db.sslmode"),
	})

	s = &Server{db: db, e: echo.New()}
	s.ConfigureServer()
	s.StartServer()
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
