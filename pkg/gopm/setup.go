package gopm

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/davidjeba/goscript/pkg/buildout"
)

// SetupOptions captures the topology and shape requested for gopm setup.
type SetupOptions struct {
	ProjectDir   string
	ProjectName  string
	Mode         string
	Type         string
	Entrypoint   string
	ManifestName string
	Force        bool
}

func parseSetupArgs(args []string) (SetupOptions, error) {
	opts := SetupOptions{
		Mode: "cs",
		Type: "app",
	}

	for i := 0; i < len(args); i++ {
		arg := strings.TrimSpace(args[i])
		if arg == "" {
			continue
		}

		switch arg {
		case "--cs":
			opts.Mode = "cs"
		case "--sw":
			opts.Mode = "sw"
		case "--force":
			opts.Force = true
		case "--mode":
			i++
			if i >= len(args) {
				return SetupOptions{}, fmt.Errorf("missing value for --mode")
			}
			opts.Mode = strings.ToLower(strings.TrimSpace(args[i]))
		case "--type", "--template":
			i++
			if i >= len(args) {
				return SetupOptions{}, fmt.Errorf("missing value for %s", arg)
			}

			projectType, err := normalizeProjectType(args[i])
			if err != nil {
				return SetupOptions{}, err
			}
			opts.Type = projectType
		case "--entrypoint":
			i++
			if i >= len(args) {
				return SetupOptions{}, fmt.Errorf("missing value for --entrypoint")
			}
			opts.Entrypoint = strings.TrimSpace(args[i])
		case "--manifest":
			i++
			if i >= len(args) {
				return SetupOptions{}, fmt.Errorf("missing value for --manifest")
			}
			opts.ManifestName = strings.TrimSpace(args[i])
		default:
			if strings.HasPrefix(arg, "--") {
				return SetupOptions{}, fmt.Errorf("unknown setup flag %q", arg)
			}
			if opts.ProjectDir != "" {
				return SetupOptions{}, fmt.Errorf("unexpected extra argument %q", arg)
			}
			opts.ProjectDir = arg
		}
	}

	if opts.Mode != "cs" && opts.Mode != "sw" {
		return SetupOptions{}, fmt.Errorf("mode must be either cs or sw")
	}

	if opts.ProjectDir == "" {
		opts.ProjectDir = "."
	}

	projectDir, err := filepath.Abs(opts.ProjectDir)
	if err != nil {
		return SetupOptions{}, fmt.Errorf("resolve project path: %w", err)
	}
	opts.ProjectDir = projectDir

	opts.ProjectName = sanitizeName(filepath.Base(projectDir))
	if opts.ProjectName == "" {
		opts.ProjectName = "goscript-app"
	}

	if opts.Entrypoint == "" {
		opts.Entrypoint = "./cmd/server"
	}

	if opts.ManifestName == "" {
		opts.ManifestName = opts.ProjectName
	}

	return opts, nil
}

func normalizeProjectType(raw string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "", "app":
		return "app", nil
	case "web", "website":
		return "website", nil
	case "erp":
		return "erp", nil
	default:
		return "", fmt.Errorf("project type must be website, app, or erp")
	}
}

func sanitizeName(raw string) string {
	raw = strings.TrimSpace(strings.ToLower(raw))
	if raw == "" {
		return ""
	}

	var b strings.Builder
	lastDash := false
	for _, r := range raw {
		switch {
		case r >= 'a' && r <= 'z':
			b.WriteRune(r)
			lastDash = false
		case r >= '0' && r <= '9':
			b.WriteRune(r)
			lastDash = false
		default:
			if !lastDash {
				b.WriteByte('-')
				lastDash = true
			}
		}
	}

	return strings.Trim(b.String(), "-")
}

