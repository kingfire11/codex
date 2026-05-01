package client

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

type codexApp struct{}

func (codexApp) Name() string { return "codex-app" }

func (codexApp) Detect() (bool, string) {
	switch runtime.GOOS {
	case "darwin":
		for _, p := range []string{"/Applications/Codex.app", filepath.Join(os.Getenv("HOME"), "Applications", "Codex.app")} {
			if _, err := os.Stat(p); err == nil {
				return true, p
			}
		}
	case "windows":
		for _, p := range []string{
			filepath.Join(os.Getenv("LOCALAPPDATA"), "Programs", "Codex"),
			filepath.Join(os.Getenv("ProgramFiles"), "Codex"),
		} {
			if p == "" {
				continue
			}
			if _, err := os.Stat(p); err == nil {
				return true, p
			}
		}
	}
	return false, "Codex App не найден"
}

func (c *codexApp) Install(token, baseURL, model string) error {
	if err := writeCodexConfig(token, baseURL, model); err != nil {
		return err
	}
	if runtime.GOOS == "darwin" {
		return installLaunchAgent(token, baseURL)
	}
	if runtime.GOOS == "windows" {
		setUserEnv("OPENAI_API_KEY", token)
		setUserEnv("OPENAI_BASE_URL", baseURL)
	}
	return nil
}

func (codexApp) Uninstall() error {
	if runtime.GOOS == "darwin" {
		h, _ := os.UserHomeDir()
		p := filepath.Join(h, "Library", "LaunchAgents", "ink.designapi.codex.plist")
		_ = exec.Command("launchctl", "unload", p).Run()
		_ = os.Remove(p)
	}
	return removeCodexFiles()
}

func installLaunchAgent(token, baseURL string) error {
	h, _ := os.UserHomeDir()
	dir := filepath.Join(h, "Library", "LaunchAgents")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	plistPath := filepath.Join(dir, "ink.designapi.codex.plist")
	plist := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0"><dict>
  <key>Label</key><string>ink.designapi.codex</string>
  <key>ProgramArguments</key>
  <array>
    <string>/bin/sh</string><string>-c</string>
    <string>launchctl setenv OPENAI_API_KEY %q; launchctl setenv OPENAI_BASE_URL %q</string>
  </array>
  <key>RunAtLoad</key><true/>
</dict></plist>
`, token, baseURL)
	if err := os.WriteFile(plistPath, []byte(plist), 0o644); err != nil {
		return err
	}
	_ = exec.Command("launchctl", "unload", plistPath).Run()
	_ = exec.Command("launchctl", "load", plistPath).Run()
	_ = exec.Command("launchctl", "setenv", "OPENAI_API_KEY", token).Run()
	_ = exec.Command("launchctl", "setenv", "OPENAI_BASE_URL", baseURL).Run()
	return nil
}

func setUserEnv(name, value string) {
	if runtime.GOOS != "windows" {
		return
	}
	_ = exec.Command("setx", name, value).Run()
}
