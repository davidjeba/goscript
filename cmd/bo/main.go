package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/davidjeba/goscript/pkg/buildout"
)

func main() {
	cfg, err := parseArgs(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		printUsage()
		os.Exit(1)
	}

	builder := buildout.NewBuilder()
	builder.DistDir = cfg.DistDir
	builder.GoBin = cfg.GoBin
	builder.DryRun = cfg.DryRun

	result, err := builder.Build(cfg.ManifestPath, cfg.Target)
	if err != nil {
		fmt.Fprintf(os.Stderr, "bo failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("bo %s -> %s\n", cfg.ManifestPath, cfg.Target)
	fmt.Printf("status: %s\n", result.Status)
	fmt.Printf("message: %s\n", result.Message)
	if result.OutputPath != "" {
		fmt.Printf("output: %s\n", result.OutputPath)
	}
	if result.BundlePath != "" {
		fmt.Printf("bundle: %s\n", result.BundlePath)
	}
}

type cliConfig struct {
	ManifestPath string
	Target       buildout.Target
	DistDir      string
	GoBin        string
	DryRun       bool
}

func parseArgs(args []string) (cliConfig, error) {
	cfg := cliConfig{
		DistDir: "dist",
		GoBin:   "go",
		Target:  buildout.TargetEXE,
	}

	var positional []string
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "":
			continue
		case "-":
			continue
		case "--dist", "-dist":
			if i+1 >= len(args) {
				return cfg, fmt.Errorf("missing value for %s", arg)
			}
			i++
			cfg.DistDir = args[i]
		case "--go", "-go":
			if i+1 >= len(args) {
				return cfg, fmt.Errorf("missing value for %s", arg)
			}
			i++
			cfg.GoBin = args[i]
		case "--dry-run":
			cfg.DryRun = true
		default:
			if strings.HasPrefix(arg, "-") {
				return cfg, fmt.Errorf("unknown flag: %s", arg)
			}
			positional = append(positional, arg)
		}
	}

	if len(positional) == 0 {
		return cfg, fmt.Errorf("missing manifest file")
	}

	cfg.ManifestPath = positional[0]
	if len(positional) >= 2 {
		target, err := buildout.ParseTarget(positional[len(positional)-1])
		if err != nil {
			return cfg, err
		}
		cfg.Target = target
	}

	return cfg, nil
}

func printUsage() {
	fmt.Println(strings.TrimSpace(`
bo - build out

Usage:
  bo <manifest> [target]
  bo <manifest> - <target>
  bo [--dist dist] [--go go] [--dry-run] <manifest> [target]

Targets:
  exe   Build a host executable
  goe   Build a portable GOE bundle
  apk   Generate an Android packaging scaffold
  ipa   Generate an iOS packaging scaffold
  dmg   Generate a macOS packaging scaffold

Examples:
  bo admin.manifest - exe
  bo admin.manifest exe
  bo admin.manifest - goe
  bo calc.manifest - exe
`))
}

