package service

import (
	"bytes"
	"finances/internal/dto"
	"finances/internal/models"
	"mime/multipart"
	"testing"
)

type mockFaturaRepo struct {
	created   []models.Fatura
	list      []models.Fatura
	byId      models.Fatura
	deletedID []string
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
	f.deletedID = append(f.deletedID, id)
	return nil
}

type mockBillRepo struct {
	created   []models.Bill
	list      []models.Bill
	parcelado []models.Bill
	fixo      []models.Bill
	filtered  []models.Bill
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

func (f *mockBillRepo) GetParcelado(id string) ([]models.Bill, error) {
	return f.parcelado, nil
}

func (f *mockBillRepo) GetFixo(id string) ([]models.Bill, error) {
	return f.fixo, nil
}
func (f *mockBillRepo) GetBillsByCategory(fatura_id string, category string) ([]models.Bill, error) {
	return f.filtered, nil
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

	id, err := service.CreateFatura(file, handler, dto.RequestFatura{
		Description: "Fatura de teste",
		Bank:        "nubank",
		Status:      "paid",
	})
	if err != nil {
		t.Fatalf("esperava criacao bem-sucedida, mas recebeu erro: %v", err)
	}

	if *id == "" {
		t.Fatal("id da fatura nao pode estar vazio")
	}

	if len(faturaRepo.created) != 1 {
		t.Fatalf("quantidade de faturas criadas inesperada: %d", len(faturaRepo.created))
	}

	if len(billRepo.created) != 1 {
		t.Fatalf("quantidade de bills criadas inesperada: %d", len(billRepo.created))
	}
}

// Teste para a função GetFatura
func TestFaturaService_GetFatura(t *testing.T) {
	faturaRepo := &mockFaturaRepo{
		byId: models.Fatura{
			Id:          "fatura-1",
			Description: "Fatura de teste",
			Status:      "paid",
			Total:       50.0,
		},
	}
	billRepo := &mockBillRepo{
		list: []models.Bill{{
			Id:       "bill-1",
			Title:    "Compra",
			Amount:   50.0,
			Method:   "fixo",
			Category: "varejo",
		}},
	}
	service := NewFaturaServiceWithDependencies(faturaRepo, billRepo)

	result, err := service.GetFatura("fatura-1")
	if err != nil {
		t.Fatalf("esperava busca bem-sucedida, mas recebeu erro: %v", err)
	}

	if len(result.Fatura.Bills) != 1 {
		t.Fatalf("quantidade de bills retornadas inesperada: %d", len(result.Fatura.Bills))
	}

	if result.Fatura.Description != "Fatura de teste" {
		t.Fatalf("descrição da fatura inesperada: %s", result.Fatura.Description)
	}

	if result.TotalGeneral != 50.0 {
		t.Fatalf("total geral inesperado: %.2f", result.TotalGeneral)
	}

	if result.TotalFixo != 50.0 {
		t.Fatalf("total fixo inesperado: %.2f", result.TotalFixo)
	}

	if result.TotalByCategory["varejo"] != 50.0 {
		t.Fatalf("total por categoria inesperado: %.2f", result.TotalByCategory["varejo"])
	}
}

func TestFaturaService_DeleteFatura(t *testing.T) {
	faturaRepo := &mockFaturaRepo{}
	billRepo := &mockBillRepo{}
	service := NewFaturaServiceWithDependencies(faturaRepo, billRepo)

	if err := service.DeleteFatura("fatura-1"); err != nil {
		t.Fatalf("esperava exclusao bem-sucedida, mas recebeu erro: %v", err)
	}

	if len(faturaRepo.deletedID) != 1 {
		t.Fatalf("quantidade de exclusoes inesperada: %d", len(faturaRepo.deletedID))
	}

	if faturaRepo.deletedID[0] != "fatura-1" {
		t.Fatalf("id excluido inesperado: %s", faturaRepo.deletedID[0])
	}
}

func TestFaturaService_GetFaturaParcelado(t *testing.T) {
	billRepo := &mockBillRepo{
		parcelado: []models.Bill{{Id: "bill-1", Title: "Uber - Parcela 1", Amount: 80.0, Method: "parcelado"}, {Id: "bill-2", Title: "Uber - Parcela 2", Amount: 20.0, Method: "parcelado"}},
	}
	service := NewFaturaServiceWithDependencies(&mockFaturaRepo{}, billRepo)

	result, err := service.GetFaturaParcelado("fatura-1")
	if err != nil {
		t.Fatalf("esperava busca de parcelado bem-sucedida, mas recebeu erro: %v", err)
	}

	if result.Total != 100.0 {
		t.Fatalf("total inesperado: %.2f", result.Total)
	}

	if len(result.Bills) != 2 {
		t.Fatalf("quantidade de bills parceladas inesperada: %d", len(result.Bills))
	}
}

func TestFaturaService_GetFaturaFixo(t *testing.T) {
	billRepo := &mockBillRepo{
		fixo: []models.Bill{{Id: "bill-1", Title: "Netflix", Amount: 30.0, Method: "fixo"}, {Id: "bill-2", Title: "Mercado", Amount: 70.0, Method: "fixo"}},
	}
	service := NewFaturaServiceWithDependencies(&mockFaturaRepo{}, billRepo)

	result, err := service.GetFaturaFixo("fatura-1")
	if err != nil {
		t.Fatalf("esperava busca de fixos bem-sucedida, mas recebeu erro: %v", err)
	}

	if result.Total != 100.0 {
		t.Fatalf("total inesperado: %.2f", result.Total)
	}

	if len(result.Bills) != 2 {
		t.Fatalf("quantidade de bills fixas inesperada: %d", len(result.Bills))
	}
}
