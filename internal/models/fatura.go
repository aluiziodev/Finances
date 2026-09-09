package models

import (
	"errors"
)

type Fatura struct {
	Id          string
	CardID      string
	Year        int
	Month       int
	Description string
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
	if f.CardID == "" {
		return errors.New("Card id is required")
	}

	if f.Description == "" {
		return errors.New("Description is required")
	}

	if f.Year < 2000 {
		return errors.New("Year is invalid")
	}

	if f.Month < 1 || f.Month > 12 {
		return errors.New("Month must be between 1 and 12")
	}

	if f.Status != "pending" && f.Status != "paid" {
		return errors.New("Invalid status")
	}

	if f.Total < 0.0 {
		return errors.New("Total cannot be negative")
	}

	return nil
}
