package service

import (
	"finances/internal/dto"
	"finances/internal/models"
	"finances/internal/repository"
)

type BillService struct {
	billRepo repository.BillRepositoryInterface
}

func NewBillServiceWithDependencies(billRepo repository.BillRepositoryInterface) *BillService {
	return &BillService{
		billRepo: billRepo,
	}
}

func NewBillService() *BillService {
	return NewBillServiceWithDependencies(repository.NewBillRepository())
}

func (s *BillService) GetBillsByCategory(fatura_id string, category string) (*dto.ResponseFatura, error) {
	bills, err := s.billRepo.GetBillsByCategory(fatura_id, category)
	if err != nil {
		return nil, err
	}

	fatura_category := models.Fatura{
		Id:    fatura_id,
		Bills: bills,
	}
	fatura_category.CalculateTotal()

	responseFatura := dto.NewResponseFatura(&fatura_category)
	return &responseFatura, nil
}
