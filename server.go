package main

import (
	"fmt"
	"net/http"

	_ "github.com/Dorogobid/EVO-test-task/docs"
	"github.com/gocarina/gocsv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type ServerInterface interface {
	ConfigureServer()
	StartServer()
	UploadCSV(c echo.Context) error
	SearchQueryToJSON(c echo.Context) error
	SearchJSONToJSON(c echo.Context) error
	SearchQueryToCSV(c echo.Context) error
	SearchJSONToCSV(c echo.Context) error
}

type Server struct {
	db DBManagerInterface
	e  *echo.Echo
}

func (s *Server) ConfigureServer() {
	s.e.Use(middleware.Recover())
	s.e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n"}))
	s.e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))

	s.e.POST("/upload", s.UploadCSV)

	s.e.GET("/search", s.SearchQueryToJSON)
	s.e.POST("/search", s.SearchJSONToJSON)

	s.e.GET("/search-csv", s.SearchQueryToCSV)
	s.e.POST("/search-csv", s.SearchJSONToCSV)

	s.e.GET("/swagger/*", echoSwagger.WrapHandler)
}

func (s *Server) StartServer() {
	s.e.Logger.Fatal(s.e.Start(viper.GetString("port")))
}

// uploadCSV ... Upload CSV file
// @Summary Import Transactions From File
// @Description Import transactions from CSV file to database
// @Tags Upload
// @Accept mpfd
// @Produce mpfd
// @Param file formData file true "Choose CSV file"
// @Success 200 {object} SucsessResponse
// @Failure 500 {object} ErrorResponse
// @Router /upload [post]
func (s *Server) UploadCSV(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	transactionsFile, err := file.Open()
	if err != nil {
		return err
	}
	defer transactionsFile.Close()

	transactions := []*Transaction{}
	if err := gocsv.Unmarshal(transactionsFile, &transactions); err != nil {
		errResp := ErrorResponse{Message: err.Error()}
		return c.JSON(http.StatusInternalServerError, errResp)
	}

	if err := s.db.LoadCSVToDB(transactions); err != nil {
		errResp := ErrorResponse{Message: err.Error()}
		return c.JSON(http.StatusInternalServerError, errResp)
	}

	sucsResp := SucsessResponse{Message: fmt.Sprintf("File %s uploaded successfully", file.Filename)}
	return c.JSON(http.StatusOK, sucsResp)
}

// searchQueryToJSON ... Get transactions
// @Summary Get transactions
// @Description Search transactions in database (query params) and return in JSON
// @Tags Search
// @Produce json
// @Param transaction_id query uint false "TransactionId"
// @Param terminal_id query []uint false "TerminalId array"
// @Param status query string false "Status"
// @Param payment_type query string false "PaymentType"
// @Param date_post_from query string false "DatePostFrom in format 'YYYY-MM-DD'"
// @Param date_post_to query string false "DatePostTo in format 'YYYY-MM-DD'"
// @Param payment_narrative query string false "PaymentNarrative substring"
// @Success 200 {object} []Transaction
// @Failure 400 {object} ErrorResponse
// @Router /search [get]
func (s *Server) SearchQueryToJSON(c echo.Context) error {
	search := new(SearchTransaction)

	err := s.bindData(c, search)
	if err != nil {
		errResp := ErrorResponse{Message: "Bad request"}
		return c.JSON(http.StatusBadRequest, errResp)
	}

	transactions, err := s.db.GetFilteredData(search)
	if err != nil {
		errResp := ErrorResponse{Message: err.Error()}
		return c.JSON(http.StatusBadRequest, errResp)
	}

	return c.JSON(http.StatusOK, transactions)
}

