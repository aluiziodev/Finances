package service

import (
	"finances/internal/categorizer"
	"finances/internal/dto"
	"finances/internal/parser"
	"finances/internal/repository"
	"mime/multipart"

	"github.com/google/uuid"
)

type FaturaService struct {
	cardRepo   repository.CardRepositoryInterface
	faturaRepo repository.FaturaRepositoryInterface
	billRepo   repository.BillRepositoryInterface
}

func NewFaturaServiceWithDependencies(cardRepo repository.CardRepositoryInterface, faturaRepo repository.FaturaRepositoryInterface, billRepo repository.BillRepositoryInterface) *FaturaService {
	return &FaturaService{
		cardRepo:   cardRepo,
		faturaRepo: faturaRepo,
		billRepo:   billRepo,
	}
}

func NewFaturaService() *FaturaService {
	return NewFaturaServiceWithDependencies(repository.NewCardRepository(), repository.NewFaturaRepository(), repository.NewBillRepository())
}

func (s *FaturaService) CreateFatura(file multipart.File, handler *multipart.FileHeader,
	req dto.RequestFatura, card_id string) (*string, error) {

	bank, err := s.cardRepo.GetBankByCardId(card_id)
	if err != nil {
		return nil, err
	}
	fatura, err := parser.ParserCSVtoModels(file, handler, req, bank)
	if err != nil {
		return nil, err
	}

	fatura.CardID = card_id

	if err := fatura.Validate(); err != nil {
		return nil, err
	}
	fatura.Id = uuid.NewString()
	fatura.CalculateTotal()

	if err := s.faturaRepo.Create(fatura); err != nil {
		return nil, err
	}

	for _, bill := range fatura.Bills {

		bill.Id = uuid.NewString()
		bill.FaturaID = fatura.Id
		bill.DefineMethod()
		bill.Category, err = categorizer.ClassifyBillTitle(bill.Title, "nubank")
		if err != nil {
			return nil, err
		}
		if err := bill.Validate(); err != nil {
			return nil, err
		}
		if err := s.billRepo.Create(bill, fatura.Id); err != nil {
			return nil, err
		}
	}

	return &fatura.Id, nil
}

func (s *FaturaService) GetAllFaturas(card_id string) (*[]dto.ResponseFatura, error) {
	faturas, err := s.faturaRepo.GetAll(card_id)
	if err != nil {
		return nil, err
	}

	responseFaturas := dto.FaturasToResponseFaturas(&faturas)

	return &responseFaturas, nil
}

func (s *FaturaService) GetFatura(id string) (*dto.Summary, error) {
	fatura, err := s.faturaRepo.Get(id)
	if err != nil {
		return nil, err
	}

	bills, err := s.billRepo.GetAllByFaturaId(id)
	if err != nil {
		return nil, err
	}

	fatura.Bills = bills

	summary := dto.Summary{}
	summary.CalculateTotal(&fatura)

	return &summary, nil
}

func (s *FaturaService) DeleteFatura(id string) error {

	if err := s.faturaRepo.Delete(id); err != nil {
		return err
	}

	return nil
}

func (s *FaturaService) UpdateFaturaPaid(id string) error {
	if err := s.faturaRepo.UpdateStatusPaid(id); err != nil {
		return err
	}

	return nil
}

func (s *FaturaService) UpdateFaturaPending(id string) error {
	if err := s.faturaRepo.UpdateStatusPending(id); err != nil {
		return err
	}

	return nil
}
