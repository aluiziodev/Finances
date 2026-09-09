package parser

import (
	"errors"
	"finances/internal/dto"
	"finances/internal/models"
	"mime/multipart"
	"strings"

	"github.com/gocarina/gocsv"
)

type NubankParser struct{}

func (nu *NubankParser) ParserCSVtoModels(file multipart.File, handler *multipart.FileHeader, req dto.RequestFatura) (models.Fatura, error) {

	if !strings.HasSuffix(strings.ToLower(handler.Filename), ".csv") {
		return models.Fatura{}, errors.New("arquivo deve ser .csv")
	}

	var bills []models.Bill
	if err := gocsv.Unmarshal(file, &bills); err != nil {
		return models.Fatura{}, err
	}

	nu.formatBills(&bills)

	fatura := models.Fatura{
		Year:        req.Year,
		Month:       req.Month,
		Description: req.Description,
		Bills:       bills,
		Status:      req.Status,
	}

	return fatura, nil

}

func (nu *NubankParser) formatBills(bills *[]models.Bill) {
	for i, bill := range *bills {
		if bill.VerifyPayment() {
			*bills = append((*bills)[:i], (*bills)[i+1:]...)
			continue
		}
	}

}
