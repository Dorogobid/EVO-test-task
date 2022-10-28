package main

import (
	"database/sql"
	"errors"

	"github.com/labstack/echo/v4"
)

func getFilteredData(s *SearchTransaction) ([]*Transaction, error) {
	transactions := []*Transaction{}
	query := "select * from transactions where"
	var isFirst bool = true

	if s.TransactionId != 0 {
		query += " transaction_id=@transaction_id"
		isFirst = false
	}
	if s.Status != "" {
		if !isFirst {query += " and"}
		query += " status=@status"
		isFirst = false
	}
	if s.TerminalId != nil && len(s.TerminalId) != 0{
		if !isFirst {query += " and"}
		query += " terminal_id in @terminal_id"
		isFirst = false
	}
	if s.PaymentType != "" {
		if !isFirst {query += " and"}
		query += " payment_type=@payment_type"
		isFirst = false
	}
	if s.DatePostFrom != "" {
		if !isFirst {query += " and"}
		query += " date_post>=@date_post_from"
		isFirst = false
	} 
	if s.DatePostTo != "" {
		if !isFirst {query += " and"}
		query += " date_post<@date_post_to"
		isFirst = false
	} 
	if s.PaymentNarrative != "" {
		if !isFirst {query += " and"}
		query += " lower(payment_narrative) like lower(@payment_narrative)"
		isFirst = false
	}

	if isFirst {
		return nil, errors.New("QUERY_PARAMS_IS_EMPTY")
	}
	
	DB.Raw(query, sql.Named("transaction_id", s.TransactionId), sql.Named("status", s.Status), 
		sql.Named("terminal_id", s.TerminalId), sql.Named("payment_type", s.PaymentType), 
		sql.Named("date_post_from", s.DatePostFrom), sql.Named("date_post_to", s.DatePostTo), 
		sql.Named("payment_narrative", "%"+s.PaymentNarrative+"%")).Find(&transactions)
	return transactions, nil
}

func bindData(c echo.Context, s *SearchTransaction) error {
	return echo.QueryParamsBinder(c).
	Uint("transaction_id", &s.TransactionId).
	BindWithDelimiter("terminal_id", &s.TerminalId, ",").
	String("status", &s.Status).
	String("payment_type", &s.PaymentType).
	String("date_post_from", &s.DatePostFrom).
	String("date_post_to", &s.DatePostTo).
	String("payment_narrative", &s.PaymentNarrative).
	BindError()
}