package main

import (
	"database/sql/driver"
	"fmt"
	"time"
)

func (date *DateTime) MarshalCSV() (string, error) {
	return date.Time.Format("2006-01-02 15:04:05"), nil
}

// Convert the CSV string as internal date
func (date *DateTime) UnmarshalCSV(csv string) (err error) {
	date.Time, err = time.Parse("2006-01-02 15:04:05", csv)
	return err
}

func (date *DateTime) Scan(b interface{}) (err error) {
    switch x := b.(type) {
    case time.Time:
        date.Time = x
	default:
    	err = fmt.Errorf("unsupported scan type %T", b)
    }
    return
}

func (date DateTime) Value() (driver.Value, error) {
    return date.Time.Format("2006-01-02 15:04:05"), nil
}