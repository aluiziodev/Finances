package categorizer

import (
	"strings"
	"testing"
)

func TestClassifyBillTitle_TableDriven(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		title    string
		bank     string
		expected string
	}{
		{name: "uber", title: "Uber *Trip", bank: "nubank", expected: "transporte"},
		{name: "uber_lower", title: "uber *trip", bank: "nubank", expected: "transporte"},
		{name: "ifood", title: "IFood - Pedido", bank: "nubank", expected: "alimentação"},
		{name: "ifood_accent", title: "ÍFood - Pedido", bank: "nubank", expected: "alimentação"},
		{name: "spotify", title: "Spotify Premium", bank: "nubank", expected: "assinaturas"},
		{name: "farmacia", title: "Drogaria São Paulo", bank: "nubank", expected: "saúde"},
		{name: "empty", title: "   ", bank: "nubank", expected: "outros"},
		{name: "no_match", title: "Pagamento de boleto sem categoria", bank: "nubank", expected: "outros"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			category, err := ClassifyBillTitle(tt.title, tt.bank)
			if err != nil {
				t.Fatalf("esperava sem erro, mas recebeu: %v", err)
			}

			if category != tt.expected {
				t.Fatalf("categoria inesperada: recebeu %q, esperava %q", category, tt.expected)
			}
		})
	}
}

func TestClassifyBillTitle_InvalidBankReturnsError(t *testing.T) {
	t.Parallel()

	_, err := ClassifyBillTitle("Spotify Premium", "banco_inexistente")
	if err == nil {
		t.Fatal("esperava erro para banco inexistente")
	}

	if !strings.Contains(err.Error(), "não foi possível ler categorias") {
		t.Fatalf("mensagem de erro inesperada: %v", err)
	}
}
