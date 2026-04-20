package buildout

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseTarget(t *testing.T) {
	target, err := ParseTarget("exe")
	if err != nil {
		t.Fatalf("ParseTarget returned error: %v", err)
	}
	if target != TargetEXE {
		t.Fatalf("expected %q, got %q", TargetEXE, target)
	}
}

func TestLoadManifest(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "admin.manifest")

	content := []byte(`{
  "name": "admin",
  "mode": "sw",
  "module": ".",
  "entrypoint": "./cmd/server",
  "paths": ["/admin", "/admin/users"]
}`)

	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatalf("failed to write manifest: %v", err)
	}

	manifest, err := LoadManifest(path)
	if err != nil {
		t.Fatalf("LoadManifest returned error: %v", err)
	}

	if manifest.Name != "admin" {
		t.Fatalf("expected name admin, got %q", manifest.Name)
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

	if manifest.BuildTarget() != "./cmd/server" {
		t.Fatalf("expected build target ./cmd/server, got %q", manifest.BuildTarget())
	}
}
