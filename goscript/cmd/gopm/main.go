// Package main implements the gopm CLI — the GoScript Package Manager and build tool.
//
// gopm is the central command-line interface for GoScript 2.0 projects. It handles:
//   - Compiling .gs (GoScript) source files to browser-compatible JavaScript
//   - Starting a development server with Hot Module Replacement (HMR)
//   - Building production-ready output with embedded .gs compilations
//   - Scaffolding new goscript projects from templates
//
// Usage:
//
//      gopm dev           Start dev server with HMR
//      gopm build         Compile .gs files for production
//      gopm compile       Compile a single .gs file to .js
//      gopm init [name]   Create a new goscript project
//      gopm version       Show version
//      gopm help          Show help
package main

import (
        "bufio"
        "embed"
        "flag"
        "fmt"
        "io/fs"
        "log"
        "net/http"
        "os"
        "os/exec"
        "path/filepath"
        "runtime"
        "strings"
        "sync"
        "time"

        gscompiler "github.com/davidjeba/goscript/pkg/gscompiler"
)

// Version is the gopm CLI version, set at build time via ldflags.
var Version = "2.0.0"

// ---------------------------------------------------------------------------
// Exit codes
// ---------------------------------------------------------------------------

const (
        exitOK       = 0
        exitUsage    = 1
        exitCompile  = 2
        exitFileIO   = 3
        exitBuild    = 4
        exitDev      = 5
        exitInit     = 6
)

// ---------------------------------------------------------------------------
// Entry point
// ---------------------------------------------------------------------------

func main() {
        if len(os.Args) < 2 {
                printHelp()
                os.Exit(exitUsage)
        }

        command := os.Args[1]
        subArgs := os.Args[2:]

        os.Exit(runCommand(command, subArgs))
}

// runCommand dispatches to the appropriate subcommand handler and returns an
// exit code.
func runCommand(command string, args []string) int {
        switch command {
        case "dev":
                return cmdDev(args)
        case "build":
                return cmdBuild(args)
        case "compile":
                return cmdCompile(args)
        case "init":
                return cmdInit(args)
        case "version", "--version", "-v":
                cmdVersion()
                return exitOK
        case "help", "--help", "-h":
                printHelp()
                return exitOK
        default:
                fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", command)
                printHelp()
                return exitUsage
        }
}

// ---------------------------------------------------------------------------
// gopm compile — compile a single .gs file to .js
// ---------------------------------------------------------------------------

// compileFlags holds the flag set for the compile subcommand.
type compileFlags struct {
        output  string
        watch   bool
        minify  bool
        sourcemap bool
}

func cmdCompile(args []string) int {
        fs := flag.NewFlagSet("compile", flag.ExitOnError)
        f := &compileFlags{}
        fs.StringVar(&f.output, "o", "", "Output file path (default: stdout)")
        fs.StringVar(&f.output, "output", "", "Output file path (default: stdout)")
        fs.BoolVar(&f.watch, "w", false, "Watch input file for changes and recompile")
        fs.BoolVar(&f.watch, "watch", false, "Watch input file for changes and recompile")
        fs.BoolVar(&f.minify, "minify", false, "Minify the output JavaScript")
        fs.BoolVar(&f.sourcemap, "sourcemap", false, "Generate a source map")
        fs.Parse(args)

        remaining := fs.Args()
        if len(remaining) < 1 {
                fmt.Fprintf(os.Stderr, "Usage: gopm compile <input.gs> [-o output.js]\n")
                return exitUsage
        }

        inputPath := remaining[0]

        if f.watch {
                return watchAndCompile(inputPath, f)
        }

        return compileFile(inputPath, f)
}

