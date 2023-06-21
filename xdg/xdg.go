package xdg

import (
	"os"
	"os/user"
	"path/filepath"
)

type Xdg struct {
	User user.User
}

func (x *Xdg) DataHome() string {
	fromEnv := os.Getenv("XDG_DATA_HOME")
	if fromEnv == "" {
		return x.DefaultDataHome()
	}

	return fromEnv
}

func (x *Xdg) DefaultDataHome() string {
	return filepath.Join(x.User.HomeDir, ".local", "share")
}
