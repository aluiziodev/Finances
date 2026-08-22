package controllers

import (
	"bytes"
	"encoding/json"
	"finances/internal/dto"
	"finances/internal/models"
	"finances/internal/service"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockControllerFaturaRepo struct {
	created []models.Fatura
	list    []models.Fatura
	byId    models.Fatura
}

func (f *mockControllerFaturaRepo) Create(fatura models.Fatura) error {
	f.created = append(f.created, fatura)
	return nil
}

func (f *mockControllerFaturaRepo) GetAll() ([]models.Fatura, error) {
	return f.list, nil
}

func (f *mockControllerFaturaRepo) Get(id string) (models.Fatura, error) {
	return f.byId, nil
}

func (f *mockControllerFaturaRepo) Delete(id string) error {
	return nil
}

type mockControllerBillRepo struct {
	created   []models.Bill
	list      []models.Bill
	parcelado []models.Bill
	fixo      []models.Bill
	filtered  []models.Bill
}

func (f *mockControllerBillRepo) Create(bill models.Bill, faturaID string) error {
	f.created = append(f.created, bill)
	return nil
}

func (f *mockControllerBillRepo) GetAllByFaturaId(faturaID string) ([]models.Bill, error) {
	return f.list, nil
}

func (f *mockControllerBillRepo) Delete(id string) error {
	return nil
}

func (f *mockControllerBillRepo) GetParcelado(id string) ([]models.Bill, error) {
	return f.parcelado, nil
}

func (f *mockControllerBillRepo) GetFixo(id string) ([]models.Bill, error) {
	return f.fixo, nil
}

func (f *mockControllerBillRepo) GetBillsByCategory(fatura_id string, category string) ([]models.Bill, error) {
	return f.filtered, nil
}

func makeMultipartCreateRequest(t *testing.T, payload dto.RequestFatura, csvContent string) *http.Request {
	t.Helper()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	dataJSON, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("erro ao serializar payload: %v", err)
	}

	if err = writer.WriteField("data", string(dataJSON)); err != nil {
		t.Fatalf("erro ao escrever campo data: %v", err)
	}

	part, err := writer.CreateFormFile("csv", "fatura.csv")
	if err != nil {
		t.Fatalf("erro ao criar arquivo no form: %v", err)
	}

	if _, err = part.Write([]byte(csvContent)); err != nil {
		t.Fatalf("erro ao escrever csv: %v", err)
	}

	if err = writer.Close(); err != nil {
		t.Fatalf("erro ao fechar multipart: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/fatura", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}

func TestCreateFaturaController(t *testing.T) {
	faturaRepo := &mockControllerFaturaRepo{}
	billRepo := &mockControllerBillRepo{}
	originalService := getFaturaService
	getFaturaService = func() *service.FaturaService {
		return service.NewFaturaServiceWithDependencies(faturaRepo, billRepo)
	}
	defer func() { getFaturaService = originalService }()

	req := makeMultipartCreateRequest(t, dto.RequestFatura{Description: "Fatura teste", Bank: "nubank", Status: "paid"}, "date,title,amount\n2024-01-01,Compra,50.00\n")
	res := httptest.NewRecorder()

	CreateFatura(res, req)

	if res.Code != http.StatusCreated {
		t.Fatalf("status inesperado: %d", res.Code)
	}

	if len(faturaRepo.created) != 1 {
		t.Fatalf("quantidade de faturas esperada: 1, obtida: %d", len(faturaRepo.created))
	}

}

func TestShowFaturasController(t *testing.T) {
	faturaRepo := &mockControllerFaturaRepo{list: []models.Fatura{{Id: "1", Description: "Fatura teste", Status: "paid", Total: 50.0}}}
	billRepo := &mockControllerBillRepo{}
	originalService := getFaturaService
	getFaturaService = func() *service.FaturaService {
		return service.NewFaturaServiceWithDependencies(faturaRepo, billRepo)
	}
	defer func() { getFaturaService = originalService }()

	req := httptest.NewRequest(http.MethodGet, "/fatura", nil)
	res := httptest.NewRecorder()

	ShowFaturas(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status inesperado: %d", res.Code)
	}

}

func TestGetFaturaController(t *testing.T) {
	faturaRepo := &mockControllerFaturaRepo{byId: models.Fatura{Id: "1", Description: "Fatura teste", Status: "paid", Total: 50.0}}
	billRepo := &mockControllerBillRepo{list: []models.Bill{{Id: "b1", Title: "Compra", Amount: 50.0}}}
	originalService := getFaturaService
	getFaturaService = func() *service.FaturaService {
		return service.NewFaturaServiceWithDependencies(faturaRepo, billRepo)
	}
	defer func() { getFaturaService = originalService }()

	req := httptest.NewRequest(http.MethodGet, "/fatura/1", nil)
	req.SetPathValue("id", "1")
	res := httptest.NewRecorder()

	GetFatura(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status inesperado: %d", res.Code)
	}

}

func TestDeleteFaturaController(t *testing.T) {
	faturaRepo := &mockControllerFaturaRepo{}
	billRepo := &mockControllerBillRepo{}
	originalService := getFaturaService
	getFaturaService = func() *service.FaturaService {
		return service.NewFaturaServiceWithDependencies(faturaRepo, billRepo)
	}
	defer func() { getFaturaService = originalService }()

	req := httptest.NewRequest(http.MethodDelete, "/fatura/1", nil)
	req.SetPathValue("id", "1")
	res := httptest.NewRecorder()

	DeleteFatura(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status inesperado: %d", res.Code)
	}

	if body := res.Body.String(); !bytes.Contains([]byte(body), []byte("Fatura deletada com sucesso")) {
		t.Fatalf("resposta inesperada: %s", body)
	}
}

func TestGetFaturaParceladoController(t *testing.T) {
	faturaRepo := &mockControllerFaturaRepo{}
	billRepo := &mockControllerBillRepo{
		parcelado: []models.Bill{{Id: "b1", Title: "Compra parcelada", Amount: 50.0}},
	}
	originalService := getFaturaService
	getFaturaService = func() *service.FaturaService {
		return service.NewFaturaServiceWithDependencies(faturaRepo, billRepo)
	}
	defer func() { getFaturaService = originalService }()

	req := httptest.NewRequest(http.MethodGet, "/fatura/1/parcelado", nil)
	req.SetPathValue("id", "1")
	res := httptest.NewRecorder()

	GetFaturaParcelado(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status inesperado: %d", res.Code)
	}

	if body := res.Body.String(); !bytes.Contains([]byte(body), []byte("Compra parcelada")) {
		t.Fatalf("resposta inesperada: %s", body)
	}
}

func TestGetFaturaFixoController(t *testing.T) {
	faturaRepo := &mockControllerFaturaRepo{}
	billRepo := &mockControllerBillRepo{
		fixo: []models.Bill{{Id: "b2", Title: "Conta fixa", Amount: 75.0}},
	}
	originalService := getFaturaService
	getFaturaService = func() *service.FaturaService {
		return service.NewFaturaServiceWithDependencies(faturaRepo, billRepo)
	}
	defer func() { getFaturaService = originalService }()

	req := httptest.NewRequest(http.MethodGet, "/fatura/1/fixo", nil)
	req.SetPathValue("id", "1")
	res := httptest.NewRecorder()

	GetFaturaFixo(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status inesperado: %d", res.Code)
	}

	if body := res.Body.String(); !bytes.Contains([]byte(body), []byte("Conta fixa")) {
		t.Fatalf("resposta inesperada: %s", body)
	}
}
