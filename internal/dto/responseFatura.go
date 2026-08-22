package dto

import "finances/internal/models"

type ResponseFatura struct {
	Description string         `json:"description"`
	Bank        string         `json:"bank"`
	Bills       []ResponseBill `json:"bills"`
	Status      string         `json:"status"`
	Total       float64        `json:"total"`
}

func NewResponseFatura(fatura *models.Fatura) ResponseFatura {
	return ResponseFatura{
		Description: fatura.Description,
		Bank:        fatura.Bank,
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
