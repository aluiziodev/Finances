package parser

import (
	"errors"
	"finances/internal/models"
	"mime/multipart"
	"strings"

	"github.com/gocarina/gocsv"
)

func ParserCSVtoModels(file multipart.File, handler *multipart.FileHeader, req models.RequestFatura) (models.Fatura, error) {

	if !strings.HasSuffix(strings.ToLower(handler.Filename), ".csv") {
		return models.Fatura{}, errors.New("arquivo deve ser .csv")
	}

	var bills []models.Bill
	if err := gocsv.Unmarshal(file, &bills); err != nil {
		return models.Fatura{}, err
	}

	formatBills(&bills)

	fatura := models.Fatura{
		Description: req.Description,
		Bank:        req.Bank,
		Bills:       bills,
		Status:      req.Status,
	}

	return fatura, nil

}

func formatBills(bills *[]models.Bill) {
	for i, bill := range *bills {
		if bill.VerifyPayment() {
			*bills = append((*bills)[:i], (*bills)[i+1:]...)
			continue
		}
	}

}
