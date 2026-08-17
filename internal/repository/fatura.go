package repository

import (
	"database/sql"
	"finances/internal/connection"
	"finances/internal/models"
)

type FaturaRepository struct {
	db *sql.DB
}

func NewFaturaRepository() *FaturaRepository {
	return &FaturaRepository{connection.DB}
}

func (repo *FaturaRepository) Create(fatura models.Fatura) error {
	_, err := repo.db.Exec(`
		INSERT INTO fatura (id, description, status, total)
		VALUES ($1, $2, $3, $4)
	`, fatura.Id, fatura.Description, fatura.Status, fatura.Total)
	return err
}
