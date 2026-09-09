package service

import (
	"finances/internal/dto"
	"finances/internal/models"
	"testing"
)

type mockCardServiceRepo struct {
	created []models.Card
	list    []models.Card
	byId    models.Card
	deleted []string
}

func (m *mockCardServiceRepo) Create(card models.Card) error {
	m.created = append(m.created, card)
	return nil
}

func (m *mockCardServiceRepo) GetAll() ([]models.Card, error) {
	return m.list, nil
}

func (m *mockCardServiceRepo) Get(id string) (models.Card, error) {
	return m.byId, nil
}

func (m *mockCardServiceRepo) GetBankByCardId(cardId string) (string, error) {
	return "nubank", nil
}

func (m *mockCardServiceRepo) Delete(id string) error {
	m.deleted = append(m.deleted, id)
	return nil
}

func (m *mockCardServiceRepo) Update(card models.Card) error { return nil }

func TestCardService_CreateCard(t *testing.T) {
	repo := &mockCardServiceRepo{}
	service := NewCardServiceWithDependencies(repo)

	id, err := service.CreateCard(dto.RequestCard{
		Name:   "Nubank",
		Bank:   "nubank",
		Limit:  5000,
		DueDay: 10,
	})
	if err != nil {
		t.Fatalf("esperava criacao bem-sucedida, mas recebeu erro: %v", err)
	}

	if *id == "" {
		t.Fatal("id do cartao nao pode estar vazio")
	}

	if len(repo.created) != 1 {
		t.Fatalf("quantidade de cartoes criados inesperada: %d", len(repo.created))
	}
}

func TestCardService_GetAllCards(t *testing.T) {
	repo := &mockCardServiceRepo{list: []models.Card{{Id: "c1", Name: "Nubank", Bank: "nubank", Limit: 5000, DueDay: 10}}}
	service := NewCardServiceWithDependencies(repo)

	cards, err := service.GetAllCards()
	if err != nil {
		t.Fatalf("esperava listagem bem-sucedida, mas recebeu erro: %v", err)
	}

	if len(*cards) != 1 {
		t.Fatalf("quantidade de cartoes inesperada: %d", len(*cards))
	}
}

func TestCardService_GetCardById(t *testing.T) {
	repo := &mockCardServiceRepo{byId: models.Card{Id: "c1", Name: "Nubank", Bank: "nubank", Limit: 5000, DueDay: 10}}
	service := NewCardServiceWithDependencies(repo)

	card, err := service.GetCardById("c1")
	if err != nil {
		t.Fatalf("esperava busca bem-sucedida, mas recebeu erro: %v", err)
	}

	if card.Name != "Nubank" {
		t.Fatalf("nome do cartao inesperado: %s", card.Name)
	}
}

func TestCardService_DeleteCard(t *testing.T) {
	repo := &mockCardServiceRepo{}
	service := NewCardServiceWithDependencies(repo)

	if err := service.DeleteCard("c1"); err != nil {
		t.Fatalf("esperava exclusao bem-sucedida, mas recebeu erro: %v", err)
	}

	if len(repo.deleted) != 1 {
		t.Fatalf("quantidade de exclusoes inesperada: %d", len(repo.deleted))
	}
}
