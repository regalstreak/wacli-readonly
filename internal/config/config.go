package config

import (
	"os"
	"path/filepath"
	"runtime"
)

// EnvStoreDir is the environment variable that overrides the default store
// directory. This is useful for Docker, CI, and multi-tenant deployments
// where the store path needs to be configured without passing --store on
// every invocation.
const EnvStoreDir = "WACLI_STORE_DIR"

// DefaultStoreDir returns the store directory to use when --store is not
// supplied. It checks WACLI_STORE_DIR first, then falls back to the XDG state
// directory on Linux or ~/.wacli-readonly on other platforms.
func DefaultStoreDir() string {
	if dir := os.Getenv(EnvStoreDir); dir != "" {
		return dir
	}
	xdgStateHome := os.Getenv("XDG_STATE_HOME")
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		if runtime.GOOS == "linux" && xdgStateHome != "" {
			return filepath.Join(xdgStateHome, "wacli-readonly")
		}
		return ".wacli-readonly"
	}
	return defaultStoreDirFor(runtime.GOOS, home, xdgStateHome, pathExists)
}

func defaultStoreDirFor(goos, home, xdgStateHome string, exists func(string) bool) string {
	legacy := filepath.Join(home, ".wacli-readonly")
	if goos != "linux" {
		return legacy
	}
	if xdgStateHome != "" {
		return filepath.Join(xdgStateHome, "wacli-readonly")
	}
	xdgDefault := filepath.Join(home, ".local", "state", "wacli-readonly")
	if exists(legacy) && !exists(xdgDefault) {
		return legacy
	}
	return xdgDefault
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