// compileFile reads a .gs source file, runs it through the gscompiler
// (lexer → parser → generator), and writes the resulting JavaScript to the
// specified output or stdout.
func compileFile(inputPath string, f *compileFlags) int {
        // Read source
        src, err := os.ReadFile(inputPath)
        if err != nil {
                fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", inputPath, err)
                return exitFileIO
        }

        // Compile: Lexer → Parser → Generator
        lexer := gscompiler.NewLexer(string(src))
        tokens := lexer.Tokenize()
        if tokens == nil {
                fmt.Fprintf(os.Stderr, "Error: lexer returned no tokens for %s\n", inputPath)
                return exitCompile
        }

        parser := gscompiler.NewParser(tokens)
        program := parser.Parse()
        if program == nil {
                fmt.Fprintf(os.Stderr, "Error: parser returned nil program for %s\n", inputPath)
                return exitCompile
        }

        generator := gscompiler.NewGenerator()
        jsOutput := generator.Generate(program)

        if f.minify {
                jsOutput = minifyJS(jsOutput)
        }

        // Write output
        if f.output == "" {
                fmt.Println(jsOutput)
        } else {
                err = os.WriteFile(f.output, []byte(jsOutput), 0644)
                if err != nil {
                        fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", f.output, err)
                        return exitFileIO
                }
                absInput, _ := filepath.Abs(inputPath)
                absOutput, _ := filepath.Abs(f.output)
                fmt.Printf("Compiled %s → %s (%d bytes)\n", absInput, absOutput, len(jsOutput))
        }

        return exitOK
}

// watchAndCompile monitors a .gs file for changes and recompiles on every
// modification, printing timing information to stderr.
func watchAndCompile(inputPath string, f *compileFlags) int {
        absPath, err := filepath.Abs(inputPath)
        if err != nil {
                fmt.Fprintf(os.Stderr, "Error resolving path: %v\n", err)
                return exitFileIO
        }

        fmt.Fprintf(os.Stderr, "Watching %s for changes...\n", absPath)

        var lastMod time.Time
        ticker := time.NewTicker(500 * time.Millisecond)
        defer ticker.Stop()

        for range ticker.C {
                info, err := os.Stat(absPath)
                if err != nil {
                        continue
                }
                if info.ModTime().After(lastMod) {
                        lastMod = info.ModTime()
                        start := time.Now()
                        code := compileFile(inputPath, f)
                        elapsed := time.Since(start)
                        if code == exitOK {
                                fmt.Fprintf(os.Stderr, "Recompiled in %s\n", elapsed)
                        }
                }
        }

        return exitOK
}

// ---------------------------------------------------------------------------
// gopm dev — development server with HMR
// ---------------------------------------------------------------------------

// devFlags holds the flag set for the dev subcommand.
type devFlags struct {
        port   int
        hmr    bool
        open   bool
        static string
}