func (pm *PackageManager) setupProject(opts SetupOptions) (string, error) {
	dirs := []string{
		opts.ProjectDir,
		filepath.Join(opts.ProjectDir, "base"),
		filepath.Join(opts.ProjectDir, "base", "config"),
		filepath.Join(opts.ProjectDir, "base", "policies"),
		filepath.Join(opts.ProjectDir, "agents"),
		filepath.Join(opts.ProjectDir, "app"),
		filepath.Join(opts.ProjectDir, "app", "modules"),
		filepath.Join(opts.ProjectDir, "app", "pages"),
		filepath.Join(opts.ProjectDir, "app", "components"),
		filepath.Join(opts.ProjectDir, "app", "services"),
		filepath.Join(opts.ProjectDir, "app", "routes"),
		filepath.Join(opts.ProjectDir, "app", "assets"),
		filepath.Join(opts.ProjectDir, "core"),
		filepath.Join(opts.ProjectDir, "tests"),
		filepath.Join(opts.ProjectDir, "docs"),
		filepath.Join(opts.ProjectDir, "deploy"),
		filepath.Join(opts.ProjectDir, "manifests"),
		filepath.Join(opts.ProjectDir, "cmd"),
		filepath.Join(opts.ProjectDir, "cmd", "server"),
	}

	if opts.Mode == "cs" {
		dirs = append(dirs,
			filepath.Join(opts.ProjectDir, "app", "api"),
			filepath.Join(opts.ProjectDir, "app", "controllers"),
			filepath.Join(opts.ProjectDir, "app", "views"),
		)
	} else {
		dirs = append(dirs,
			filepath.Join(opts.ProjectDir, "app", "topology"),
			filepath.Join(opts.ProjectDir, "app", "sync"),
			filepath.Join(opts.ProjectDir, "app", "swarm-policies"),
			filepath.Join(opts.ProjectDir, "app", "trust"),
		)
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return "", fmt.Errorf("create %s: %w", dir, err)
		}
	}

	if err := writeSetupStub(filepath.Join(opts.ProjectDir, "base", "README.md"), baseReadmeStub(opts), opts.Force); err != nil {
		return "", err
	}

	if err := writeSetupStub(filepath.Join(opts.ProjectDir, "agents", "README.md"), agentsReadmeStub(), opts.Force); err != nil {
		return "", err
	}

	manifest := buildout.Manifest{
		Name:        opts.ManifestName,
		Mode:        opts.Mode,
		Output:      opts.ManifestName,
		Module:      ".",
		Entrypoint:  opts.Entrypoint,
		BaseDir:     "base",
		AgentsDir:   "agents",
		Description: fmt.Sprintf("%s scaffold generated by gopm setup", opts.Type),
		Pages:       defaultPages(opts.Type),
		Metadata: map[string]string{
			"projectType": opts.Type,
			"setupMode":   opts.Mode,
		},
	}

	manifestPath := filepath.Join(opts.ProjectDir, "manifests", opts.ManifestName+".manifest")
	if err := writeManifestFile(manifestPath, manifest, opts.Force); err != nil {
		return "", err
	}

	return manifestPath, nil
}

func defaultPages(projectType string) []string {
	switch projectType {
	case "website":
		return []string{"/"}
	case "erp":
		return []string{"/", "/dashboard", "/modules"}
	default:
		return []string{"/"}
	}
}

func writeSetupStub(path, contents string, force bool) error {
	if !force {
		if _, err := os.Stat(path); err == nil {
			return nil
		} else if !os.IsNotExist(err) {
			return fmt.Errorf("check %s: %w", path, err)
		}
	}

	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

func writeManifestFile(path string, manifest buildout.Manifest, force bool) error {
	if !force {
		if _, err := os.Stat(path); err == nil {
			return fmt.Errorf("manifest already exists: %s (use --force to overwrite)", path)
		} else if !os.IsNotExist(err) {
			return fmt.Errorf("check %s: %w", path, err)
		}
	}

	manifest.Normalize(path)

	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal manifest: %w", err)
	}

	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write manifest %s: %w", path, err)
	}

	return nil
}

func baseReadmeStub(opts SetupOptions) string {
	return fmt.Sprintf("# Project Base Guidance\n\nThis project is scaffolded in `%s` mode as a `%s` project.\n\nUse this folder for project-local AI guidance that extends the shared GoScript `base/` contract.\n", opts.Mode, opts.Type)
}

func agentsReadmeStub() string {
	return "# Runtime Agents\n\nUse this folder only for autonomous roles that exist inside the application at runtime.\n"
}
