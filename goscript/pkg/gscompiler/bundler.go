package gscompiler

import (
        "fmt"
        "os"
        "path/filepath"
        "regexp"
        "sort"
        "strings"
)

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

// Bundle represents a compiled bundle of .gs files merged into a single
// JavaScript output.  It carries per-file metadata, the final concatenated
// JS source and an optional source map for debugging.
type Bundle struct {
        Files     []BundleFile      // Individual compiled files (in bundle order)
        Output    string            // Final concatenated JavaScript
        SourceMap map[string]string // .gs filename → original source (populated with WithSourceMap)
        Stats     BundleStats       // Summary statistics
}

// BundleFile represents a single compiled .gs file within a Bundle.
type BundleFile struct {
        Path       string // Original .gs file path (relative or absolute)
        ModuleName string // Short display name (basename without extension)
        JavaScript string // Compiled JavaScript output for this file
        Source     string // Original .gs source (populated with WithSourceMap)
        Imports    []string // Imports found in this file (goscript/* paths)
        SourceSize int    // Original .gs source size in bytes
        JSSize     int    // Compiled JavaScript size in bytes
}

// BundleStats holds summary statistics for a bundle.
type BundleStats struct {
        FileCount int // Total number of .gs files bundled
        GSSize    int // Total .gs source size in bytes
        JSSize    int // Total compiled JS size in bytes
        HasRuntime bool // Whether the runtime preamble is included
}

// BundleError represents a compilation error for a specific file.
type BundleError struct {
        File string // .gs file path
        Err  error  // Underlying error
}

func (e *BundleError) Error() string {
        if e.File != "" {
                return fmt.Sprintf("%s: %s", e.File, e.Err.Error())
        }
        return e.Err.Error()
}

func (e *BundleError) Unwrap() error {
        return e.Err
}

// ---------------------------------------------------------------------------
// Bundler
// ---------------------------------------------------------------------------

// Bundler manages the bundling process — discovering .gs files, compiling
// them to JavaScript, resolving imports, and producing a single output bundle.
type Bundler struct {
        EntryPoints    []string // Entry .gs files (sourced first in output)
        IncludeRuntime bool     // Whether to include __gs runtime preamble
        Minify         bool     // Whether to minify output
        SourceMap      bool     // Whether to collect source map data
        OutDir         string   // Output directory (if set, writes bundle.js)
}

// BundleOption is a functional option for configuring a Bundler.
type BundleOption func(*Bundler)

// WithRuntime prepends the goscript client runtime (embedded via go:embed)
// to the bundle output.  The runtime provides __gs.useState, __gs.h,
// __gs.mount, and all other client-side helpers.
func WithRuntime() BundleOption {
        return func(b *Bundler) {
                b.IncludeRuntime = true
        }
}

// WithMinify enables simple minification of the output JavaScript.
// It removes single-line comments, collapses whitespace, and strips
// leading indentation.  This is NOT a full minifier — for production
// use a dedicated tool like terser or uglify-js.
func WithMinify() BundleOption {
        return func(b *Bundler) {
                b.Minify = true
        }
}

// WithSourceMap collects the original .gs source for each file into
// Bundle.SourceMap.  The source map is a simple filename→source lookup
// (not a VLQ-encoded .map file).
func WithSourceMap() BundleOption {
        return func(b *Bundler) {
                b.SourceMap = true
        }
}

// WithEntryPoint sets one or more entry .gs files.  Entry files are
// compiled first and placed at the top of the bundle output (after the
// runtime preamble, if enabled).  Non-entry files follow in alphabetical
// order.
func WithEntryPoint(paths ...string) BundleOption {
        return func(b *Bundler) {
                b.EntryPoints = append(b.EntryPoints, paths...)
        }
}

// WithOutDir sets the output directory.  When set, the bundler writes
// the final JavaScript to <dir>/bundle.js after bundling.
func WithOutDir(dir string) BundleOption {
        return func(b *Bundler) {
                b.OutDir = dir
        }
}

// newBundler creates a Bundler with the given options applied.
func newBundler(opts ...BundleOption) *Bundler {
        b := &Bundler{}
        for _, opt := range opts {
                opt(b)
        }
        return b
}

