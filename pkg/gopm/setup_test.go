package gopm

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/davidjeba/goscript/pkg/buildout"
)

func TestParseSetupArgsDefaults(t *testing.T) {
	opts, err := parseSetupArgs([]string{"demo-app"})
	if err != nil {
		t.Fatalf("parseSetupArgs returned error: %v", err)
	}

	if opts.Mode != "cs" {
		t.Fatalf("expected mode cs, got %q", opts.Mode)
	}

	if opts.Type != "app" {
		t.Fatalf("expected type app, got %q", opts.Type)
	}

	if filepath.Base(opts.ProjectDir) != "demo-app" {
		t.Fatalf("expected project dir to end with demo-app, got %q", opts.ProjectDir)
	}
}

func TestParseSetupArgsSwarmERP(t *testing.T) {
	opts, err := parseSetupArgs([]string{"--sw", "--type", "erp", "mesh-suite"})
	if err != nil {
		t.Fatalf("parseSetupArgs returned error: %v", err)
	}

	if opts.Mode != "sw" {
		t.Fatalf("expected mode sw, got %q", opts.Mode)
	}

	if opts.Type != "erp" {
		t.Fatalf("expected type erp, got %q", opts.Type)
	}
}

func TestSetupProjectWritesManifest(t *testing.T) {
	root := t.TempDir()
	projectDir := filepath.Join(root, "erp-suite")

	pm := NewPackageManager()
	manifestPath, err := pm.setupProject(SetupOptions{
		ProjectDir:   projectDir,
		ProjectName:  "erp-suite",
		Mode:         "sw",
		Type:         "erp",
		Entrypoint:   "./cmd/server",
		ManifestName: "erp-suite",
	})
	if err != nil {
		t.Fatalf("setupProject returned error: %v", err)
	}

	raw, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatalf("failed to read manifest: %v", err)
	}

	var manifest buildout.Manifest
	if err := json.Unmarshal(raw, &manifest); err != nil {
		t.Fatalf("failed to decode manifest: %v", err)
	}

	if manifest.Mode != "sw" {
		t.Fatalf("expected mode sw, got %q", manifest.Mode)
	}

	if manifest.BaseDir != "base" {
		t.Fatalf("expected base dir base, got %q", manifest.BaseDir)
	}

	if manifest.AgentsDir != "agents" {
		t.Fatalf("expected agents dir agents, got %q", manifest.AgentsDir)
	}

	requiredDirs := []string{
		filepath.Join(projectDir, "base"),
		filepath.Join(projectDir, "agents"),
		filepath.Join(projectDir, "app", "swarm-policies"),
		filepath.Join(projectDir, "manifests"),
	}

	for _, dir := range requiredDirs {
		if info, err := os.Stat(dir); err != nil || !info.IsDir() {
			t.Fatalf("expected directory %s to exist", dir)
		}
	}
}
