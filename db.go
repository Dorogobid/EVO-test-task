package main

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func connectToDb() *gorm.DB {
	dsn := "host=postgres user=evo password=evo dbname=evo port=5432 sslmode=disable"
	var db *gorm.DB
	var err error
	for {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			fmt.Println("Connection to database failed. Trying to reconnect...")
		} else {
			fmt.Println("Connected to database done.")
			break
		}
		time.Sleep(time.Second * 5)
	}

	err = db.AutoMigrate(&Transaction{})
	if err != nil {
		panic("Could not run migration.")
	}
	return db
}

func loadCSVToDB(transactions []*Transaction) error {
	tx := DB.Begin()
	defer func() {
	  if r := recover(); r != nil {
		tx.Rollback()
	  }
	}()
  
	if err := tx.Error; err != nil {
	  return err
	}
  
	for _, transaction := range transactions {
		if err := tx.Create(&transaction).Error; err != nil {
			tx.Rollback()
			return err
		 }
	}
	return tx.Commit().Error
}