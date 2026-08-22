package parser

import (
	"bytes"
	"finances/internal/dto"
	"mime/multipart"
	"strings"
	"testing"
)

func createCSVMultipartFile(t *testing.T, filename string, content string) (multipart.File, *multipart.FileHeader) {
	t.Helper()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile("csv", filename)
	if err != nil {
		t.Fatalf("erro ao criar arquivo multipart: %v", err)
	}

	if _, err = part.Write([]byte(content)); err != nil {
		t.Fatalf("erro ao escrever arquivo multipart: %v", err)
	}

	if err = writer.Close(); err != nil {
		t.Fatalf("erro ao fechar writer: %v", err)
	}

	reader := multipart.NewReader(bytes.NewReader(body.Bytes()), writer.Boundary())
	form, err := reader.ReadForm(10 << 20)
	if err != nil {
		t.Fatalf("erro ao ler form multipart: %v", err)
	}

	if len(form.File["csv"]) == 0 {
		t.Fatal("arquivo csv nao encontrado no form")
	}

	file, err := form.File["csv"][0].Open()
	if err != nil {
		t.Fatalf("erro ao abrir arquivo csv: %v", err)
	}

	return file, form.File["csv"][0]
}

func TestParserCSVtoModels_ValidCSV(t *testing.T) {
	file, handler := createCSVMultipartFile(t, "fatura.csv", "date,title,amount\n2024-01-01,Compra,\"100,50\"\n")
	defer file.Close()

	fatura, err := ParserCSVtoModels(file, handler, dto.RequestFatura{
		Description: "Fatura de teste",
		Bank:        "Banco Teste",
		Status:      "paid",
	})
	if err != nil {
		t.Fatalf("esperava parse valido, mas recebeu erro: %v", err)
	}

	if fatura.Description != "Fatura de teste" {
		t.Fatalf("descricao inesperada: %s", fatura.Description)
	}

	if len(fatura.Bills) != 1 {
		t.Fatalf("quantidade de bills inesperada: %d", len(fatura.Bills))
	}

	if fatura.Bills[0].Amount != 100.5 {
		t.Fatalf("total inesperado: %v", fatura.Bills[0].Amount)
	}
}

func TestParserCSVtoModels_InvalidExtension(t *testing.T) {
	file, handler := createCSVMultipartFile(t, "fatura.txt", "date,title,amount\n2024-01-01,Compra,100.50\n")
	defer file.Close()

	_, err := ParserCSVtoModels(file, handler, dto.RequestFatura{
		Description: "Fatura de teste",
		Bank:        "Banco Teste",
		Status:      "paid",
	})
	if err == nil {
		t.Fatal("esperava erro para extensao invalida")
	}

	if !strings.Contains(err.Error(), "arquivo deve ser .csv") {
		t.Fatalf("mensagem inesperada: %v", err)
	}
}
