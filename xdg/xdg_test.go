package xdg

import (
	"os/user"
	"testing"
)

func setup(t *testing.T, envs map[string]string) {
	for key, value := range envs {
		t.Setenv(key, value)
	}
}

func TestDataHome(t *testing.T) {
	usr := user.User{HomeDir: "/home/hashirafan"}
	x := Xdg{
		User: usr,
	}
	testCases := []struct {
		description string
		want        string
		env         map[string]string
	}{
		{"empty", "/home/hashirafan/.local/share", map[string]string{"XDG_DATA_HOME": ""}},
		{"special", "/home/hashirafan/opinionated/data", map[string]string{"XDG_DATA_HOME": "/home/hashirafan/opinionated/data"}},
		{"special and others", "/home/hashirafan/opinionated/data", map[string]string{
			"XDG_DATA_HOME":   "/home/hashirafan/opinionated/data",
			"XDG_CONFIG_HOME": "/home/hashirafan/.config",
			"XDG_STATE_HOME":  "/home/hashirafan/.local/state",
			"XDG_CACHE_HOME":  "/home/hashirafan/.cache",
			"XDG_DATA_DIRS":   "/usr/local/share:/usr/share",
		}},
		{"empty and others", "/home/hashirafan/.local/share", map[string]string{
			"XDG_DATA_HOME":   "",
			"XDG_CONFIG_HOME": "/home/hashirafan/.config",
			"XDG_STATE_HOME":  "/home/hashirafan/.local/state",
			"XDG_CACHE_HOME":  "/home/hashirafan/.cache",
			"XDG_DATA_DIRS":   "/usr/local/share:/usr/share",
		}},
	}

	for _, tc := range testCases {
		setup(t, tc.env)
		if got := x.DataHome(); got != tc.want {
			t.Fatalf("xdg returned wrong DataHome when %s. [got] %s [want] %s", tc.description, got, tc.want)
		}
	}
}
