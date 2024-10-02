package config

import (
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

// Returns path of the config file based on the OS
func GetConfigFilePath() string {
	var configDir string

	switch runtime.GOOS {
	case "windows":
		configDir = os.Getenv("APPDATA")
	case "darwin":
		usr, _ := user.Current()
		configDir = filepath.Join(usr.HomeDir, "Library", "Application Support")
	case "linux":
		usr, _ := user.Current()
		configDir = filepath.Join(usr.HomeDir, ".config")
	default:
		configDir = "."
	}

	return filepath.Join(configDir, "github.com/Debug-Studios/kinsyn", "config.yaml")
}