func cmdDev(args []string) int {
        fs := flag.NewFlagSet("dev", flag.ExitOnError)
        f := &devFlags{
                port: 8080,
                hmr:  true,
        }
        fs.IntVar(&f.port, "port", 8080, "Dev server port")
        fs.IntVar(&f.port, "p", 8080, "Dev server port")
        fs.BoolVar(&f.hmr, "hmr", true, "Enable Hot Module Replacement")
        fs.BoolVar(&f.open, "open", false, "Open browser automatically")
        fs.StringVar(&f.static, "static", "./static", "Directory for static assets")
        fs.Parse(args)

        // Discover .gs files to watch
        gsFiles, err := findGSFiles(".")
        if err != nil {
                fmt.Fprintf(os.Stderr, "Error scanning for .gs files: %v\n", err)
                return exitDev
        }

        if len(gsFiles) == 0 {
                fmt.Fprintln(os.Stderr, "No .gs files found in the current directory.")
                fmt.Fprintln(os.Stderr, "Create a .gs file or run: gopm init myapp")
        }

        fmt.Fprintf(os.Stderr, "Found %d .gs file(s) to watch\n", len(gsFiles))

        // Compile all .gs files into memory for the dev server
        compiled := &sync.Map{}
        for _, gf := range gsFiles {
                js, err := compileGSToJS(gf)
                if err != nil {
                        fmt.Fprintf(os.Stderr, "Warning: could not compile %s: %v\n", gf, err)
                        continue
                }
                compiled.Store(gf, js)
        }

        mux := http.NewServeMux()

        // Serve compiled .gs files as .js
        mux.HandleFunc("/__gs/", func(w http.ResponseWriter, r *http.Request) {
                gsPath := strings.TrimPrefix(r.URL.Path, "/__gs/")
                gsPath = filepath.Clean(gsPath)
                jsPath := strings.TrimSuffix(gsPath, ".js") + ".gs"

                if js, ok := compiled.Load(jsPath); ok {
                        w.Header().Set("Content-Type", "application/javascript")
                        w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
                        fmt.Fprint(w, js)
                        return
                }

                http.NotFound(w, r)
        })

        // HMR WebSocket endpoint
        mux.HandleFunc("/__hmr", func(w http.ResponseWriter, r *http.Request) {
                // Upgrade to WebSocket (using gorilla/websocket is optional; for a
                // basic implementation we can use the standard library's hijacker)
                upgrader := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                        w.Header().Set("Content-Type", "text/event-stream")
                        w.Header().Set("Cache-Control", "no-cache")
                        w.Header().Set("Connection", "keep-alive")
                        flusher, ok := w.(http.Flusher)
                        if !ok {
                                http.Error(w, "streaming not supported", http.StatusInternalServerError)
                                return
                        }
                        fmt.Fprint(w, "data: {\"type\":\"connected\",\"message\":\"gopm HMR\"}\n\n")
                        flusher.Flush()

                        // Keep connection alive for HMR events
                        ctx := r.Context()
                        ticker := time.NewTicker(15 * time.Second)
                        defer ticker.Stop()
                        for {
                                select {
                                case <-ctx.Done():
                                        return
                                case <-ticker.C:
                                        fmt.Fprint(w, ": ping\n\n")
                                        flusher.Flush()
                                }
                        }
                })
                upgrader.ServeHTTP(w, r)
        })

        // Serve static files
        if _, err := os.Stat(f.static); err == nil {
                fileServer := http.FileServer(http.Dir(f.static))
                mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
        }

        // Fallback: serve index.html for SPA-like behavior
        mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
                indexPath := filepath.Join(f.static, "index.html")
                if _, err := os.Stat(indexPath); err == nil {
                        w.Header().Set("Content-Type", "text/html")
                        http.ServeFile(w, r, indexPath)
                        return
                }
                http.NotFound(w, r)
        })

        addr := fmt.Sprintf(":%d", f.port)

        fmt.Fprintf(os.Stderr, `
  ╔═══════════════════════════════════════╗
  ║     gopm dev server                  ║
  ║     GoScript %s                     ║
  ╠═══════════════════════════════════════╣
  ║  Local:  http://localhost:%-12s ║
  ║  HMR:    %s%-28s ║
  ║  Files:  %-30s ║
  ╚═══════════════════════════════════════╝
`, Version, f.port, boolStr(f.hmr), "enabled", fmt.Sprintf("%d .gs file(s)", len(gsFiles)))

        // Background watcher
        go func() {
                watchGSFilesForHMR(".", compiled, gsFiles)
        }()

        if err := http.ListenAndServe(addr, mux); err != nil {
                fmt.Fprintf(os.Stderr, "Dev server error: %v\n", err)
                return exitDev
        }

        return exitOK
}

// ---------------------------------------------------------------------------
// gopm build — compile all .gs files for production
// ---------------------------------------------------------------------------

// buildFlags holds the flag set for the build subcommand.
type buildFlags struct {
        output string
        minify bool
}

