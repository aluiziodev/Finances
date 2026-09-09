package repository

import (
	"database/sql"
	"finances/internal/connection"
	"finances/internal/models"
)

type FaturaRepositoryInterface interface {
	Create(fatura models.Fatura) error
	GetAll(card_id string) ([]models.Fatura, error)
	Get(id string) (models.Fatura, error)
	Delete(id string) error
	UpdateStatusPaid(id string) error
	UpdateStatusPending(id string) error
}

type FaturaRepository struct {
	db *sql.DB
}

func NewFaturaRepository() *FaturaRepository {
	return &FaturaRepository{connection.DB}
}

func (repo *FaturaRepository) Create(fatura models.Fatura) error {
	_, err := repo.db.Exec(`
		INSERT INTO fatura (id, card_id, year, month, description, status, total)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, fatura.Id, fatura.CardID, fatura.Year, fatura.Month, fatura.Description, fatura.Status, fatura.Total)
	return err
}

func (repo *FaturaRepository) GetAll(card_id string) ([]models.Fatura, error) {
	rows, err := repo.db.Query(`
		SELECT f.id, f.card_id, f.year, f.month, f.description, f.status, f.total FROM fatura f
		WHERE f.card_id = $1
	`, card_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var faturas []models.Fatura

	for rows.Next() {
		var fatura models.Fatura
		if err = rows.Scan(&fatura.Id, &fatura.CardID, &fatura.Year, &fatura.Month, &fatura.Description, &fatura.Status, &fatura.Total); err != nil {
			return nil, err
		}

		faturas = append(faturas, fatura)
	}

	return faturas, nil
}

func (repo *FaturaRepository) Get(id string) (models.Fatura, error) {
	row, err := repo.db.Query(`
		SELECT f.id, f.card_id, f.year, f.month, f.description, f.status, f.total FROM fatura f
		WHERE f.id = $1
	`, id)
	if err != nil {
		return models.Fatura{}, err
	}
	defer row.Close()

	var fatura models.Fatura

	if row.Next() {
		if err = row.Scan(&fatura.Id, &fatura.CardID, &fatura.Year, &fatura.Month, &fatura.Description, &fatura.Status, &fatura.Total); err != nil {
			return models.Fatura{}, err
		}
		return fatura, nil
	}
	return models.Fatura{}, sql.ErrNoRows
}

func (repo *FaturaRepository) Delete(id string) error {
	_, err := repo.db.Exec(`
		DELETE FROM fatura WHERE id = $1
	`, id)
	return err
}

func (repo *FaturaRepository) UpdateStatusPaid(id string) error {
	_, err := repo.db.Exec(`
		UPDATE fatura
		SET status = 'paid'
		WHERE id = $1
	`)
	return err
}

func (repo *FaturaRepository) UpdateStatusPending(id string) error {
	_, err := repo.db.Exec(`
		UPDATE fatura
		SET status = 'pending'
		WHERE id = $1
	`)
	return err
}
