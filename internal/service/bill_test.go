package service

import (
	"finances/internal/models"
	"testing"
)

func TestBillService_GetFaturaParcelado(t *testing.T) {
	billRepo := &mockBillRepo{
		parcelado: []models.Bill{{Id: "bill-1", Title: "Uber - Parcela 1", Amount: 80.0, Method: "parcelado"}, {Id: "bill-2", Title: "Uber - Parcela 2", Amount: 20.0, Method: "parcelado"}},
	}
	service := NewBillServiceWithDependencies(billRepo)

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

func TestBillService_GetFaturaFixo(t *testing.T) {
	billRepo := &mockBillRepo{
		fixo: []models.Bill{{Id: "bill-1", Title: "Netflix", Amount: 30.0, Method: "fixo"}, {Id: "bill-2", Title: "Mercado", Amount: 70.0, Method: "fixo"}},
	}
	service := NewBillServiceWithDependencies(billRepo)

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