func cmdBuild(args []string) int {
        fs := flag.NewFlagSet("build", flag.ExitOnError)
        f := &buildFlags{
                output: "./build",
        }
        fs.StringVar(&f.output, "o", "./build", "Build output directory")
        fs.StringVar(&f.output, "output", "./build", "Build output directory")
        fs.BoolVar(&f.minify, "minify", false, "Minify output JavaScript")
        fs.Parse(args)

        // Find all .gs files
        gsFiles, err := findGSFiles(".")
        if err != nil {
                fmt.Fprintf(os.Stderr, "Error scanning for .gs files: %v\n", err)
                return exitBuild
        }

        if len(gsFiles) == 0 {
                fmt.Fprintln(os.Stderr, "No .gs files found. Nothing to build.")
                return exitBuild
        }

        fmt.Fprintf(os.Stderr, "Building %d .gs file(s)...\n", len(gsFiles))

        // Create output directory
        if err := os.MkdirAll(f.output, 0755); err != nil {
                fmt.Fprintf(os.Stderr, "Error creating build directory: %v\n", err)
                return exitBuild
        }

        successCount := 0
        errorCount := 0
        totalBytes := 0

        for _, gf := range gsFiles {
                js, err := compileGSToJS(gf)
                if err != nil {
                        fmt.Fprintf(os.Stderr, "  ✗ %s: %v\n", gf, err)
                        errorCount++
                        continue
                }

                if f.minify {
                        js = minifyJS(js)
                }

                // Build the output path: mirror the source tree under the build dir
                outPath := filepath.Join(f.output, strings.TrimSuffix(gf, ".gs")+".js")
                outDir := filepath.Dir(outPath)
                if err := os.MkdirAll(outDir, 0755); err != nil {
                        fmt.Fprintf(os.Stderr, "  ✗ %s: cannot create directory: %v\n", gf, err)
                        errorCount++
                        continue
                }

                if err := os.WriteFile(outPath, []byte(js), 0644); err != nil {
                        fmt.Fprintf(os.Stderr, "  ✗ %s: %v\n", gf, err)
                        errorCount++
                        continue
                }

                fmt.Fprintf(os.Stderr, "  ✓ %s → %s (%d bytes)\n", gf, outPath, len(js))
                totalBytes += len(js)
                successCount++
        }

        // Copy static files if ./static exists
        staticDir := "./static"
        if info, err := os.Stat(staticDir); err == nil && info.IsDir() {
                buildStatic := filepath.Join(f.output, "static")
                if err := copyDir(staticDir, buildStatic); err != nil {
                        fmt.Fprintf(os.Stderr, "Warning: could not copy static files: %v\n", err)
                } else {
                        fmt.Fprintf(os.Stderr, "  Copied static/ → %s\n", buildStatic)
                }
        }

        fmt.Fprintf(os.Stderr, "\nBuild complete: %d succeeded, %d failed, %d total bytes\n",
                successCount, errorCount, totalBytes)

        if errorCount > 0 {
                return exitBuild
        }

        // Generate a manifest
        manifest := generateManifest(gsFiles, f.output)
        manifestPath := filepath.Join(f.output, "manifest.json")
        if err := os.WriteFile(manifestPath, []byte(manifest), 0644); err != nil {
                fmt.Fprintf(os.Stderr, "Warning: could not write manifest: %v\n", err)
        }

        return exitOK
}

// ---------------------------------------------------------------------------
// gopm init — scaffold a new goscript project
// ---------------------------------------------------------------------------

func cmdInit(args []string) int {
        var name string
        if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
                name = args[0]
        }

        if name == "" {
                fmt.Print("Project name: ")
                reader := bufio.NewReader(os.Stdin)
                input, _ := reader.ReadString('\n')
                name = strings.TrimSpace(input)
        }

        if name == "" {
                fmt.Fprintln(os.Stderr, "Error: project name is required")
                return exitInit
        }

        // Sanitize name
        name = sanitizeProjectName(name)
        projectDir := "./" + name

        if _, err := os.Stat(projectDir); err == nil {
                fmt.Fprintf(os.Stderr, "Error: directory %s already exists\n", projectDir)
                return exitInit
        }

        fmt.Fprintf(os.Stderr, "Creating goscript project: %s\n", name)

        // Create directory structure
        dirs := []string{
                filepath.Join(projectDir, "app"),
                filepath.Join(projectDir, "app/pages"),
                filepath.Join(projectDir, "app/components"),
                filepath.Join(projectDir, "app/layouts"),
                filepath.Join(projectDir, "app/api"),
                filepath.Join(projectDir, "public"),
                filepath.Join(projectDir, "static"),
        }

        for _, dir := range dirs {
                if err := os.MkdirAll(dir, 0755); err != nil {
                        fmt.Fprintf(os.Stderr, "Error creating directory %s: %v\n", dir, err)
                        return exitInit
                }
        }

        // Write go.mod
        goMod := fmt.Sprintf(moduleTemplate, name)
        if err := os.WriteFile(filepath.Join(projectDir, "go.mod"), []byte(goMod), 0644); err != nil {
                fmt.Fprintf(os.Stderr, "Error writing go.mod: %v\n", err)
                return exitInit
        }

        // Write main.go
        mainGo := fmt.Sprintf(mainTemplate, name)
        if err := os.WriteFile(filepath.Join(projectDir, "main.go"), []byte(mainGo), 0644); err != nil {
                fmt.Fprintf(os.Stderr, "Error writing main.go: %v\n", err)
                return exitInit
        }

        // Write home page
        homePage := fmt.Sprintf(homePageTemplate, name)
        if err := os.WriteFile(filepath.Join(projectDir, "app/pages/home.gs"), []byte(homePage), 0644); err != nil {
                fmt.Fprintf(os.Stderr, "Error writing home page: %v\n", err)
                return exitInit
        }

        // Write root layout
        rootLayout := layoutTemplate
        if err := os.WriteFile(filepath.Join(projectDir, "app/layouts/root.gs"), []byte(rootLayout), 0644); err != nil {
                fmt.Fprintf(os.Stderr, "Error writing root layout: %v\n", err)
                return exitInit
        }

        // Write .gitignore
        gitignore := gitignoreTemplate
        if err := os.WriteFile(filepath.Join(projectDir, ".gitignore"), []byte(gitignore), 0644); err != nil {
                fmt.Fprintf(os.Stderr, "Error writing .gitignore: %v\n", err)
                return exitInit
        }

        fmt.Fprintf(os.Stderr, `
  ✓ Project created successfully!

  Next steps:
    cd %s
    gopm dev          # Start the development server
    gopm build        # Build for production
    go run main.go    # Run the Go server directly
`, name)

        return exitOK
}