// ---------------------------------------------------------------------------
// Public API
// ---------------------------------------------------------------------------

// BundleFromDir bundles all .gs files found recursively in the given
// directory.  Files are sorted alphabetically unless entry points are
// specified via WithEntryPoint, in which case entries come first.
//
// The returned Bundle contains the concatenated JavaScript, per-file
// metadata, and optional source map data.
func BundleFromDir(dir string, opts ...BundleOption) (*Bundle, error) {
        b := newBundler(opts...)

        // Discover .gs files recursively
        gsFiles, err := discoverGSFiles(dir)
        if err != nil {
                return nil, fmt.Errorf("discovering .gs files in %s: %w", dir, err)
        }

        if len(gsFiles) == 0 {
                return nil, fmt.Errorf("no .gs files found in %s", dir)
        }

        // Sort files, but put entry points first
        gsFiles = orderFiles(gsFiles, b.EntryPoints)

        return b.buildBundle(gsFiles)
}

// BundleFiles bundles specific .gs files given by their paths.
// Entry points (if any) are placed first in the output.
func BundleFiles(paths []string, opts ...BundleOption) (*Bundle, error) {
        if len(paths) == 0 {
                return nil, fmt.Errorf("no files specified for bundling")
        }

        b := newBundler(opts...)

        // Resolve to absolute paths and deduplicate
        absPaths := make([]string, 0, len(paths))
        seen := make(map[string]bool)
        for _, p := range paths {
                abs, err := filepath.Abs(p)
                if err != nil {
                        return nil, fmt.Errorf("resolving path %s: %w", p, err)
                }
                if !seen[abs] {
                        seen[abs] = true
                        absPaths = append(absPaths, abs)
                }
        }

        absPaths = orderFiles(absPaths, b.EntryPoints)

        return b.buildBundle(absPaths)
}

// ---------------------------------------------------------------------------
// Bundle construction
// ---------------------------------------------------------------------------

// buildBundle compiles the given .gs files and assembles the final bundle.
func (b *Bundler) buildBundle(gsFiles []string) (*Bundle, error) {
        bundle := &Bundle{
                Files:     make([]BundleFile, 0, len(gsFiles)),
                SourceMap: make(map[string]string),
        }

        // Phase 1: Compile each .gs file
        compileErrors := make([]BundleError, 0)
        for _, path := range gsFiles {
                bf, err := b.compileFile(path)
                if err != nil {
                        compileErrors = append(compileErrors, BundleError{
                                File: path,
                                Err:  err,
                        })
                        continue
                }
                bundle.Files = append(bundle.Files, *bf)
        }

        // If every file failed, return the first error
        if len(compileErrors) == len(gsFiles) {
                return nil, fmt.Errorf("all files failed to compile; first error: %w", compileErrors[0])
        }

        // Phase 2: Resolve imports (collect stdlib references across all files)
        allImports := b.resolveImports(bundle.Files)

        // Phase 3: Assemble output
        var sb strings.Builder

        // Header
        fmt.Fprintf(&sb, "// Generated by goscript v2.0\n")
        fmt.Fprintf(&sb, "// Bundle: %d files\n\n", len(bundle.Files))

        // Runtime preamble
        if b.IncludeRuntime {
                sb.WriteString("// --- Goscript Runtime ---\n")
                sb.WriteString("// Runtime: __gs client runtime\n")
                sb.WriteString("(function() {\n")
                sb.WriteString(RuntimeJS)
                sb.WriteString("\n})();\n\n")
                bundle.Stats.HasRuntime = true
        }

        // Stdlib usage comment
        if len(allImports) > 0 {
                sort.Strings(allImports)
                unique := dedup(allImports)
                fmt.Fprintf(&sb, "// Stdlib imports used: %s\n\n", strings.Join(unique, ", "))
        }

        // Compiled files
        for _, bf := range bundle.Files {
                fmt.Fprintf(&sb, "// --- File: %s ---\n", bf.ModuleName)
                if len(bf.Imports) > 0 {
                        fmt.Fprintf(&sb, "// Imports: %s\n", strings.Join(bf.Imports, ", "))
                }
                js := bf.JavaScript
                if b.Minify {
                        js = minifyJS(js)
                }
                sb.WriteString(js)
                sb.WriteString("\n\n")
        }

        // Footer
        sb.WriteString("// --- End of goscript bundle ---\n")

        // Apply minification to runtime if requested
        output := sb.String()
        if b.Minify && b.IncludeRuntime {
                // The runtime itself was already embedded; minify the full output
                output = minifyJS(output)
        }

        bundle.Output = output

        // Compute stats
        for _, bf := range bundle.Files {
                bundle.Stats.FileCount++
                bundle.Stats.GSSize += bf.SourceSize
                bundle.Stats.JSSize += bf.JSSize
        }

        // Write to disk if OutDir is set
        if b.OutDir != "" {
                if err := b.writeBundle(bundle); err != nil {
                        return nil, fmt.Errorf("writing bundle to %s: %w", b.OutDir, err)
                }
        }

        // Build source map from compiled files
        if b.SourceMap {
                for _, bf := range bundle.Files {
                        bundle.SourceMap[bf.ModuleName] = bf.Source
                }
        } else {
                bundle.SourceMap = nil
        }

        return bundle, nil
}