// searchJSONToJSON ... Get transactions
// @Summary Get transactions
// @Description Search transactions in database (JSON body) and return in JSON
// @Tags Search
// @Accept json
// @Produce json
// @Param request body SearchTransaction true "Request body example"
// @Success 200 {object} []Transaction
// @Failure 400 {object} ErrorResponse
// @Router /search [post]
func (s *Server) SearchJSONToJSON(c echo.Context) error {
	search := new(SearchTransaction)

	err := c.Bind(&search)
	if err != nil {
		errResp := ErrorResponse{Message: "Bad request"}
		return c.JSON(http.StatusBadRequest, errResp)
	}

	transactions, err := s.db.GetFilteredData(search)
	if err != nil {
		errResp := ErrorResponse{Message: err.Error()}
		return c.JSON(http.StatusBadRequest, errResp)
	}

	return c.JSON(http.StatusOK, transactions)
}

// searchQueryToCSV ... Get transactions to CSV file
// @Summary Get transactions to CSV file
// @Description Search transactions in database (query params) and return CSV file
// @Tags Search-csv
// @Produce plain
// @Param transaction_id query uint false "TransactionId"
// @Param terminal_id query []uint false "TerminalId array"
// @Param status query string false "Status"
// @Param payment_type query string false "PaymentType"
// @Param date_post_from query string false "DatePostFrom in format 'YYYY-MM-DD'"
// @Param date_post_to query string false "DatePostTo in format 'YYYY-MM-DD'"
// @Param payment_narrative query string false "PaymentNarrative substring"
// @Success 200 {file} file "CSV file"
// @Failure 400 {object} ErrorResponse
// @Router /search-csv [get]
func (s *Server) SearchQueryToCSV(c echo.Context) error {
	search := new(SearchTransaction)
	err := s.bindData(c, search)
	if err != nil {
		errResp := ErrorResponse{Message: "Bad request"}
		return c.JSON(http.StatusBadRequest, errResp)
	}
	transactions, err := s.db.GetFilteredData(search)
	if err != nil {
		errResp := ErrorResponse{Message: err.Error()}
		return c.JSON(http.StatusBadRequest, errResp)
	}

	csvContent, err := gocsv.MarshalString(&transactions)
	if err != nil {
		errResp := ErrorResponse{Message: err.Error()}
		return c.JSON(http.StatusBadRequest, errResp)
	}

	return c.Blob(http.StatusOK, "text/csv", []byte(csvContent))
}

// searchJSONToCSV ... Get transactions to CSV file
// @Summary Get transactions to CSV file
// @Description Search transactions in database (JSON body) and return CSV file
// @Tags Search-csv
// @Accept json
// @Produce plain
// @Param request body SearchTransaction true "Request body example"
// @Success 200 {file} file "CSV file"
// @Failure 400 {object} ErrorResponse
// @Router /search-csv [post]
func (s *Server) SearchJSONToCSV(c echo.Context) error {
	search := new(SearchTransaction)
	err := c.Bind(&search)
	if err != nil {
		errResp := ErrorResponse{Message: "Bad request"}
		return c.JSON(http.StatusBadRequest, errResp)
	}
	transactions, err := s.db.GetFilteredData(search)
	if err != nil {
		errResp := ErrorResponse{Message: err.Error()}
		return c.JSON(http.StatusBadRequest, errResp)
	}

	csvContent, err := gocsv.MarshalString(&transactions)
	if err != nil {
		errResp := ErrorResponse{Message: err.Error()}
		return c.JSON(http.StatusBadRequest, errResp)
	}

	return c.Blob(http.StatusOK, "text/csv", []byte(csvContent))
}

func (s *Server) bindData(c echo.Context, search *SearchTransaction) error {
	return echo.QueryParamsBinder(c).
		Uint("transaction_id", &search.TransactionId).
		BindWithDelimiter("terminal_id", &search.TerminalId, ",").
		BindWithDelimiter("terminal_id[]", &search.TerminalId, ",").
		String("status", &search.Status).
		String("payment_type", &search.PaymentType).
		String("date_post_from", &search.DatePostFrom).
		String("date_post_to", &search.DatePostTo).
		String("payment_narrative", &search.PaymentNarrative).
		BindError()
}
