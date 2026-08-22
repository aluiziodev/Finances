package service

import (
	"finances/internal/categorizer"
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

	fatura, err := parser.ParserCSVtoModels(file, handler, req)
	if err != nil {
		return nil, err
	}

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
		bill.DefineMethod()
		bill.Category, err = categorizer.ClassifyBillTitle(bill.Title, fatura.Bank)
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

	return &fatura, nil
}

func (s *FaturaService) GetAllFaturas() (*[]models.Fatura, error) {
	faturas, err := s.faturaRepo.GetAll()
	if err != nil {
		return nil, err
	}

	return &faturas, nil
}

func (s *FaturaService) GetFatura(id string) (*models.Fatura, error) {
	fatura, err := s.faturaRepo.Get(id)
	if err != nil {
		return nil, err
	}

	bills, err := s.billRepo.GetAllByFaturaId(id)
	if err != nil {
		return nil, err
	}

	fatura.Bills = bills

	return &fatura, nil
}

func (s *FaturaService) DeleteFatura(id string) error {

	if err := s.faturaRepo.Delete(id); err != nil {
		return err
	}

	return nil
}

func (s *FaturaService) GetFaturaParcelado(id string) (*models.Fatura, error) {
	bills, err := s.billRepo.GetParcelado(id)
	if err != nil {
		return nil, err
	}

	fatura_parcelado := models.Fatura{
		Id:    id,
		Bills: bills,
	}
	fatura_parcelado.CalculateTotal()
	return &fatura_parcelado, nil
}

func (s *FaturaService) GetFaturaFixo(id string) (*models.Fatura, error) {
	bills, err := s.billRepo.GetFixo(id)
	if err != nil {
		return nil, err
	}

	fatura_fixo := models.Fatura{
		Id:    id,
		Bills: bills,
	}
	fatura_fixo.CalculateTotal()
	return &fatura_fixo, nil
}

func (s *FaturaService) CalculateSummary(fatura *models.Fatura) *models.Summary {
	summary := models.Summary{}
	summary.CalculateTotal(fatura)

	return &summary
}
