package client

import (
	"os"
	"path/filepath"

	"github.com/kingfire11/codex/installer/internal/shell"
)

type codexVSCode struct{}

func (codexVSCode) Name() string { return "codex-vscode" }

func (codexVSCode) Detect() (bool, string) {
	h, _ := os.UserHomeDir()
	candidates := []string{
		filepath.Join(h, ".vscode", "extensions"),
		filepath.Join(h, ".vscode-server", "extensions"),
		filepath.Join(h, "Library", "Application Support", "Code", "User"),
		filepath.Join(h, "AppData", "Roaming", "Code", "User"),
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return true, p
		}
	}
	return false, "VS Code не обнаружен"
}

func (c *codexVSCode) Install(token, baseURL, model string) error {
	if err := writeCodexConfig(token, baseURL, model); err != nil {
		return err
	}
	if !isWindows() {
		envPath := codexDir() + "/designapi.env"
		_ = shell.AddEnvHook(envPath)
		// VS Code Server (remote)
		h, _ := os.UserHomeDir()
		serverEnv := filepath.Join(h, ".vscode-server", "server-env-setup")
		if _, err := os.Stat(filepath.Dir(serverEnv)); err == nil {
			_ = shell.AppendOnce(serverEnv, "# designapi.ink\n. \""+envPath+"\"\n")
		}
	}
	return nil
}

func (codexVSCode) Uninstall() error { return removeCodexFiles() }
