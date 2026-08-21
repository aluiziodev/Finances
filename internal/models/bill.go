package models

import (
	"strconv"
	"strings"
)

type Bill struct {
	Id     string `json:"id"`
	Date   string `json:"date" csv:"date"`
	Title  string `json:"title" csv:"title"`
	Amount Value  `json:"amount" csv:"amount"`
	Method string `json:"method"`
}

func (b *Bill) DefineMethod() {
	if strings.Contains(b.Title, " - Parcela ") {
		b.Method = "parcelado"
	} else {
		b.Method = "fixo"
	}
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