// compileFile reads a single .gs file, compiles it, and returns a BundleFile.
func (b *Bundler) compileFile(path string) (*BundleFile, error) {
        // Read source
        source, err := os.ReadFile(path)
        if err != nil {
                return nil, fmt.Errorf("reading %s: %w", path, err)
        }

        // Compile: Lexer → Parser → Generator
        input := string(source)

        tokens, err := NewLexer(input).Tokenize()
        if err != nil {
                return nil, fmt.Errorf("lexing: %w", err)
        }

        program, err := NewParser(tokens).Parse()
        if err != nil {
                return nil, fmt.Errorf("parsing: %w", err)
        }

        gen := NewGenerator()
        javascript, err := gen.Generate(program)
        if err != nil {
                return nil, fmt.Errorf("generating: %w", err)
        }

        // Extract imports from the AST
        imports := extractImports(program)

        bf := &BundleFile{
                Path:       path,
                ModuleName: strings.TrimSuffix(filepath.Base(path), ".gs"),
                JavaScript: javascript,
                Source:     input,
                Imports:    imports,
                SourceSize: len(source),
                JSSize:     len(javascript),
        }

        return bf, nil
}

// writeBundle writes the bundle output to <OutDir>/bundle.js.
func (b *Bundler) writeBundle(bundle *Bundle) error {
        if err := os.MkdirAll(b.OutDir, 0755); err != nil {
                return fmt.Errorf("creating output directory: %w", err)
        }
        outPath := filepath.Join(b.OutDir, "bundle.js")
        return os.WriteFile(outPath, []byte(bundle.Output), 0644)
}

// ---------------------------------------------------------------------------
// File discovery
// ---------------------------------------------------------------------------

// discoverGSFiles walks a directory tree and returns all .gs file paths,
// sorted alphabetically.
func discoverGSFiles(dir string) ([]string, error) {
        var files []string

        err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
                if err != nil {
                        return err
                }
                if d.IsDir() {
                        // Skip hidden directories (e.g. .git)
                        name := d.Name()
                        if strings.HasPrefix(name, ".") {
                                return filepath.SkipDir
                        }
                        return nil
                }
                if strings.HasSuffix(strings.ToLower(d.Name()), ".gs") {
                        files = append(files, path)
                }
                return nil
        })

        if err != nil {
                return nil, err
        }

        sort.Strings(files)
        return files, nil
}

// orderFiles places entry points first (in the order they were given),
// followed by remaining files in alphabetical order.
func orderFiles(allFiles, entryPoints []string) []string {
        if len(entryPoints) == 0 {
                return allFiles
        }

        var ordered []string
        var rest []string
        used := make(map[string]bool)

        // Add entry points first
        for _, ep := range entryPoints {
                abs, _ := filepath.Abs(ep)
                // Find matching file
                for _, f := range allFiles {
                        if f == ep || f == abs {
                                if !used[f] {
                                        used[f] = true
                                        ordered = append(ordered, f)
                                }
                                break
                        }
                }
        }

        // Add remaining files
        for _, f := range allFiles {
                if !used[f] {
                        rest = append(rest, f)
                }
        }

        return append(ordered, rest...)
}

