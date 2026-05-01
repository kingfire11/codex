package client

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/kingfire11/codex/installer/internal/backup"
)

func codexDir() string {
	h, _ := os.UserHomeDir()
	return filepath.Join(h, ".codex")
}

func writeCodexConfig(token, baseURL, model string) error {
	dir := codexDir()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	if err := backup.SnapshotIfExists(dir, []string{"config.toml", "auth.json"}); err != nil {
		return err
	}

	cfg := fmt.Sprintf(`model = %q
model_provider = "designapi"

[model_providers.designapi]
name = "DesignAPI"
base_url = %q
wire_api = "responses"
env_key = "OPENAI_API_KEY"
`, model, baseURL)
	if err := os.WriteFile(filepath.Join(dir, "config.toml"), []byte(cfg), 0o644); err != nil {
		return err
	}

	auth := fmt.Sprintf(`{"OPENAI_API_KEY":%q}`+"\n", token)
	if err := os.WriteFile(filepath.Join(dir, "auth.json"), []byte(auth), 0o600); err != nil {
		return err
	}

	envFile := fmt.Sprintf("export OPENAI_API_KEY=%q\nexport OPENAI_BASE_URL=%q\n", token, baseURL)
	if err := os.WriteFile(filepath.Join(dir, "designapi.env"), []byte(envFile), 0o600); err != nil {
		return err
	}
	return nil
}

func removeCodexFiles() error {
	dir := codexDir()
	for _, f := range []string{"config.toml", "auth.json", "designapi.env"} {
		_ = os.Remove(filepath.Join(dir, f))
	}
	return nil
}

func isWindows() bool { return runtime.GOOS == "windows" }