// ---------------------------------------------------------------------------
// gopm version
// ---------------------------------------------------------------------------

func cmdVersion() {
        fmt.Printf("gopm %s (%s/%s)\n", Version, runtime.GOOS, runtime.GOARCH)
        fmt.Printf("goscript %s\n", "2.0.0")
}

// ---------------------------------------------------------------------------
// Help
// ---------------------------------------------------------------------------

func printHelp() {
        help := `
gopm — GoScript Package Manager & Build Tool

Usage:
  gopm <command> [options]

Commands:
  compile <file.gs>    Compile a single .gs file to JavaScript
                       Output to stdout or -o file.js

  dev                  Start development server with HMR
                       Compiles .gs → .js on change, serves assets

  build                Compile all .gs files for production
                       Outputs to ./build by default

  init [name]          Create a new goscript project from template
                       Sets up directory structure and starter files

  version              Show gopm and goscript version

  help                 Show this help message

Compile Flags:
  -o, --output <path>  Write output to file instead of stdout
  -w, --watch          Watch input file and recompile on change
  --minify             Minify the generated JavaScript
  --sourcemap          Generate a source map file

Dev Server Flags:
  -p, --port <port>    Server port (default: 8080)
  --hmr                Enable Hot Module Replacement (default: true)
  --static <dir>       Static files directory (default: ./static)
  --open               Open browser on start

Build Flags:
  -o, --output <dir>   Build output directory (default: ./build)
  --minify             Minify all generated JavaScript

Examples:
  gopm compile app/pages/home.gs -o build/home.js
  gopm compile counter.gs --watch
  gopm dev --port 3000
  gopm build --minify -o ./dist
  gopm init myapp
  gopm version

GS Language (.gs → .js):
  .gs files are GoScript's own language that compiles to browser JavaScript.
  Use the gscompiler pipeline:
    Lexer.Tokenize() → Parser.Parse() → Generator.Generate()

  gopm handles this automatically — just write .gs and run gopm dev.
`
        fmt.Println(strings.TrimSpace(help))
}

// ---------------------------------------------------------------------------
// Helper: compile a .gs file to JS string
// ---------------------------------------------------------------------------

func compileGSToJS(path string) (string, error) {
        src, err := os.ReadFile(path)
        if err != nil {
                return "", fmt.Errorf("read %s: %w", path, err)
        }

        lexer := gscompiler.NewLexer(string(src))
        tokens := lexer.Tokenize()
        if tokens == nil {
                return "", fmt.Errorf("lexer produced no tokens for %s", path)
        }

        parser := gscompiler.NewParser(tokens)
        program := parser.Parse()
        if program == nil {
                return "", fmt.Errorf("parser produced no AST for %s", path)
        }

        generator := gscompiler.NewGenerator()
        js := generator.Generate(program)

        return js, nil
}

// ---------------------------------------------------------------------------
// Helper: find all .gs files in a directory tree
// ---------------------------------------------------------------------------

func findGSFiles(root string) ([]string, error) {
        var files []string
        err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
                if err != nil {
                        return err
                }
                if !d.IsDir() && strings.HasSuffix(path, ".gs") {
                        rel, err := filepath.Rel(root, path)
                        if err != nil {
                                return err
                        }
                        // Skip hidden directories and build output
                        if strings.Contains(rel, "/.") || strings.Contains(rel, "build/") {
                                return nil
                        }
                        files = append(files, rel)
                }
                return nil
        })
        return files, err
}