// ---------------------------------------------------------------------------
// Import resolution
// ---------------------------------------------------------------------------

// extractImports returns the import paths from a parsed Program.
// Only goscript/* stdlib imports are collected; local relative imports
// are tracked but left for future resolution.
func extractImports(program *Program) []string {
        var imports []string
        for _, imp := range program.Imports {
                path := strings.Trim(imp.Path.Value, `"`)
                if isStdlibImport(path) {
                        imports = append(imports, path)
                }
        }
        return imports
}

// isStdlibImport returns true if the import path is a goscript stdlib
// import (e.g. "goscript/dom", "goscript/state") or a Go standard
// library import that the generator handles (e.g. "fmt", "strings").
func isStdlibImport(path string) bool {
        if strings.HasPrefix(path, "goscript/") {
                return true
        }
        // Check Go stdlib imports that have mappings
        _, ok := stdlibMapping[path]
        return ok
}

// resolveImports collects all stdlib import paths across all bundle files
// and returns a deduplicated, sorted list.
func (b *Bundler) resolveImports(files []BundleFile) []string {
        var all []string
        seen := make(map[string]bool)
        for _, f := range files {
                for _, imp := range f.Imports {
                        if !seen[imp] {
                                seen[imp] = true
                                all = append(all, imp)
                        }
                }
        }
        sort.Strings(all)
        return all
}

// ---------------------------------------------------------------------------
// Minification
// ---------------------------------------------------------------------------

// minifyJS performs simple JavaScript minification:
//   - Removes single-line comments (// ...)
//   - Removes multi-line comments (/* ... */)
//   - Collapses multiple blank lines into one
//   - Strips leading whitespace from each line
//   - Removes trailing whitespace
//
// This is a lightweight minifier for development use.  For production
// bundling, use a dedicated tool like terser.
func minifyJS(input string) string {
        var result strings.Builder
        lines := strings.Split(input, "\n")

        // Regex patterns
        singleLineComment := regexp.MustCompile(`^\s*//.*$`)
        blockCommentStart := regexp.MustCompile(`^\s*/\*`)
        blockCommentEnd := regexp.MustCompile(`\*/`)
        blankLine := regexp.MustCompile(`^\s*$`)

        inBlockComment := false
        prevBlank := false

        for _, line := range lines {
                // Handle block comments
                if inBlockComment {
                        if blockCommentEnd.MatchString(line) {
                                inBlockComment = false
                        }
                        continue
                }

                if blockCommentStart.MatchString(line) && !blockCommentEnd.MatchString(line) {
                        inBlockComment = true
                        continue
                }

                // Skip full-line single-line comments (but not URLs like http://)
                trimmed := strings.TrimSpace(line)
                if singleLineComment.MatchString(line) && !strings.Contains(trimmed, "://") {
                        // Keep file separator comments (// --- File: name ---)
                        if strings.Contains(trimmed, "---") {
                                result.WriteString(trimmed)
                                result.WriteByte('\n')
                                prevBlank = false
                                continue
                        }
                        continue
                }

                // Collapse blank lines
                if blankLine.MatchString(line) {
                        if !prevBlank {
                                result.WriteByte('\n')
                                prevBlank = true
                        }
                        continue
                }

                // Strip leading whitespace
                stripped := strings.TrimLeft(line, " \t")
                // Strip trailing whitespace
                stripped = strings.TrimRight(stripped, " \t\r")

                result.WriteString(stripped)
                result.WriteByte('\n')
                prevBlank = false
        }

        return strings.TrimRight(result.String(), " \t\n")
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// dedup returns a new slice with duplicate strings removed, preserving order.
func dedup(items []string) []string {
        seen := make(map[string]bool)
        var result []string
        for _, item := range items {
                if !seen[item] {
                        seen[item] = true
                        result = append(result, item)
                }
        }
        return result
}
