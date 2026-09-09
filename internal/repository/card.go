package repository

import (
	"database/sql"
	"finances/internal/connection"
	"finances/internal/models"
)

type CardRepositoryInterface interface {
	Create(card models.Card) error
	GetAll() ([]models.Card, error)
	Get(id string) (models.Card, error)
	GetBankByCardId(cardId string) (string, error)
	Delete(id string) error
	Update(card models.Card) error
}

type CardRepository struct {
	db *sql.DB
}

func NewCardRepository() *CardRepository {
	return &CardRepository{db: connection.DB}
}

func (repo *CardRepository) Create(card models.Card) error {
	_, err := repo.db.Exec(`
		INSERT INTO card (id, name, bank, credit_limit, due_day)
		VALUES ($1, $2, $3, $4, $5)
	`, card.Id, card.Name, card.Bank, card.Limit, card.DueDay)
	return err
}

func (repo *CardRepository) GetAll() ([]models.Card, error) {
	rows, err := repo.db.Query(`
		SELECT c.id, c.name, c.bank, c.credit_limit, c.due_day FROM card c
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []models.Card

	for rows.Next() {
		var card models.Card
		if err = rows.Scan(&card.Id, &card.Name, &card.Bank, &card.Limit, &card.DueDay); err != nil {
			return nil, err
		}

		cards = append(cards, card)
	}

	return cards, nil
}

func (repo *CardRepository) Get(id string) (models.Card, error) {
	row := repo.db.QueryRow(`
		SELECT c.id, c.name, c.bank, c.credit_limit, c.due_day FROM card c
		WHERE c.id = $1
	`, id)

	var card models.Card
	if err := row.Scan(&card.Id, &card.Name, &card.Bank, &card.Limit, &card.DueDay); err != nil {
		return models.Card{}, err
	}

	return card, nil
}

func (repo *CardRepository) GetBankByCardId(cardId string) (string, error) {
	row := repo.db.QueryRow(`
		SELECT c.bank FROM card c
		WHERE c.id = $1
	`, cardId)

	var bank string
	if err := row.Scan(&bank); err != nil {
		return "", err
	}

	return bank, nil
}
func (repo *CardRepository) Delete(id string) error {
	_, err := repo.db.Exec(`
		DELETE FROM card WHERE id = $1
	`, id)
	return err
}

func (repo *CardRepository) Update(card models.Card) error {
	_, err := repo.db.Exec(`
		UPDATE card SET name = $1, bank = $2, credit_limit = $3, due_day = $4 WHERE id = $5
	`, card.Name, card.Bank, card.Limit, card.DueDay, card.Id)
	return err
}
