package main

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBConfig struct {
	Host     string
	Username string
	Password string
	DBName   string
	Port     string
	SSLMode  string
}

type DBManagerInterface interface {
	ConnectToDb(cfg DBConfig)
	LoadCSVToDB(transactions []*Transaction) error
	GetFilteredData(s *SearchTransaction) ([]*Transaction, error)
}

type DBManager struct {
	db *gorm.DB
}

func (db *DBManager) ConnectToDb(cfg DBConfig) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", cfg.Host, cfg.Username, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode)
	var err error
	for {
		db.db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			fmt.Println("Connection to database failed. Trying to reconnect...")
		} else {
			fmt.Println("Connection to database done.")
			break
		}
		time.Sleep(time.Second * 5)
	}

	err = db.db.AutoMigrate(&Transaction{})
	if err != nil {
		panic("Could not run migration.")
	}
}

func (db *DBManager) LoadCSVToDB(transactions []*Transaction) error {
	tx := db.db.Begin()
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

func (db *DBManager) GetFilteredData(s *SearchTransaction) ([]*Transaction, error) {
	transactions := []*Transaction{}
	query := "select * from transactions where"
	var isFirst bool = true

	if s.TransactionId != 0 {
		query += " transaction_id=@transaction_id"
		isFirst = false
	}
	if s.Status != "" {
		if !isFirst {
			query += " and"
		}
		query += " status=@status"
		isFirst = false
	}
	if s.TerminalId != nil && len(s.TerminalId) != 0 {
		if !isFirst {
			query += " and"
		}
		query += " terminal_id in @terminal_id"
		isFirst = false
	}
	if s.PaymentType != "" {
		if !isFirst {
			query += " and"
		}
		query += " payment_type=@payment_type"
		isFirst = false
	}
	if s.DatePostFrom != "" {
		if !isFirst {
			query += " and"
		}
		query += " date_post>=@date_post_from"
		isFirst = false
	}
	if s.DatePostTo != "" {
		if !isFirst {
			query += " and"
		}
		query += " date_post<@date_post_to"
		isFirst = false
	}
	if s.PaymentNarrative != "" {
		if !isFirst {
			query += " and"
		}
		query += " lower(payment_narrative) like lower(@payment_narrative)"
		isFirst = false
	}

	if isFirst {
		return nil, errors.New("QUERY_PARAMS_IS_EMPTY")
	}

	db.db.Raw(query, sql.Named("transaction_id", s.TransactionId), sql.Named("status", s.Status),
		sql.Named("terminal_id", s.TerminalId), sql.Named("payment_type", s.PaymentType),
		sql.Named("date_post_from", s.DatePostFrom), sql.Named("date_post_to", s.DatePostTo),
		sql.Named("payment_narrative", "%"+s.PaymentNarrative+"%")).Find(&transactions)
	return transactions, nil
}
