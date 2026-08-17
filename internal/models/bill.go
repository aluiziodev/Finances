package models

import (
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type Bill struct {
	Id     uuid.UUID `json:"id"`
	Date   string    `csv:"date"`
	Title  string    `csv:"title"`
	Amount Value     `csv:"amount"`
}

type Value float64

func (v *Value) UnmarshalCSV(str string) error {
	value := strings.TrimSpace(str)
	value = strings.ReplaceAll(str, " ", "")
	value = strings.ReplaceAll(value, ".", "")
	value = strings.ReplaceAll(value, ",", ".")

	parsedValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return err
	}

	*v = Value(parsedValue)
	return nil
}

func (v Value) MarshalCSV() (string, error) {
	return strconv.FormatFloat(float64(v), 'f', 2, 64), nil
}
