package dto

import "finances/internal/models"

var validCategories = []string{
	"transporte",
	"alimentação",
	"mercado",
	"saúde",
	"assinaturas",
	"vestuario",
	"celular",
	"entretenimento",
	"varejo",
	"moradia",
	"educacao",
	"viagem",
	"servicos",
	"outros",
}

type Summary struct {
	Fatura          ResponseFatura     `json:"fatura"`
	TotalParcelado  float64            `json:"total_parcelado"`
	TotalFixo       float64            `json:"total_fixo"`
	TotalGeneral    float64            `json:"total_geral"`
	TotalByCategory map[string]float64 `json:"total_by_category"`
}

func (s *Summary) CalculateTotal(fatura *models.Fatura) {
	s.initialize()
	s.TotalGeneral = fatura.Total
	s.Fatura = NewResponseFatura(fatura)
	for _, bill := range fatura.Bills {
		if bill.Method == "parcelado" {
			s.TotalParcelado += float64(bill.Amount)
		} else {
			s.TotalFixo += float64(bill.Amount)
		}

		s.TotalByCategory[bill.Category] += float64(bill.Amount)
	}
}

func (s *Summary) initialize() {
	s.TotalParcelado = 0
	s.TotalFixo = 0
	s.TotalGeneral = 0

	s.TotalByCategory = make(map[string]float64)

	for _, category := range validCategories {
		s.TotalByCategory[category] = 0
	}
}
