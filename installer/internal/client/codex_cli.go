package client

import (
	"os/exec"

	"github.com/kingfire11/codex/installer/internal/shell"
)

type codexCLI struct{}

func (codexCLI) Name() string { return "codex-cli" }

func (codexCLI) Detect() (bool, string) {
	p, err := exec.LookPath("codex")
	if err != nil {
		return false, "не найден в PATH"
	}
	return true, p
}

func (c *codexCLI) Install(token, baseURL, model string) error {
	if err := writeCodexConfig(token, baseURL, model); err != nil {
		return err
	}
	if !isWindows() {
		_ = shell.AddEnvHook(codexDir() + "/designapi.env")
	}
	return nil
}

func (codexCLI) Uninstall() error { return removeCodexFiles() }
