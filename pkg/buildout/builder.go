package buildout

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// Builder executes build-out exports from a manifest.
type Builder struct {
	GoBin   string
	DistDir string
	DryRun  bool
	Stdout  io.Writer
	Stderr  io.Writer
}

// NewBuilder creates a builder with sensible defaults.
func NewBuilder() *Builder {
	return &Builder{
		GoBin:   "go",
		DistDir: "dist",
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
	}
}

// Build exports the requested target from a manifest.
func (b *Builder) Build(manifestPath string, target Target) (BuildResult, error) {
	manifest, err := LoadManifest(manifestPath)
	if err != nil {
		return BuildResult{}, err
	}

	distDir := b.DistDir
	if distDir == "" {
		distDir = "dist"
	}

	artifactRoot := manifest.PlannedOutputDir(distDir)
	if err := os.MkdirAll(artifactRoot, 0o755); err != nil {
		return BuildResult{}, err
	}

	plan := manifest.Plan(manifestPath, target, distDir)
	if err := writeJSONFile(filepath.Join(artifactRoot, "build-plan.json"), plan); err != nil {
		return BuildResult{}, err
	}
	if err := writeJSONFile(filepath.Join(artifactRoot, "manifest.normalized.json"), manifest); err != nil {
		return BuildResult{}, err
	}

	switch target {
	case TargetEXE:
		return b.buildExecutable(manifestPath, manifest, target, artifactRoot)
	case TargetGOE:
		return b.buildGOE(manifestPath, manifest, target, artifactRoot, plan)
	case TargetAPK, TargetIPA, TargetDMG:
		return b.scaffoldPlatformBundle(manifestPath, manifest, target, artifactRoot)
	default:
		return BuildResult{}, fmt.Errorf("unsupported target %q", target)
	}
}

func (b *Builder) buildExecutable(manifestPath string, manifest Manifest, target Target, artifactRoot string) (BuildResult, error) {
	outputPath := filepath.Join(artifactRoot, manifest.BinaryName())
	if b.DryRun {
		return BuildResult{
			ManifestPath: manifestPath,
			Manifest:     manifest,
			Target:       target,
			BuildTarget:  manifest.BuildTarget(),
			Status:       "dry-run",
			Message:      "dry run complete; executable not built",
			OutputPath:   outputPath,
			Artifacts:    []string{filepath.Join(artifactRoot, "build-plan.json"), filepath.Join(artifactRoot, "manifest.normalized.json")},
		}, nil
	}

	if err := b.runGoBuild(filepath.Dir(manifestPath), manifest.BuildTarget(), outputPath); err != nil {
		return BuildResult{}, err
	}

	return BuildResult{
		ManifestPath: manifestPath,
		Manifest:     manifest,
		Target:       target,
		BuildTarget:  manifest.BuildTarget(),
		Status:       "built",
		Message:      "executable built successfully",
		OutputPath:   outputPath,
		Artifacts: []string{
			outputPath,
			filepath.Join(artifactRoot, "build-plan.json"),
			filepath.Join(artifactRoot, "manifest.normalized.json"),
		},
	}, nil
}

func (b *Builder) buildGOE(manifestPath string, manifest Manifest, target Target, artifactRoot string, plan BuildPlan) (BuildResult, error) {
	tempDir, err := os.MkdirTemp("", "bo-goe-*")
	if err != nil {
		return BuildResult{}, err
	}
	defer os.RemoveAll(tempDir)

	binaryPath := filepath.Join(tempDir, manifest.BinaryName())
	if runtime.GOOS != "windows" && strings.HasSuffix(strings.ToLower(binaryPath), ".exe") {
		binaryPath = strings.TrimSuffix(binaryPath, ".exe")
	}

	if b.DryRun {
		return BuildResult{
			ManifestPath: manifestPath,
			Manifest:     manifest,
			Target:       target,
			BuildTarget:  manifest.BuildTarget(),
			Status:       "dry-run",
			Message:      "dry run complete; GOE bundle not built",
			OutputPath:   binaryPath,
		}, nil
	}

	if err := b.runGoBuild(filepath.Dir(manifestPath), manifest.BuildTarget(), binaryPath); err != nil {
		return BuildResult{}, err
	}

	bundlePath := filepath.Join(artifactRoot, manifest.Output+".goe")
	if err := b.writeGOEBundle(bundlePath, manifestPath, manifest, binaryPath, plan); err != nil {
		return BuildResult{}, err
	}

	return BuildResult{
		ManifestPath: manifestPath,
		Manifest:     manifest,
		Target:       target,
		BuildTarget:  manifest.BuildTarget(),
		Status:       "built",
		Message:      "goe bundle created successfully",
		OutputPath:   binaryPath,
		BundlePath:   bundlePath,
		Artifacts: []string{
			bundlePath,
			filepath.Join(artifactRoot, "build-plan.json"),
			filepath.Join(artifactRoot, "manifest.normalized.json"),
		},
	}, nil
}

