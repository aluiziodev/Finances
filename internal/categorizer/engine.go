package categorizer

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var (
	compiledCache = make(map[string]map[string][]*regexp.Regexp)
	cacheMu       sync.RWMutex
)

func ClassifyBillTitle(title string, bank string) (string, error) {
	compiled, err := loadCompiledCategories(bank)
	if err != nil {
		return "", err
	}

	original := strings.ToLower(strings.TrimSpace(title))
	normalized := normalize(title)

	if original == "" && normalized == "" {
		return "outros", nil
	}

	for categoria, regexs := range compiled {
		for _, re := range regexs {
			if re == nil {
				continue
			}

			if original != "" && re.MatchString(original) {
				return categoria, nil
			}
			if normalized != "" && re.MatchString(normalized) {
				return categoria, nil
			}
		}
	}

	return "outros", nil
}

func loadCompiledCategories(bank string) (map[string][]*regexp.Regexp, error) {
	cacheMu.RLock()
	if m, ok := compiledCache[bank]; ok {
		cacheMu.RUnlock()
		return m, nil
	}
	cacheMu.RUnlock()

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

	compiled := make(map[string][]*regexp.Regexp, len(categories))
	for categoria, padroes := range categories {
		regs := make([]*regexp.Regexp, 0, len(padroes))
		for _, padrao := range padroes {
			re, err := regexp.Compile("(?i)" + padrao)
			if err != nil {
				regs = append(regs, nil)
				continue
			}
			regs = append(regs, re)
		}
		compiled[categoria] = regs
	}

	cacheMu.Lock()
	compiledCache[bank] = compiled
	cacheMu.Unlock()

	return compiled, nil
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

func normalize(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	t := transform.Chain(norm.NFD)
	res, _, err := transform.String(t, s)
	if err != nil {
		res = s
	}

	var b []rune
	for _, r := range res {
		if unicode.Is(unicode.Mn, r) {
			continue
		}
		if unicode.IsPunct(r) || unicode.IsSymbol(r) {
			continue
		}
		if r == '\u00A0' {
			r = ' '
		}
		b = append(b, r)
	}
	out := strings.Join(strings.Fields(string(b)), " ")
	return out
}
