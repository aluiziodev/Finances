package categorizer

import (
	"strings"
	"testing"
)

func TestClassifyBillTitle_MatchesKnownPattern(t *testing.T) {
	tests := []struct {
		name     string
		title    string
		bank     string
		expected string
	}{
		{name: "uber", title: "Uber *Trip", bank: "nubank", expected: "transporte"},
		{name: "ifood", title: "IFood - Pedido", bank: "nubank", expected: "alimentação"},
		{name: "spotify", title: "Spotify Premium", bank: "nubank", expected: "assinaturas"},
		{name: "farmacia", title: "Drogaria São Paulo", bank: "nubank", expected: "saúde"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

func TestClassifyBillTitle_EmptyTitleReturnsOthers(t *testing.T) {
	category, err := ClassifyBillTitle("   ", "nubank")
	if err != nil {
		t.Fatalf("esperava sem erro, mas recebeu: %v", err)
	}

	if category != "outros" {
		t.Fatalf("categoria inesperada: recebeu %q, esperava %q", category, "outros")
	}
}

func TestClassifyBillTitle_NoMatchReturnsOthers(t *testing.T) {
	category, err := ClassifyBillTitle("Pagamento de boleto sem categoria", "nubank")
	if err != nil {
		t.Fatalf("esperava sem erro, mas recebeu: %v", err)
	}

	if category != "outros" {
		t.Fatalf("categoria inesperada: recebeu %q, esperava %q", category, "outros")
	}
}

func TestClassifyBillTitle_InvalidBankReturnsError(t *testing.T) {
	_, err := ClassifyBillTitle("Spotify Premium", "banco_inexistente")
	if err == nil {
		t.Fatal("esperava erro para banco inexistente")
	}

	if !strings.Contains(err.Error(), "não foi possível ler categorias") {
		t.Fatalf("mensagem de erro inesperada: %v", err)
	}
}
