package client

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kingfire11/codex/installer/internal/backup"
)

type openCode struct{}

func (openCode) Name() string { return "opencode" }

func ocConfigPath() string {
	h, _ := os.UserHomeDir()
	return filepath.Join(h, ".config", "opencode", "opencode.json")
}

func (openCode) Detect() (bool, string) {
	dir := filepath.Dir(ocConfigPath())
	if _, err := os.Stat(dir); err == nil {
		return true, dir
	}
	return false, "OpenCode не настроен (~/.config/opencode не существует)"
}

func (o *openCode) Install(token, baseURL, model string) error {
	cfgPath := ocConfigPath()
	dir := filepath.Dir(cfgPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	if err := backup.SnapshotIfExists(dir, []string{"opencode.json"}); err != nil {
		return err
	}
	cfg := fmt.Sprintf(`{
  "$schema": "https://opencode.ai/config.json",
  "provider": {
    "designapi": {
      "npm": "@ai-sdk/openai-compatible",
      "name": "DesignAPI",
      "options": {
        "baseURL": %q,
        "apiKey": %q
      },
      "models": {
        %q: { "name": %q }
      }
    }
  }
}
`, baseURL, token, model, model)
	return os.WriteFile(cfgPath, []byte(cfg), 0o600)
}

func (openCode) Uninstall() error {
	_ = os.Remove(ocConfigPath())
	return nil
}
