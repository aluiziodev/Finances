package parser

import (
	"finances/internal/models"
	"os"

	"github.com/gocarina/gocsv"
)

func ParserCSVtoModels(path string, description string, status string) (*models.Fatura, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var bills []models.Bill
	if err := gocsv.UnmarshalFile(file, &bills); err != nil {
		return nil, err
	}

	fatura := models.Fatura{
		Description: description,
		Bills:       bills,
		Status:      status,
	}
	fatura.CalculateTotal()
	if err := fatura.Validate(); err != nil {
		return nil, err
	}

	return &fatura, nil
}
