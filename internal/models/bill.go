package models

import (
	"fmt"
	"strconv"
	"strings"
)

type Bill struct {
	Id       string `json:"id"`
	Date     string `json:"date" csv:"date"`
	Title    string `json:"title" csv:"title"`
	Amount   Value  `json:"amount" csv:"amount"`
	Method   string `json:"method"`
	Category string `json:"category"`
}

func (b *Bill) DefineMethod() {
	if strings.Contains(b.Title, " - Parcela ") {
		b.Method = "parcelado"
	} else {
		b.Method = "fixo"
	}
}

func (b *Bill) Validate() error {
	if b.Title == "" {
		return fmt.Errorf("Campo title obrigatorio!!")
	}

	if b.Method != "parcelado" && b.Method != "fixo" {
		return fmt.Errorf("Campo method invalido!")
	}

	return nil
}

func (b *Bill) VerifyPayment() bool {
	if b.Title == "Pagamento recebido" {
		return true
	}
	return false
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
