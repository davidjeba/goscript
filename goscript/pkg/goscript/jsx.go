package goscript

import (
	"fmt"
	"regexp"
	"strings"
)

// JSXParser parses JSX-like syntax into CreateElement calls
type JSXParser struct {
	// Configuration options
	options map[string]interface{}
}

// NewJSXParser creates a new JSX parser
func NewJSXParser(options map[string]interface{}) *JSXParser {
	return &JSXParser{
		options: options,
	}
}

// ParseJSX parses a JSX-like string into Go code
func (p *JSXParser) ParseJSX(jsx string) (string, error) {
	// This is a simplified implementation
	// A real implementation would use a proper parser
	
	// Replace self-closing tags
	selfClosingRegex := regexp.MustCompile(`<([a-zA-Z][a-zA-Z0-9]*)\s*([^>]*?)\s*/>`)
	jsx = selfClosingRegex.ReplaceAllString(jsx, `CreateElement("$1", Props{$2})`)
	
	// Replace opening tags with attributes
	openingTagRegex := regexp.MustCompile(`<([a-zA-Z][a-zA-Z0-9]*)\s*([^>]*?)>`)
	jsx = openingTagRegex.ReplaceAllString(jsx, `CreateElement("$1", Props{$2}, `)
	
	// Replace closing tags
	closingTagRegex := regexp.MustCompile(`</([a-zA-Z][a-zA-Z0-9]*)>`)
	jsx = closingTagRegex.ReplaceAllString(jsx, `)`)
	
	// Handle text nodes
	// This is a very simplified approach
	textNodeRegex := regexp.MustCompile(`>([^<]+)<`)
	jsx = textNodeRegex.ReplaceAllString(jsx, `, "$1"<`)
	
	return jsx, nil
}

// ParseAttributes parses JSX attributes into Props
func (p *JSXParser) ParseAttributes(attrs string) (Props, error) {
	props := Props{}
	
	if attrs == "" {
		return props, nil
	}
	
	// This is a simplified implementation
	// A real implementation would handle quotes, expressions, etc.
	attrRegex := regexp.MustCompile(`([a-zA-Z][a-zA-Z0-9]*)\s*=\s*"([^"]*)"`)
	matches := attrRegex.FindAllStringSubmatch(attrs, -1)
	
	for _, match := range matches {
		if len(match) == 3 {
			key := match[1]
			value := match[2]
			props[key] = value
		}
	}
	
	return props, nil
}

// TranspileGSX transpiles a .gsx file to a .go file
func TranspileGSX(gsxContent string) (string, error) {
	// This is a simplified implementation
	// A real implementation would use a proper parser
	
	// Extract imports
	importsRegex := regexp.MustCompile(`import\s+\(([^)]+)\)`)
	importsMatch := importsRegex.FindStringSubmatch(gsxContent)
	imports := ""
	if len(importsMatch) > 1 {
		imports = importsMatch[1]
	}
	
	// Extract package
	packageRegex := regexp.MustCompile(`package\s+([a-zA-Z][a-zA-Z0-9]*)`)
	packageMatch := packageRegex.FindStringSubmatch(gsxContent)
	packageName := ""
	if len(packageMatch) > 1 {
		packageName = packageMatch[1]
	}
	
	// Extract function declarations
	funcRegex := regexp.MustCompile(`func\s+([a-zA-Z][a-zA-Z0-9]*)\s*\(([^)]*)\)\s*string\s*{([^}]+)}`)
	funcMatches := funcRegex.FindAllStringSubmatch(gsxContent, -1)
	
	var functions []string
	for _, match := range funcMatches {
		if len(match) > 3 {
			funcName := match[1]
			funcParams := match[2]
			funcBody := match[3]
			
			// Parse JSX in function body
			parser := NewJSXParser(nil)
			parsedBody, err := parser.ParseJSX(funcBody)
			if err != nil {
				return "", fmt.Errorf("error parsing JSX in function %s: %v", funcName, err)
			}
			
			function := fmt.Sprintf("func %s(%s) string {\n%s\n}", funcName, funcParams, parsedBody)
			functions = append(functions, function)
		}
	}
	
	// Combine everything
	result := fmt.Sprintf("package %s\n\nimport (%s)\n\n%s\n", 
		packageName, imports, strings.Join(functions, "\n\n"))
	
	return result, nil
}

// GSXCompiler compiles .gsx files to .go files
type GSXCompiler struct {
	// Configuration options
	options map[string]interface{}
}

// NewGSXCompiler creates a new GSX compiler
func NewGSXCompiler(options map[string]interface{}) *GSXCompiler {
	return &GSXCompiler{
		options: options,
	}
}

// CompileFile compiles a .gsx file to a .go file
func (c *GSXCompiler) CompileFile(inputPath, outputPath string) error {
	// In a real implementation, this would:
	// 1. Read the input file
	// 2. Parse the GSX content
	// 3. Generate Go code
	// 4. Write to the output file
	
	// For now, we'll just return a placeholder error
	return fmt.Errorf("GSX compilation not fully implemented")
}