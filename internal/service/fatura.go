package service

import (
	"finances/internal/models"
	"finances/internal/parser"
	"finances/internal/repository"
	"mime/multipart"

	"github.com/google/uuid"
)

func CreateFatura(file multipart.File, handler *multipart.FileHeader,
	req models.RequestFatura) (*models.Fatura, error) {

	fatura, err := parser.ParserCSVtoModels(file, handler, req.Description, req.Status)
	if err != nil {
		return nil, err
	}
	fatura.Id = uuid.NewString()

	repo := repository.NewFaturaRepository()
	if err := repo.Create(fatura); err != nil {
		return nil, err
	}

	repo_bill := repository.NewBillRepository()

	for _, bill := range fatura.Bills {
		bill.Id = uuid.NewString()
		if err := repo_bill.Create(bill, fatura.Id); err != nil {
			return nil, err
		}
	}

	return &fatura, nil
}

func GetAllFaturas() ([]models.Fatura, error) {
	repo := repository.NewFaturaRepository()
	faturas, err := repo.GetAll()
	if err != nil {
		return nil, err
	}

	return faturas, nil
}

func GetFatura(id string) (models.Fatura, error) {
	repo := repository.NewFaturaRepository()
	fatura, err := repo.Get(id)
	if err != nil {
		return models.Fatura{}, err
	}

	repo_bill := repository.NewBillRepository()
	bills, err := repo_bill.GetAllByFaturaId(id)
	if err != nil {
		return models.Fatura{}, err
	}

	fatura.Bills = bills

	return fatura, nil
}
