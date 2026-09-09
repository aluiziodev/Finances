package service

import (
	"finances/internal/dto"
	"finances/internal/repository"

	"github.com/google/uuid"
)

type CardService struct {
	cardRepo repository.CardRepositoryInterface
}

func NewCardServiceWithDependencies(cardRepo repository.CardRepositoryInterface) *CardService {
	return &CardService{
		cardRepo: cardRepo,
	}
}

func NewCardService() *CardService {
	return NewCardServiceWithDependencies(repository.NewCardRepository())
}

func (s *CardService) CreateCard(card dto.RequestCard) (*string, error) {
	model := card.ToModel()
	model.Id = uuid.NewString()
	if err := model.Validate(); err != nil {
		return nil, err
	}
	if err := s.cardRepo.Create(model); err != nil {
		return nil, err
	}
	return &model.Id, nil
}

func (s *CardService) GetAllCards() (*[]dto.ResponseCard, error) {
	cards, err := s.cardRepo.GetAll()
	if err != nil {
		return nil, err
	}

	responseCards := dto.CardsToResponseCards(&cards)
	return &responseCards, nil
}

func (s *CardService) GetCardById(id string) (*dto.ResponseCard, error) {
	card, err := s.cardRepo.Get(id)
	if err != nil {
		return nil, err
	}

	responseCard := dto.NewResponseCard(&card)
	return &responseCard, nil
}

func (s *CardService) DeleteCard(id string) error {
	return s.cardRepo.Delete(id)
}
