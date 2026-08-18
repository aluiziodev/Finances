package repository

import (
	"database/sql"
	"finances/internal/connection"
	"finances/internal/models"
)

type BillRepository struct {
	db *sql.DB
}

func NewBillRepository() *BillRepository {
	return &BillRepository{connection.DB}
}

func (repo *BillRepository) Create(bill models.Bill, fatura_id string) error {
	_, err := repo.db.Exec(`
		INSERT INTO bill (id, title, date, amount, fatura)
		VALUES ($1, $2, $3, $4, $5)
	`, bill.Id, bill.Title, bill.Date, bill.Amount, fatura_id)
	return err
}

func (repo *BillRepository) GetAllByFaturaId(fatura_id string) ([]models.Bill, error) {
	rows, err := repo.db.Query(`
		SELECT b.id, b.title, b.date, b.amount FROM bill b
		WHERE b.fatura = $1
	`, fatura_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bills []models.Bill

	for rows.Next() {
		var bill models.Bill
		if err = rows.Scan(&bill.Id, &bill.Title, &bill.Date, &bill.Amount); err != nil {
			return nil, err
		}

		bills = append(bills, bill)
	}

	return bills, nil
}
