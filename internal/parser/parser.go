package parser

import (
	"finances/internal/dto"
	"finances/internal/models"
	"fmt"
	"mime/multipart"
	"strings"
)

type ParserInterface interface {
	ParserCSVtoModels(file multipart.File, handler *multipart.FileHeader, req dto.RequestFatura) (models.Fatura, error)
}

var parsers = map[string]func() ParserInterface{
	"nubank": func() ParserInterface { return &NubankParser{} },
}

func newParser(bank string) (ParserInterface, error) {
	key := strings.ToLower(strings.TrimSpace(bank))

	if factory, ok := parsers[key]; ok {
		return factory(), nil
	}

	return nil, fmt.Errorf("parser para o banco %q não encontrado", bank)
}

func ParserCSVtoModels(file multipart.File, handler *multipart.FileHeader, req dto.RequestFatura) (models.Fatura, error) {
	parser, err := newParser(req.Bank)
	if err != nil {
		return models.Fatura{}, err
	}

	return parser.ParserCSVtoModels(file, handler, req)
}
