package parser

import (
	"errors"
	"finances/internal/models"
	"mime/multipart"
	"strings"

	"github.com/gocarina/gocsv"
)

func ParserCSVtoModels(file multipart.File, handler *multipart.FileHeader, description string, status string) (models.Fatura, error) {

	if !strings.HasSuffix(strings.ToLower(handler.Filename), ".csv") {
		return models.Fatura{}, errors.New("arquivo deve ser .csv")
	}

	var bills []models.Bill
	if err := gocsv.Unmarshal(file, &bills); err != nil {
		return models.Fatura{}, err
	}

	fatura := models.Fatura{
		Description: description,
		Bills:       bills,
		Status:      status,
	}

	fatura.CalculateTotal()
	if err := fatura.Validate(); err != nil {
		return models.Fatura{}, err
	}

	return fatura, nil

}
