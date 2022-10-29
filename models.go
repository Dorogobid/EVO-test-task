package main

type ErrorResponse struct {
	Message string `json:"message" example:"Error message"`
}

type SucsessResponse struct {
	Message string `json:"message" example:"Sucsess message"`
}

type Transaction struct {
	TransactionId 		uint		`gorm:"primaryKey;type:serial" json:"transaction_id" example:"1"`
	RequestId 			uint		`json:"request_id" example:"20020"`
	TerminalId 			uint		`json:"terminal_id" example:"3506"`
	PartnerObjectId 	uint		`json:"partner_object_id" example:"1111"`
	AmountTotal 		float32		`gorm:"type:numeric(8,2)" json:"amount_total" example:"1899"`
	AmountOriginal		float32		`gorm:"type:numeric(8,2)" json:"amount_original" example:"1899"`
	CommissionPS 		float32		`gorm:"type:numeric(8,2)" json:"commission_ps" example:"1.33"`
	CommissionClient 	float32		`gorm:"type:numeric(8,2)" json:"commission_client" example:"0"`
	CommissionProvider 	float32		`gorm:"type:numeric(8,2)" json:"commission_provider" example:"-3.8"`
	DateInput 			string		`gorm:"type:timestamp without time zone" json:"date_input" example:"2022-08-23T09:04:49Z"`
	DatePost 			string		`gorm:"type:timestamp without time zone" json:"date_post" example:"2022-08-23T09:04:50Z"`
	Status 				string		`json:"status" example:"accepted"`
	PaymentType 		string		`json:"payment_type" example:"cash"`
	PaymentNumber 		string		`json:"payment_number" example:"PS16698705"`
	ServiceId 			uint		`json:"service_id" example:"14480"`
	Service 			string		`json:"service" example:"Поповнення карток"`
	PayeeId 			uint		`json:"payee_id" example:"19237155"`
	PayeeName 			string		`json:"payee_name" example:"privat"`
	PayeeBankMfo 		uint		`json:"payee_bank_mfo" example:"304801"`
	PayeeBankAccount 	string		`json:"payee_bank_account" example:"UA713949358919023"`
	PaymentNarrative 	string		`json:"payment_narrative" example:"Перерахування коштів згідно договору про надання послуг А11/27122 від 19.11.2020 р."`
}

type SearchTransaction struct {
	TransactionId 		uint		`query:"transaction_id" json:"transaction_id" example:"0"`
	TerminalId 			[]uint		`query:"terminal_id" json:"terminal_id" example:"3509,3510"`
	Status 				string		`query:"status" json:"status" example:"accepted"`
	PaymentType 		string		`query:"payment_type" json:"payment_type" example:""`
	DatePostFrom		string		`query:"date_post_from" json:"date_post_from" example:"2022-08-17"`
	DatePostTo			string		`query:"date_post_to" json:"date_post_to" example:""`
	PaymentNarrative	string		`query:"payment_narrative" json:"payment_narrative" example:""`
}