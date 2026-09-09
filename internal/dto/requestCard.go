package dto

import "finances/internal/models"

type RequestCard struct {
	Name   string  `json:"name"`
	Bank   string  `json:"bank"`
	Limit  float64 `json:"limit"`
	DueDay int     `json:"due_day"`
}

func (r *RequestCard) ToModel() models.Card {
	return models.Card{
		Name:   r.Name,
		Bank:   r.Bank,
		Limit:  r.Limit,
		DueDay: r.DueDay,
	}
}
