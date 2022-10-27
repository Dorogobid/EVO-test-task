package main

import "database/sql"

func getFilteredData(s *SearchTransaction) []*Transaction {
	transactions := []*Transaction{}
	query := "select * from transactions where"
	isFirst := true

	if s.TransactionId != 0 {
		query += " transaction_id=@transaction_id"
		isFirst = false
	}
	if s.Status != "" {
		if !isFirst {query += " and"}
		query += " status=@status"
		isFirst = false
	}
	if s.TerminalId != nil {
		if !isFirst {query += " and"}
		query += " terminal_id in @terminal_id"
		isFirst = false
	}
	if s.PaymentType != "" {
		if !isFirst {query += " and"}
		query += " payment_type=@payment_type"
		isFirst = false
	}
	if s.DatePostFrom != "" && s.DatePostTo != "" {
		if !isFirst {query += " and"}
		query += " date_post>=@date_post_from and date_post<@date_post_to"
		isFirst = false
	} 
	if s.PaymentNarrative != "" {
		if !isFirst {query += " and"}
		query += " lower(payment_narrative) like lower(@payment_narrative)"
		isFirst = false
	}
	
	DB.Raw(query, sql.Named("transaction_id", s.TransactionId), sql.Named("status", s.Status), 
		sql.Named("terminal_id", s.TerminalId), sql.Named("payment_type", s.PaymentType), 
		sql.Named("date_post_from", s.DatePostFrom), sql.Named("date_post_to", s.DatePostTo), 
		sql.Named("payment_narrative", "%"+s.PaymentNarrative+"%")).Find(&transactions)
	return transactions
}