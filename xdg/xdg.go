package xdg

import (
	"os"
	"path/filepath"
)

func DataHome() (string, error) {
	fromEnv := os.Getenv("XDG_DATA_HOME")
	if fromEnv == "" {
		return DefaultDataHome()
	}

	return fromEnv, nil
}

func DefaultDataHome() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".local", "share"), nil
}
