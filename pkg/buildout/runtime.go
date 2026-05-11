package buildout

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

func runtimeGOOS() string {
	return runtime.GOOS
}

func nowRFC3339() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func findModuleRoot(startDir string) (string, error) {
	dir := startDir
	for {
		modPath := filepath.Join(dir, "go.mod")
		if info, err := os.Stat(modPath); err == nil && !info.IsDir() {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("unable to locate go.mod from %s", startDir)
		}

		dir = parent
	}
}