// ---------------------------------------------------------------------------
// Helper: HMR watcher for dev server
// ---------------------------------------------------------------------------

func watchGSFilesForHMR(root string, compiled *sync.Map, knownFiles []string) {
        snapshots := make(map[string]time.Time)
        for _, f := range knownFiles {
                if info, err := os.Stat(filepath.Join(root, f)); err == nil {
                        snapshots[f] = info.ModTime()
                }
        }

        ticker := time.NewTicker(500 * time.Millisecond)
        defer ticker.Stop()

        for range ticker.C {
                files, err := findGSFiles(root)
                if err != nil {
                        continue
                }

                for _, gf := range files {
                        fullPath := filepath.Join(root, gf)
                        info, err := os.Stat(fullPath)
                        if err != nil {
                                continue
                        }

                        prevMod, exists := snapshots[gf]
                        if !exists || info.ModTime().After(prevMod) {
                                snapshots[gf] = info.ModTime()

                                start := time.Now()
                                js, err := compileGSToJS(gf)
                                elapsed := time.Since(start)

                                if err != nil {
                                        log.Printf("[HMR] Compile error %s: %v", gf, err)
                                        continue
                                }

                                compiled.Store(gf, js)
                                log.Printf("[HMR] Recompiled %s in %s (%d bytes)", gf, elapsed, len(js))
                        }
                }
        }
}

// ---------------------------------------------------------------------------
// Helper: simple JS minifier (removes comments and extra whitespace)
// ---------------------------------------------------------------------------

func minifyJS(input string) string {
        // Remove single-line comments
        result := removeSingleLineComments(input)
        // Remove multi-line comments
        result = removeMultiLineComments(result)
        // Collapse whitespace
        result = collapseWhitespace(result)
        return result
}

func removeSingleLineComments(input string) string {
        var sb strings.Builder
        lines := strings.Split(input, "\n")
        for _, line := range lines {
                if idx := strings.Index(line, "//"); idx >= 0 {
                        // Check if the // is inside a string
                        inString := false
                        quoteChar := byte(0)
                        for i := 0; i < idx; i++ {
                                c := line[i]
                                if inString {
                                        if c == quoteChar {
                                                inString = false
                                        }
                                } else {
                                        if c == '"' || c == '\'' || c == '`' {
                                                inString = true
                                                quoteChar = c
                                        }
                                }
                        }
                        if !inString {
                                line = line[:idx]
                        }
                }
                sb.WriteString(line)
                sb.WriteByte('\n')
        }
        return sb.String()
}

func removeMultiLineComments(input string) string {
        var sb strings.Builder
        i := 0
        for i < len(input) {
                if i+1 < len(input) && input[i] == '/' && input[i+1] == '*' {
                        end := strings.Index(input[i+2:], "*/")
                        if end >= 0 {
                                i += end + 4
                                continue
                        }
                }
                sb.WriteByte(input[i])
                i++
        }
        return sb.String()
}

func collapseWhitespace(input string) string {
        var sb strings.Builder
        prevSpace := false
        for _, c := range input {
                if c == '\n' || c == '\r' || c == '\t' {
                        if !prevSpace {
                                sb.WriteByte(' ')
                                prevSpace = true
                        }
                        continue
                }
                if c == ' ' {
                        if prevSpace {
                                continue
                        }
                        prevSpace = true
                } else {
                        prevSpace = false
                }
                sb.WriteRune(c)
        }
        return sb.String()
}

// ---------------------------------------------------------------------------
// Helper: copy a directory recursively
// ---------------------------------------------------------------------------

func copyDir(src, dst string) error {
        return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
                if err != nil {
                        return err
                }
                rel, err := filepath.Rel(src, path)
                if err != nil {
                        return err
                }
                target := filepath.Join(dst, rel)

                if d.IsDir() {
                        return os.MkdirAll(target, 0755)
                }

                data, err := os.ReadFile(path)
                if err != nil {
                        return err
                }
                return os.WriteFile(target, data, 0644)
        })
}

// ---------------------------------------------------------------------------
// Helper: generate build manifest
// ---------------------------------------------------------------------------

