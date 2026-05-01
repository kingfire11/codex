package shell

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const marker = "# designapi.ink"

func detectRC() (string, string) {
	h, _ := os.UserHomeDir()
	sh := filepath.Base(os.Getenv("SHELL"))
	switch sh {
	case "zsh":
		return sh, filepath.Join(h, ".zshrc")
	case "fish":
		return sh, filepath.Join(h, ".config", "fish", "config.fish")
	case "bash":
		return sh, filepath.Join(h, ".bashrc")
	}
	if _, err := os.Stat(filepath.Join(h, ".zshrc")); err == nil {
		return "zsh", filepath.Join(h, ".zshrc")
	}
	return "bash", filepath.Join(h, ".bashrc")
}

func AddEnvHook(envFile string) error {
	sh, rc := detectRC()
	hook := fmt.Sprintf("\n%s\n[ -f %q ] && . %q\n", marker, envFile, envFile)
	if sh == "fish" {
		hook = fmt.Sprintf("\n%s\ntest -f %q; and source %q\n", marker, envFile, envFile)
	}
	return AppendOnce(rc, hook)
}

func AppendOnce(path, content string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	if data, err := os.ReadFile(path); err == nil && strings.Contains(string(data), marker) {
		return nil
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(content)
	return err
}
