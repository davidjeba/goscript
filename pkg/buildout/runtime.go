package buildout

import (
	"runtime"
	"time"
)

func runtimeGOOS() string {
	return runtime.GOOS
}

func nowRFC3339() string {
	return time.Now().UTC().Format(time.RFC3339)
}

