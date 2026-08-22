package models

import (
	"errors"
)

type Fatura struct {
	Id          string
	Description string
	Bank        string
	Bills       []Bill
	Total       float64
	Status      string
}

func (f *Fatura) CalculateTotal() {
	var total float64
	for _, bill := range f.Bills {
		total += float64(bill.Amount)
	}
	f.Total = total
}

func (f *Fatura) Validate() error {
	if f.Description == "" {
		return errors.New("Campo descriçao obrigatorio!!")
	}

	if f.Status != "peding" && f.Status != "paid" {
		return errors.New("Campo status invalido!")
	}

	if f.Total < 0.0 {
		return errors.New("Campo total nao pode ser negativo!")
	}

	return nil

}
