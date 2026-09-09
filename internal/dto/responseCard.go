package dto

import "finances/internal/models"

type ResponseCard struct {
	Name   string  `json:"name"`
	Bank   string  `json:"bank"`
	Limit  float64 `json:"limit"`
	DueDay int     `json:"due_day"`
}

func NewResponseCard(card *models.Card) ResponseCard {
	return ResponseCard{
		Name:   card.Name,
		Bank:   card.Bank,
		Limit:  card.Limit,
		DueDay: card.DueDay,
	}
}

func CardsToResponseCards(cards *[]models.Card) []ResponseCard {
	responseCards := make([]ResponseCard, len(*cards))

	for i, card := range *cards {
		responseCards[i] = NewResponseCard(&card)
	}

	return responseCards
}
