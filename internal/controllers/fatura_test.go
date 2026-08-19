package controllers

import (
	"bytes"
	"encoding/json"
	"finances/internal/models"
	"finances/internal/service"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

type fakeControllerFaturaRepo struct {
	created []models.Fatura
	list    []models.Fatura
	byId    models.Fatura
}

func (f *fakeControllerFaturaRepo) Create(fatura models.Fatura) error {
	f.created = append(f.created, fatura)
	return nil
}

func (f *fakeControllerFaturaRepo) GetAll() ([]models.Fatura, error) {
	return f.list, nil
}

func (f *fakeControllerFaturaRepo) Get(id string) (models.Fatura, error) {
	return f.byId, nil
}

type fakeControllerBillRepo struct {
	created []models.Bill
	list    []models.Bill
}

func (f *fakeControllerBillRepo) Create(bill models.Bill, faturaID string) error {
	f.created = append(f.created, bill)
	return nil
}

func (f *fakeControllerBillRepo) GetAllByFaturaId(faturaID string) ([]models.Bill, error) {
	return f.list, nil
}

func makeMultipartCreateRequest(t *testing.T, payload models.RequestFatura, csvContent string) *http.Request {
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
	faturaRepo := &fakeControllerFaturaRepo{}
	billRepo := &fakeControllerBillRepo{}
	originalService := getFaturaService
	getFaturaService = func() *service.FaturaService {
		return service.NewFaturaServiceWithDependencies(faturaRepo, billRepo)
	}

	req := makeMultipartCreateRequest(t, models.RequestFatura{Description: "Fatura teste", Status: "paid"}, "date,title,amount\n2024-01-01,Compra,50.00\n")
	res := httptest.NewRecorder()

	CreateFatura(res, req)

	if res.Code != http.StatusCreated {
		t.Fatalf("status inesperado: %d", res.Code)
	}

	if len(faturaRepo.created) != 1 {
		t.Fatalf("quantidade de faturas esperada: 1, obtida: %d", len(faturaRepo.created))
	}

	getFaturaService = originalService
}

func TestShowFaturasController(t *testing.T) {
	faturaRepo := &fakeControllerFaturaRepo{list: []models.Fatura{{Id: "1", Description: "Fatura teste", Status: "paid", Total: 50.0}}}
	billRepo := &fakeControllerBillRepo{}
	originalService := getFaturaService
	getFaturaService = func() *service.FaturaService {
		return service.NewFaturaServiceWithDependencies(faturaRepo, billRepo)
	}

	req := httptest.NewRequest(http.MethodGet, "/fatura", nil)
	res := httptest.NewRecorder()

	ShowFaturas(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status inesperado: %d", res.Code)
	}

	getFaturaService = originalService
}

func TestGetFaturaController(t *testing.T) {
	faturaRepo := &fakeControllerFaturaRepo{byId: models.Fatura{Id: "1", Description: "Fatura teste", Status: "paid", Total: 50.0}}
	billRepo := &fakeControllerBillRepo{list: []models.Bill{{Id: "b1", Title: "Compra", Amount: 50.0}}}
	originalService := getFaturaService
	getFaturaService = func() *service.FaturaService {
		return service.NewFaturaServiceWithDependencies(faturaRepo, billRepo)
	}

	req := httptest.NewRequest(http.MethodGet, "/fatura/1", nil)
	req.SetPathValue("id", "1")
	res := httptest.NewRecorder()

	GetFatura(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status inesperado: %d", res.Code)
	}

	getFaturaService = originalService
}