func generateManifest(files []string, outputDir string) string {
        var sb strings.Builder
        sb.WriteString("{\n")
        sb.WriteString(fmt.Sprintf("  \"version\": \"2.0.0\",\n"))
        sb.WriteString(fmt.Sprintf("  \"timestamp\": \"%s\",\n", time.Now().UTC().Format(time.RFC3339)))
        sb.WriteString("  \"files\": [\n")
        for i, f := range files {
                outPath := filepath.Join(outputDir, strings.TrimSuffix(f, ".gs")+".js")
                sb.WriteString(fmt.Sprintf("    {\"source\": \"%s\", \"output\": \"%s\"}", f, outPath))
                if i < len(files)-1 {
                        sb.WriteString(",")
                }
                sb.WriteString("\n")
        }
        sb.WriteString("  ]\n")
        sb.WriteString("}\n")
        return sb.String()
}

// ---------------------------------------------------------------------------
// Helper: sanitize project name
// ---------------------------------------------------------------------------

func sanitizeProjectName(name string) string {
        name = strings.TrimSpace(name)
        name = strings.ToLower(name)
        name = strings.ReplaceAll(name, " ", "-")
        // Remove non-alphanumeric characters except hyphens
        var sb strings.Builder
        for _, c := range name {
                if (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '-' {
                        sb.WriteRune(c)
                }
        }
        result := sb.String()
        if result == "" {
                result = "goscript-app"
        }
        return result
}

// ---------------------------------------------------------------------------
// Helper: bool to string
// ---------------------------------------------------------------------------

func boolStr(b bool) string {
        if b {
                return "true"
        }
        return "false"
}

// ---------------------------------------------------------------------------
// Project template strings
// ---------------------------------------------------------------------------

const moduleTemplate = `module %s

go 1.22

require github.com/davidjeba/goscript v2.0.0
`

const mainTemplate = `package main

import (
        "fmt"
        "log"
        "net/http"
        "os"

        goscript "github.com/davidjeba/goscript/pkg/goscript"
)

func main() {
        port := "8080"
        if p := os.Getenv("PORT"); p != "" {
                port = p
        }

        router := goscript.NewAppRouter("/")

        // Home page
        router.RegisterRoute("/", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
                metadata := goscript.NewMetadata().
                        SetTitle("%s — GoScript App").
                        SetDescription("Built with GoScript 2.0").
                        Build()

                html := fmt.Sprintf("%%s<!DOCTYPE html>
<html lang=\"en\">
<head>
        <meta charset=\"utf-8\">
        <meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">
        %%s
</head>
<body>
        <div style=\"max-width:800px;margin:2rem auto;font-family:system-ui,sans-serif\">
                <h1>Welcome to %s</h1>
                <p>Edit <code>app/pages/home.gs</code> and save to see changes.</p>
        </div>
</body>
</html>", metadata.Render())

                w.Header().Set("Content-Type", "text/html")
                fmt.Fprint(w, html)
        }, []string{"GET"})

        fmt.Printf("Server running at http://localhost:%%s\n", port)
        log.Fatal(http.ListenAndServe(":"+port, router))
}
`

const homePageTemplate = `package main

import "goscript/dom"

// Home page — this .gs file compiles to JavaScript via gopm
func Home() dom.Element {
        return dom.CreateElement("div", dom.Props{
                "class": "page",
        },
                dom.CreateElement("h1", nil, "Welcome to %s"),
                dom.CreateElement("p", nil, "Edit this file and save to see hot reload."),
        )
}
`

const layoutTemplate = `package main

import "goscript/dom"

// Root layout — wraps all pages
func RootLayout(children ...dom.Element) dom.Element {
        return dom.CreateElement("div", dom.Props{
                "class": "layout",
        },
                dom.CreateElement("header", nil,
                        dom.CreateElement("h1", nil, "%s"),
                ),
                dom.CreateElement("main", nil, children...),
                dom.CreateElement("footer", nil,
                        dom.CreateElement("p", nil, "Powered by GoScript"),
                ),
        )
}
`

const gitignoreTemplate = `# Build output
/build/
/dist/

# Dependencies
/vendor/

# IDE
.idea/
.vscode/
*.swp
*.swo

# OS
.DS_Store
Thumbs.db

# Environment
.env
.env.local
`
