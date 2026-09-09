package controllers

import (
	"bytes"
	"encoding/json"
	"finances/internal/dto"
	"finances/internal/models"
	"finances/internal/service"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockControllerCardRepo struct {
	created []models.Card
	list    []models.Card
	byId    models.Card
	deleted []string
	bank    string
}

func (m *mockControllerCardRepo) Create(card models.Card) error {
	m.created = append(m.created, card)
	return nil
}

func (m *mockControllerCardRepo) GetAll() ([]models.Card, error) {
	return m.list, nil
}

func (m *mockControllerCardRepo) Get(id string) (models.Card, error) {
	return m.byId, nil
}

func (m *mockControllerCardRepo) GetBankByCardId(cardId string) (string, error) {
	return m.bank, nil
}

func (m *mockControllerCardRepo) Delete(id string) error {
	m.deleted = append(m.deleted, id)
	return nil
}

func (m *mockControllerCardRepo) Update(card models.Card) error { return nil }

func TestCreateCardController(t *testing.T) {
	repo := &mockControllerCardRepo{}
	originalService := getCardService
	getCardService = func() *service.CardService {
		return service.NewCardServiceWithDependencies(repo)
	}
	defer func() { getCardService = originalService }()

	payload := dto.RequestCard{Name: "Nubank", Bank: "nubank", Limit: 5000, DueDay: 10}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/card", bytes.NewReader(body))
	res := httptest.NewRecorder()

	CreateCard(res, req)

	if res.Code != http.StatusCreated {
		t.Fatalf("status inesperado: %d", res.Code)
	}

	if len(repo.created) != 1 {
		t.Fatalf("quantidade de cartoes esperada: 1, obtida: %d", len(repo.created))
	}
}

func TestShowCardsController(t *testing.T) {
	repo := &mockControllerCardRepo{list: []models.Card{{Id: "c1", Name: "Nubank", Bank: "nubank", Limit: 5000, DueDay: 10}}}
	originalService := getCardService
	getCardService = func() *service.CardService {
		return service.NewCardServiceWithDependencies(repo)
	}
	defer func() { getCardService = originalService }()

	req := httptest.NewRequest(http.MethodGet, "/card", nil)
	res := httptest.NewRecorder()

	ShowCards(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status inesperado: %d", res.Code)
	}
}

func TestGetCardByIdController(t *testing.T) {
	repo := &mockControllerCardRepo{byId: models.Card{Id: "c1", Name: "Nubank", Bank: "nubank", Limit: 5000, DueDay: 10}}
	originalService := getCardService
	getCardService = func() *service.CardService {
		return service.NewCardServiceWithDependencies(repo)
	}
	defer func() { getCardService = originalService }()

	req := httptest.NewRequest(http.MethodGet, "/card/c1", nil)
	req.SetPathValue("id", "c1")
	res := httptest.NewRecorder()

	GetCardById(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status inesperado: %d", res.Code)
	}
}

func TestDeleteCardController(t *testing.T) {
	repo := &mockControllerCardRepo{}
	originalService := getCardService
	getCardService = func() *service.CardService {
		return service.NewCardServiceWithDependencies(repo)
	}
	defer func() { getCardService = originalService }()

	req := httptest.NewRequest(http.MethodDelete, "/card/c1", nil)
	req.SetPathValue("id", "c1")
	res := httptest.NewRecorder()

	DeleteCard(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status inesperado: %d", res.Code)
	}
}
