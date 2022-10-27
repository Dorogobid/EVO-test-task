package main

import (
	"fmt"
	"net/http"

	"github.com/gocarina/gocsv"
	"github.com/labstack/echo/v4"
)

// Upload CSV endpoint
func uploadCSV(c echo.Context) error {
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
		return c.String(http.StatusInternalServerError, err.Error())
	}

	if err := loadCSVToDB(transactions); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.HTML(http.StatusOK, fmt.Sprintf("<p>File %s uploaded successfully</p>", file.Filename))
}

func searchToJSON(c echo.Context) error {
	s := new(SearchTransaction)
	if err := c.Bind(s); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	transactions := getFilteredData(s)

	return c.JSON(http.StatusOK, transactions)
}

func searchToCSV(c echo.Context) error {
	s := new(SearchTransaction)
	if err := c.Bind(s); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	transactions := getFilteredData(s)

	csvContent, err := gocsv.MarshalString(&transactions) 
	if err != nil {
		return c.String(http.StatusInternalServerError, "can not marshal data")
	}
	
	return c.Blob(http.StatusOK, "text/csv", []byte(csvContent))
}