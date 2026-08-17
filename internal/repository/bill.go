package repository

import (
	"database/sql"
	"finances/internal/connection"
	"finances/internal/models"

	"github.com/google/uuid"
)

type BillRepository struct {
	db *sql.DB
}

func NewBillRepository() *BillRepository {
	return &BillRepository{connection.DB}
}

func (repo *BillRepository) Create(bill models.Bill, fatura_id uuid.UUID) error {
	_, err := repo.db.Exec(`
		INSERT INTO bill (id, title, date, amount, fatura)
		VALUES ($1, $2, $3, $4, $5)
	`, bill.Id, bill.Title, bill.Date, bill.Amount, fatura_id)
	return err
}
