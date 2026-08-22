package categorizer

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

func ClassifyBillTitle(title string, bank string) (string, error) {
	categorias, err := loadBillCategories(bank)
	if err != nil {
		return "", err
	}

	texto := strings.ToLower(strings.TrimSpace(title))
	if texto == "" {
		return "outros", nil
	}

	for categoria, padroes := range categorias {
		for _, padrao := range padroes {
			re, err := regexp.Compile("(?i)" + padrao)
			if err != nil {
				continue
			}
			if re.MatchString(texto) {
				return categoria, nil
			}
		}
	}

	return "outros", nil
}

func loadBillCategories(bank string) (map[string][]string, error) {
	_, filename, _, _ := runtime.Caller(0)
	baseDir := filepath.Dir(filename)
	path := filepath.Join(baseDir, "categorization", bank+"Categorization.json")

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("não foi possível ler categorias: %w", err)
	}

	var categories map[string][]string
	if err := json.Unmarshal(content, &categories); err != nil {
		return nil, fmt.Errorf("falha ao ler categorias: %w", err)
	}

	return categories, nil
}
