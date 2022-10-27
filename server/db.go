package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func connectToDb() *gorm.DB {
	dsn := "host=localhost user=evo password=evo dbname=evo port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Connection to database failed.")
	}

	db.AutoMigrate(&Transaction{})
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