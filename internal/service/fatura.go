package service

import (
	"finances/internal/models"
	"finances/internal/parser"
	"finances/internal/repository"
	"mime/multipart"

	"github.com/google/uuid"
)

type FaturaService struct {
	faturaRepo repository.FaturaRepositoryInterface
	billRepo   repository.BillRepositoryInterface
}

func NewFaturaServiceWithDependencies(faturaRepo repository.FaturaRepositoryInterface, billRepo repository.BillRepositoryInterface) *FaturaService {
	return &FaturaService{
		faturaRepo: faturaRepo,
		billRepo:   billRepo,
	}
}

func NewFaturaService() *FaturaService {
	return NewFaturaServiceWithDependencies(repository.NewFaturaRepository(), repository.NewBillRepository())
}

func (s *FaturaService) CreateFatura(file multipart.File, handler *multipart.FileHeader,
	req models.RequestFatura) (*models.Fatura, error) {

	fatura, err := parser.ParserCSVtoModels(file, handler, req.Description, req.Status)
	if err != nil {
		return nil, err
	}
	fatura.Id = uuid.NewString()

	if err := s.faturaRepo.Create(fatura); err != nil {
		return nil, err
	}

	for _, bill := range fatura.Bills {
		bill.Id = uuid.NewString()
		if err := s.billRepo.Create(bill, fatura.Id); err != nil {
			return nil, err
		}
	}

	return &fatura, nil
}

func (s *FaturaService) GetAllFaturas() ([]models.Fatura, error) {
	faturas, err := s.faturaRepo.GetAll()
	if err != nil {
		return nil, err
	}

	return faturas, nil
}

func (s *FaturaService) GetFatura(id string) (models.Fatura, error) {
	fatura, err := s.faturaRepo.Get(id)
	if err != nil {
		return models.Fatura{}, err
	}

	bills, err := s.billRepo.GetAllByFaturaId(id)
	if err != nil {
		return models.Fatura{}, err
	}

	fatura.Bills = bills

	return fatura, nil
}

func CreateFatura(file multipart.File, handler *multipart.FileHeader,
	req models.RequestFatura) (*models.Fatura, error) {
	return NewFaturaService().CreateFatura(file, handler, req)
}

func GetAllFaturas() ([]models.Fatura, error) {
	return NewFaturaService().GetAllFaturas()
}

func GetFatura(id string) (models.Fatura, error) {
	return NewFaturaService().GetFatura(id)
}