func (b *Builder) scaffoldPlatformBundle(manifestPath string, manifest Manifest, target Target, artifactRoot string) (BuildResult, error) {
	platformRoot := filepath.Join(artifactRoot, strings.ToLower(string(target)))
	if err := os.MkdirAll(platformRoot, 0o755); err != nil {
		return BuildResult{}, err
	}

	readme := fmt.Sprintf(
		"BO generated a %s scaffold for %s.\n\n"+
			"Manifest: %s\n"+
			"Build target: %s\n\n"+
			"This artifact is a build contract, not a finalized mobile/desktop packager yet.\n",
		strings.ToUpper(string(target)),
		manifest.Name,
		manifestPath,
		manifest.BuildTarget(),
	)

	if err := os.WriteFile(filepath.Join(platformRoot, "README.txt"), []byte(readme), 0o644); err != nil {
		return BuildResult{}, err
	}
	if err := writeJSONFile(filepath.Join(platformRoot, "manifest.json"), manifest); err != nil {
		return BuildResult{}, err
	}

	return BuildResult{
		ManifestPath: manifestPath,
		Manifest:     manifest,
		Target:       target,
		BuildTarget:  manifest.BuildTarget(),
		Status:       "scaffolded",
		Message:      fmt.Sprintf("%s scaffold generated; native packaging can be layered on later", strings.ToUpper(string(target))),
		OutputPath:   platformRoot,
		Artifacts: []string{
			platformRoot,
			filepath.Join(artifactRoot, "build-plan.json"),
			filepath.Join(artifactRoot, "manifest.normalized.json"),
		},
	}, nil
}

func (b *Builder) runGoBuild(workDir, buildTarget, outputPath string) error {
	if err := os.MkdirAll(filepath.Dir(outputPath), 0o755); err != nil {
		return err
	}

	args := []string{"build", "-trimpath", "-o", outputPath, "-ldflags", "-s -w", buildTarget}
	cmd := exec.Command(b.GoBin, args...)
	cmd.Dir = workDir
	cmd.Stdout = b.Stdout
	cmd.Stderr = b.Stderr
	cmd.Env = os.Environ()

	return cmd.Run()
}

func (b *Builder) writeGOEBundle(bundlePath, manifestPath string, manifest Manifest, binaryPath string, plan BuildPlan) error {
	if err := os.MkdirAll(filepath.Dir(bundlePath), 0o755); err != nil {
		return err
	}

	file, err := os.Create(bundlePath)
	if err != nil {
		return err
	}
	defer file.Close()

	zipWriter := zip.NewWriter(file)
	defer zipWriter.Close()

	if err := addFileToZip(zipWriter, "manifest.json", mustMarshal(manifest)); err != nil {
		return err
	}

	plan.Notes = []string{"portable GOE bundle with embedded binary and manifest"}
	if err := addFileToZip(zipWriter, "build-plan.json", mustMarshal(plan)); err != nil {
		return err
	}

	if err := addBinaryToZip(zipWriter, filepath.Join("bin", manifest.BinaryName()), binaryPath); err != nil {
		return err
	}

	runtimeNote := fmt.Sprintf(
		"GOE bundle for %s\n\n"+
			"Target: %s\n"+
			"Build target: %s\n"+
			"Created: %s\n\n"+
			"This package is a portable app bundle format. The Go runtime is already compiled into the binary itself.\n"+
			"Future runtime loaders can add sandboxing, launch metadata, or device-specific adapters.\n",
		manifest.Name,
		TargetGOE,
		manifest.BuildTarget(),
		time.Now().UTC().Format(time.RFC3339),
	)
	if err := addFileToZip(zipWriter, "README.txt", []byte(runtimeNote)); err != nil {
		return err
	}

	return nil
}

func addFileToZip(zipWriter *zip.Writer, name string, content []byte) error {
	writer, err := zipWriter.Create(name)
	if err != nil {
		return err
	}

	_, err = writer.Write(content)
	return err
}

func addBinaryToZip(zipWriter *zip.Writer, name, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer, err := zipWriter.Create(name)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, file)
	return err
}

func mustMarshal(value interface{}) []byte {
	data, err := jsonMarshal(value)
	if err != nil {
		return []byte("{}")
	}
	return data
}
