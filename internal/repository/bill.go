package repository

import (
	"database/sql"
	"finances/internal/connection"
	"finances/internal/models"
)

type BillRepositoryInterface interface {
	Create(bill models.Bill, fatura_id string) error
	GetAllByFaturaId(fatura_id string) ([]models.Bill, error)
	Delete(id string) error
	GetParcelado(fatura_id string) ([]models.Bill, error)
	GetFixo(fatura_id string) ([]models.Bill, error)
	GetBillsByCategory(fatura_id string, category string) ([]models.Bill, error)
}

type BillRepository struct {
	db *sql.DB
}

func NewBillRepository() *BillRepository {
	return &BillRepository{connection.DB}
}

func (repo *BillRepository) Create(bill models.Bill, fatura_id string) error {
	_, err := repo.db.Exec(`
		INSERT INTO bill (id, title, date, amount, fatura, method, category)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, bill.Id, bill.Title, bill.Date, bill.Amount, fatura_id, bill.Method, bill.Category)
	return err
}

func (repo *BillRepository) GetAllByFaturaId(fatura_id string) ([]models.Bill, error) {
	rows, err := repo.db.Query(`
		SELECT b.id, b.title, b.date, b.amount, b.category, b.method FROM bill b
		WHERE b.fatura = $1
	`, fatura_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bills []models.Bill

	for rows.Next() {
		var bill models.Bill
		if err = rows.Scan(&bill.Id, &bill.Title, &bill.Date, &bill.Amount, &bill.Category, &bill.Method); err != nil {
			return nil, err
		}

		bills = append(bills, bill)
	}

	return bills, nil
}

func (repo *BillRepository) Delete(id string) error {
	_, err := repo.db.Exec(`
		DELETE FROM bill WHERE id = $1
	`, id)
	return err
}

func (repo *BillRepository) GetParcelado(fatura_id string) ([]models.Bill, error) {
	rows, err := repo.db.Query(`
		SELECT b.id, b.title, b.date, b.amount, b.category, b.method FROM bill b
		WHERE b.fatura = $1 AND b.method = 'parcelado'
	`, fatura_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bills []models.Bill

	for rows.Next() {
		var bill models.Bill
		if err = rows.Scan(&bill.Id, &bill.Title, &bill.Date, &bill.Amount, &bill.Category, &bill.Method); err != nil {
			return nil, err
		}

		bills = append(bills, bill)
	}

	return bills, nil
}

func (repo *BillRepository) GetFixo(fatura_id string) ([]models.Bill, error) {
	rows, err := repo.db.Query(`
		SELECT b.id, b.title, b.date, b.amount, b.category, b.method FROM bill b
		WHERE b.fatura = $1 AND b.method = 'fixo'
	`, fatura_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bills []models.Bill

	for rows.Next() {
		var bill models.Bill
		if err = rows.Scan(&bill.Id, &bill.Title, &bill.Date, &bill.Amount, &bill.Category, &bill.Method); err != nil {
			return nil, err
		}

		bills = append(bills, bill)
	}

	return bills, nil
}

func (repo *BillRepository) GetBillsByCategory(fatura_id string, category string) ([]models.Bill, error) {
	rows, err := repo.db.Query(`
		SELECT b.id, b.title, b.date, b.amount, b.category, b.method FROM bill b
		WHERE b.fatura = $1 AND b.category = $2
	`, fatura_id, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bills []models.Bill

	for rows.Next() {
		var bill models.Bill
		if err = rows.Scan(&bill.Id, &bill.Title, &bill.Date, &bill.Amount, &bill.Category, &bill.Method); err != nil {
			return nil, err
		}

		bills = append(bills, bill)
	}

	return bills, nil
}
