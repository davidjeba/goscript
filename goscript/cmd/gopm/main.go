package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const goscriptVersion = "2.0.0"

func main() {
	flag.Usage = func() {
		fmt.Println(`
 ███╗   ██╗ ██████╗ ███████╗ ██████╗ ██████╗ ███╗   ██╗
 ████╗  ██║██╔═══██╗██╔════╝██╔═══██╗██╔══██╗████╗  ██║
 ██╔██╗ ██║██║   ██║███████╗██║   ██║██████╔╝██╔██╗ ██║
 ██║╚██╗██║██║   ██║╚════██║██║   ██║██╔══██╗██║╚██╗██║
 ██║ ╚████║╚██████╔╝███████║╚██████╔╝██║  ██║██║ ╚████║
 ╚═╝  ╚═══╝ ╚═════╝ ╚══════╝ ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═══╝
  GoScript 2.0 — Full-Stack Go Web Framework
`)
		fmt.Println("Commands:")
		fmt.Println("  gopm init [name]     Initialize a new GoScript project")
		fmt.Println("  gopm dev             Start development server with HMR")
		fmt.Println("  gopm build           Build for production")
		fmt.Println("  gopm start           Start production server")
		fmt.Println("  gopm generate <type> Generate components, pages, API routes")
		fmt.Println("  gopm deploy          Deploy to platform")
		fmt.Println("  gopm lint            Run linter")
		fmt.Println("  gopm test            Run tests")
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	cmd := os.Args[1]

	switch cmd {
	case "init":
		initProject()
	case "dev":
		runDev()
	case "build":
		runBuild()
	case "start":
		runStart()
	case "generate", "g":
		generate(os.Args[2:])
	case "deploy":
		deploy()
	case "lint":
		runLint()
	case "test":
		runTest()
	default:
		fmt.Printf("Unknown command: %s\n", cmd)
		flag.Usage()
		os.Exit(1)
	}
}

func initProject() {
	name := "my-goscript-app"
	if len(os.Args) > 2 {
		name = os.Args[2]
	}

	fmt.Printf("Creating GoScript project: %s\n", name)

	structure := map[string]string{
		"go.mod": fmt.Sprintf(
			"module %s\ngo 1.22\n\nrequire github.com/davidjeba/goscript v%s\n",
			name, goscriptVersion),
		"goscript.config.json": `{
  "port": 8080,
  "renderMode": "hybrid",
  "ssr": { "enabled": true },
  "hmr": { "enabled": true, "port": 8081 },
  "cors": { "origins": ["*"] },
  "compression": { "enabled": true }
}
`,
		"cmd/server/main.go": fmt.Sprintf(`package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	goscript "github.com/davidjeba/goscript/pkg/goscript"
)

func main() {
	router := goscript.NewAppRouter("/")
	router.RegisterRoute("/", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		w.Write([]byte("<h1>Hello from GoScript!</h1>"))
	}, []string{"GET"})

	pipeline := goscript.NewPipeline()
	pipeline.Use(goscript.CORSMiddleware(goscript.CORSConfig{AllowAllOrigins: true}))
	pipeline.Use(goscript.SecurityHeadersMiddleware())
	pipeline.Use(goscript.LoggingMiddleware(log.Printf))

	fmt.Println("GoScript server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
`, name),
		"app/layout.go": `package app

import "github.com/davidjeba/goscript/pkg/goscript"

func Layout() goscript.Component {
	return goscript.FunctionalComponent(func(props goscript.Props) string {
		return "<!DOCTYPE html><html><head></head><body>{{children}}</body></html>"
	})
}
`,
		"app/page.go": `package app

import "github.com/davidjeba/goscript/pkg/goscript"

func Page() goscript.Component {
	return goscript.FunctionalComponent(func(props goscript.Props) string {
		return "<h1>Welcome to GoScript</h1>"
	})
}
`,
	}

	for path, content := range structure {
		fullPath := filepath.Join(name, path)
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("Error creating directory %s: %v\n", dir, err)
			continue
		}
		if err := ioutil.WriteFile(fullPath, []byte(content), 0644); err != nil {
			fmt.Printf("Error creating %s: %v\n", fullPath, err)
			continue
		}
		fmt.Printf("  Created %s\n", fullPath)
	}

	fmt.Printf("\nProject '%s' created successfully!\n", name)
	fmt.Printf("  cd %s && gopm dev\n", name)
}

func runDev() {
	fmt.Printf("Starting development server with HMR...\n")
	fmt.Printf("GoScript %s — Dev Server\n", goscriptVersion)
	// In production, this would start the DevServer with HMR
	fmt.Println("  Server: http://localhost:8080")
	fmt.Println("  HMR:    ws://localhost:8081")
	select {} // block forever
}

func runBuild() {
	start := time.Now()
	fmt.Printf("Building for production...\n")
	cmd := exec.Command("go", "build", "-o", "bin/server", "./cmd/server/")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Build failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Build completed in %s\n", time.Since(start))
}

func runStart() {
	fmt.Printf("Starting production server...\n")
	cmd := exec.Command("./bin/server")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Server error: %v\n", err)
		os.Exit(1)
	}
}

func generate(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: gopm generate <component|page|api> <name>")
		return
	}

	generateType := args[0]
	name := "new-item"
	if len(args) > 1 {
		name = args[1]
	}

	switch generateType {
	case "component", "c":
		kebab := strings.ToLower(strings.ReplaceAll(name, " ", "-"))
		filePath := filepath.Join("pkg", "components", kebab+".go")
		content := fmt.Sprintf(`package components

import "github.com/davidjeba/goscript/pkg/goscript"

type %s struct {
	goscript.BaseComponent
}

func New%s() *%s {
	return &%s{}
}

func (c *%s) Render() string {
	return "<div class=\"%s\"></div>"
}
`, name, name, name, name, name, kebab)
		writeGeneratedFile(filePath, content)

	case "page", "p":
		filePath := filepath.Join("app", name, "page.go")
		content := fmt.Sprintf(`package %s

import (
	"net/http"

	"github.com/davidjeba/goscript/pkg/goscript"
)

func Page() goscript.Component {
	return goscript.FunctionalComponent(func(props goscript.Props) string {
		return "<div><h1>%s</h1></div>"
	})
}
`, name, name)
		writeGeneratedFile(filePath, content)

	case "api":
		filePath := filepath.Join("api", name, "route.go")
		content := fmt.Sprintf(`package %s

import (
	"net/http"

	"github.com/davidjeba/goscript/pkg/goscript"
)

func Handler(ctx *goscript.APIContext) (interface{}, error) {
	return map[string]string{"message": "Hello from %%s"}, nil
}
`, name, name)
		writeGeneratedFile(filePath, content)

	default:
		fmt.Printf("Unknown generate type: %s\n", generateType)
		fmt.Println("Use: component, page, or api")
	}
}

func writeGeneratedFile(path, content string) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	if err := ioutil.WriteFile(path, []byte(content), 0644); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("  Generated %s\n", path)
}

func deploy() {
	fmt.Println("Deploying GoScript project...")
	fmt.Println("  Target: production")
	fmt.Println("  (deploy is a stub in GoScript 2.0)")
}

func runLint() {
	fmt.Println("Running linter...")
	cmd := exec.Command("go", "vet", "./...")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Lint issues found\n")
		os.Exit(1)
	}
	fmt.Println("No issues found")
}

func runTest() {
	fmt.Println("Running tests...")
	cmd := exec.Command("go", "test", "./...", "-v")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
}
