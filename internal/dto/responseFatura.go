package dto

import "finances/internal/models"

type ResponseFatura struct {
	Year        int            `json:"year"`
	Month       int            `json:"month"`
	Description string         `json:"description"`
	Bills       []ResponseBill `json:"bills"`
	Status      string         `json:"status"`
	Total       float64        `json:"total"`
}

func NewResponseFatura(fatura *models.Fatura) ResponseFatura {
	return ResponseFatura{
		Year:        fatura.Year,
		Month:       fatura.Month,
		Description: fatura.Description,
		Bills:       BillsToResponseBills(&fatura.Bills),
		Status:      fatura.Status,
		Total:       fatura.Total,
	}
}

func FaturasToResponseFaturas(faturas *[]models.Fatura) []ResponseFatura {
	responseFaturas := make([]ResponseFatura, len(*faturas))

	for i, fatura := range *faturas {
		responseFaturas[i] = NewResponseFatura(&fatura)
	}

	return responseFaturas
}
