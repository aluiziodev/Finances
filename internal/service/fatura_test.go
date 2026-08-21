package service

import (
	"bytes"
	"finances/internal/models"
	"mime/multipart"
	"testing"
)

type mockFaturaRepo struct {
	created []models.Fatura
	list    []models.Fatura
	byId    models.Fatura
}

func (f *mockFaturaRepo) Create(fatura models.Fatura) error {
	f.created = append(f.created, fatura)
	return nil
}

func (f *mockFaturaRepo) GetAll() ([]models.Fatura, error) {
	return f.list, nil
}

func (f *mockFaturaRepo) Get(id string) (models.Fatura, error) {
	return f.byId, nil
}

func (f *mockFaturaRepo) Delete(id string) error {
	return nil
}

type mockBillRepo struct {
	created []models.Bill
	list    []models.Bill
}

func (f *mockBillRepo) Create(bill models.Bill, faturaID string) error {
	f.created = append(f.created, bill)
	return nil
}

func (f *mockBillRepo) GetAllByFaturaId(faturaID string) ([]models.Bill, error) {
	return f.list, nil
}

func (f *mockBillRepo) Delete(id string) error {
	return nil
}

func createMultipartCSV(t *testing.T, filename string, content string) (multipart.File, *multipart.FileHeader) {
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
		t.Fatalf("erro ao ler form: %v", err)
	}

	if len(form.File["csv"]) == 0 {
		t.Fatal("arquivo csv nao encontrado")
	}

	file, err := form.File["csv"][0].Open()
	if err != nil {
		t.Fatalf("erro ao abrir csv: %v", err)
	}

	return file, form.File["csv"][0]
}

// Teste para a função CreateFatura
func TestFaturaService_CreateFatura(t *testing.T) {
	faturaRepo := &mockFaturaRepo{}
	billRepo := &mockBillRepo{}
	service := NewFaturaServiceWithDependencies(faturaRepo, billRepo)

	file, handler := createMultipartCSV(t, "fatura.csv", "date,title,amount\n2024-01-01,Compra,50.00\n")
	defer file.Close()

	result, err := service.CreateFatura(file, handler, models.RequestFatura{
		Description: "Fatura de teste",
		Status:      "paid",
	})
	if err != nil {
		t.Fatalf("esperava criacao bem-sucedida, mas recebeu erro: %v", err)
	}

	if result.Id == "" {
		t.Fatal("id da fatura nao pode estar vazio")
	}

	if len(faturaRepo.created) != 1 {
		t.Fatalf("quantidade de faturas criadas inesperada: %d", len(faturaRepo.created))
	}

	if len(billRepo.created) != 1 {
		t.Fatalf("quantidade de bills criadas inesperada: %d", len(billRepo.created))
	}
}

// Teste para a função GetAllFaturas
func TestFaturaService_GetFatura(t *testing.T) {
	faturaRepo := &mockFaturaRepo{
		byId: models.Fatura{Id: "fatura-1", Description: "Fatura de teste", Status: "paid"},
	}
	billRepo := &mockBillRepo{
		list: []models.Bill{{Id: "bill-1", Title: "Compra", Amount: 50.0}},
	}
	service := NewFaturaServiceWithDependencies(faturaRepo, billRepo)

	result, err := service.GetFatura("fatura-1")
	if err != nil {
		t.Fatalf("esperava busca bem-sucedida, mas recebeu erro: %v", err)
	}

	if result.Id != "fatura-1" {
		t.Fatalf("id inesperado: %s", result.Id)
	}

	if len(result.Bills) != 1 {
		t.Fatalf("quantidade de bills retornadas inesperada: %d", len(result.Bills))
	}
}
