package dto

import "finances/internal/models"

type ResponseBill struct {
	Date     string  `json:"date"`
	Title    string  `json:"title"`
	Amount   float64 `json:"amount"`
	Method   string  `json:"method"`
	Category string  `json:"category"`
}

func NewResponseBill(bill *models.Bill) ResponseBill {
	return ResponseBill{
		Date:     bill.Date,
		Title:    bill.Title,
		Amount:   float64(bill.Amount),
		Method:   bill.Method,
		Category: bill.Category,
	}
}

func BillsToResponseBills(bills *[]models.Bill) []ResponseBill {
	responseBills := make([]ResponseBill, len(*bills))

	for i, bill := range *bills {
		responseBills[i] = NewResponseBill(&bill)
	}

	return responseBills
}
