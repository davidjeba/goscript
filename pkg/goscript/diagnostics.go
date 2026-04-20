package goscript

import "strings"

// Diagnostic describes a source or manifest issue.
type Diagnostic struct {
	File       string `json:"file,omitempty"`
	Code       string `json:"code,omitempty"`
	Severity   string `json:"severity,omitempty"`
	Message    string `json:"message"`
	Suggestion string `json:"suggestion,omitempty"`
	Line       int    `json:"line,omitempty"`
	Column     int    `json:"column,omitempty"`
}

// AnalyzeSource produces lightweight diagnostics for a GoScript source file.
func AnalyzeSource(fileName, source string) []Diagnostic {
	diagnostics := make([]Diagnostic, 0)

	if strings.Contains(source, "TODO") {
		diagnostics = append(diagnostics, Diagnostic{
			File:     fileName,
			Code:     "todo-comment",
			Severity: "info",
			Message:  "source contains TODO markers",
		})
	}

	if strings.Contains(source, "not fully implemented") {
		diagnostics = append(diagnostics, Diagnostic{
			File:       fileName,
			Code:       "stubbed-implementation",
			Severity:   "warning",
			Message:    "source contains a stubbed implementation",
			Suggestion: "replace the stub with a complete implementation or mark it as experimental",
		})
	}

	if strings.Contains(source, "panic(") {
		diagnostics = append(diagnostics, Diagnostic{
			File:       fileName,
			Code:       "panic-path",
			Severity:   "warning",
			Message:    "source includes panic paths",
			Suggestion: "prefer Result or Diagnostic returns for recoverable errors",
		})
	}

	return diagnostics
}

