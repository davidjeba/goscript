package goscript

import (
	"regexp"
	"sort"
	"strings"
	"sync"
)

// Symbol describes an indexed language symbol.
type Symbol struct {
	Name     string            `json:"name"`
	Kind     string            `json:"kind"`
	File     string            `json:"file,omitempty"`
	Line     int               `json:"line,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// DocumentIndex stores symbols for completion and navigation.
type DocumentIndex struct {
	mu      sync.RWMutex
	symbols map[string]Symbol
}

// NewDocumentIndex creates an empty index.
func NewDocumentIndex() *DocumentIndex {
	return &DocumentIndex{
		symbols: make(map[string]Symbol),
	}
}

// IndexSource scans source and records simple symbols.
func (i *DocumentIndex) IndexSource(fileName, source string) []Symbol {
	lines := strings.Split(source, "\n")
	results := make([]Symbol, 0)

	funcRe := regexp.MustCompile(`^\s*func\s+([A-Za-z_][A-Za-z0-9_]*)`)
	typeRe := regexp.MustCompile(`^\s*type\s+([A-Za-z_][A-Za-z0-9_]*)`)
	varRe := regexp.MustCompile(`^\s*var\s+([A-Za-z_][A-Za-z0-9_]*)`)
	constRe := regexp.MustCompile(`^\s*const\s+([A-Za-z_][A-Za-z0-9_]*)`)

	for idx, line := range lines {
		if match := funcRe.FindStringSubmatch(line); len(match) == 2 {
			results = append(results, Symbol{Name: match[1], Kind: "func", File: fileName, Line: idx + 1})
		}
		if match := typeRe.FindStringSubmatch(line); len(match) == 2 {
			results = append(results, Symbol{Name: match[1], Kind: "type", File: fileName, Line: idx + 1})
		}
		if match := varRe.FindStringSubmatch(line); len(match) == 2 {
			results = append(results, Symbol{Name: match[1], Kind: "var", File: fileName, Line: idx + 1})
		}
		if match := constRe.FindStringSubmatch(line); len(match) == 2 {
			results = append(results, Symbol{Name: match[1], Kind: "const", File: fileName, Line: idx + 1})
		}
	}

	i.mu.Lock()
	for _, symbol := range results {
		i.symbols[symbol.Name] = symbol
	}
	i.mu.Unlock()

	return results
}

// Lookup returns a symbol by name.
func (i *DocumentIndex) Lookup(name string) (Symbol, bool) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	symbol, ok := i.symbols[name]
	return symbol, ok
}

// Complete returns symbol names that share a prefix.
func (i *DocumentIndex) Complete(prefix string) []string {
	i.mu.RLock()
	defer i.mu.RUnlock()

	results := make([]string, 0)
	for name := range i.symbols {
		if strings.HasPrefix(name, prefix) {
			results = append(results, name)
		}
	}

	sort.Strings(results)
	return results
}

