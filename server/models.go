package main

import (
	"time"
)

type DateTime struct {
	time.Time
}

type Transaction struct {
	TransactionId 		uint		`gorm:"primaryKey"`
	RequestId 			uint
	TerminalId 			uint
	PartnerObjectId 	uint
	AmountTotal 		float32		`gorm:"type:numeric(8,2)"`
	AmountOriginal		float32		`gorm:"type:numeric(8,2)"`
	CommissionPS 		float32		`gorm:"type:numeric(8,2)"`
	CommissionClient 	float32		`gorm:"type:numeric(8,2)"`
	CommissionProvider 	float32		`gorm:"type:numeric(8,2)"`
	DateInput 			DateTime	`gorm:"type:timestamp without time zone"`
	DatePost 			DateTime	`gorm:"type:timestamp without time zone"`
	Status 				string
	PaymentType 		string
	PaymentNumber 		string
	ServiceId 			uint
	Service 			string
	PayeeId 			uint
	PayeeName 			string
	PayeeBankMfo 		uint
	PayeeBankAccount 	string
	PaymentNarrative 	string
}

type SearchTransaction struct {
	TransactionId 		uint		`json:"transaction_id" form:"transaction_id" query:"transaction_id"`
	TerminalId 			[]uint		`json:"terminal_id" form:"terminal_id" query:"terminal_id"`
	Status 				string		`json:"status" form:"status" query:"status"`
	PaymentType 		string		`json:"payment_type" form:"payment_type" query:"payment_type"`
	DatePostFrom		string		`json:"date_post_from" form:"date_post_from" query:"date_post_from"`
	DatePostTo			string		`json:"date_post_to" form:"date_post_to" query:"date_post_to"`
	PaymentNarrative	string		`json:"payment_narrative" form:"payment_narrative" query:"payment_narrative"`
}