package repository

import (
	"database/sql"
	"finances/internal/connection"
	"finances/internal/models"
)

type FaturaRepositoryInterface interface {
	Create(fatura models.Fatura) error
	GetAll() ([]models.Fatura, error)
	Get(id string) (models.Fatura, error)
	Delete(id string) error
}

type FaturaRepository struct {
	db *sql.DB
}

func NewFaturaRepository() *FaturaRepository {
	return &FaturaRepository{connection.DB}
}

func (repo *FaturaRepository) Create(fatura models.Fatura) error {
	_, err := repo.db.Exec(`
		INSERT INTO fatura (id, description, bank, status, total)
		VALUES ($1, $2, $3, $4, $5)
	`, fatura.Id, fatura.Description, fatura.Bank, fatura.Status, fatura.Total)
	return err
}

func (repo *FaturaRepository) GetAll() ([]models.Fatura, error) {
	rows, err := repo.db.Query(`
		SELECT f.id, f.description, f.bank, f.status, f.total FROM fatura f
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var faturas []models.Fatura

	for rows.Next() {
		var fatura models.Fatura
		if err = rows.Scan(&fatura.Id, &fatura.Description, &fatura.Bank, &fatura.Status, &fatura.Total); err != nil {
			return nil, err
		}

		faturas = append(faturas, fatura)
	}

	return faturas, nil
}

func (repo *FaturaRepository) Get(id string) (models.Fatura, error) {
	row, err := repo.db.Query(`
		SELECT f.id, f.description,f.bank, f.status, f.total FROM fatura f
		WHERE f.id = $1
	`, id)
	if err != nil {
		return models.Fatura{}, err
	}
	defer row.Close()

	var fatura models.Fatura

	if row.Next() {
		if err = row.Scan(&fatura.Id, &fatura.Description, &fatura.Bank, &fatura.Status, &fatura.Total); err != nil {
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
