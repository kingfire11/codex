package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/kingfire11/codex/installer/internal/client"
	"github.com/kingfire11/codex/installer/internal/verify"
)

var version = "dev"

const defaultBaseURL = "https://api.designapi.ink/v1"
const defaultModel = "gpt-5.5"

func main() {
	var (
		token     string
		baseURL   string
		model     string
		clientArg string
		yes       bool
		showVer   bool
	)
	flag.StringVar(&token, "token", os.Getenv("OPENAI_API_KEY"), "API token (designapi.ink)")
	flag.StringVar(&baseURL, "base-url", defaultBaseURL, "API base URL")
	flag.StringVar(&model, "model", defaultModel, "default model")
	flag.StringVar(&clientArg, "client", "", "comma-separated: codex-cli, codex-vscode, codex-app, opencode, all")
	flag.BoolVar(&yes, "yes", false, "non-interactive: install without prompting")
	flag.BoolVar(&showVer, "version", false, "print version")
	flag.Parse()

	if showVer {
		fmt.Println("designapi-installer", version)
		return
	}

	cmd := "install"
	if flag.NArg() > 0 {
		cmd = flag.Arg(0)
	}

	switch cmd {
	case "install":
		runInstall(token, baseURL, model, clientArg, yes)
	case "doctor":
		runDoctor(token, baseURL, model)
	case "uninstall":
		runUninstall()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", cmd)
		fmt.Fprintln(os.Stderr, "usage: designapi-installer [install|doctor|uninstall] [flags]")
		os.Exit(2)
	}
}

func prompt(label, def string, secret bool) string {
	if def != "" {
		fmt.Printf("%s [%s]: ", label, mask(def, secret))
	} else {
		fmt.Printf("%s: ", label)
	}
	r := bufio.NewReader(os.Stdin)
	line, _ := r.ReadString('\n')
	line = strings.TrimSpace(line)
	if line == "" {
		return def
	}
	return line
}

func mask(s string, secret bool) string {
	if !secret || len(s) < 8 {
		return s
	}
	return s[:4] + "…" + s[len(s)-4:]
}

func runInstall(token, baseURL, model, clientArg string, yes bool) {
	if token == "" {
		token = prompt("Введите API-токен designapi.ink", "", true)
	}
	if token == "" {
		fmt.Fprintln(os.Stderr, "Токен обязателен.")
		os.Exit(1)
	}

	all := client.All()
	var selected []client.Installer
	if clientArg == "" {
		fmt.Println()
		fmt.Println("Найденные клиенты:")
		var detected []client.Installer
		for _, c := range all {
			ok, where := c.Detect()
			mark := "·"
			if ok {
				mark = "✓"
				detected = append(detected, c)
			}
			fmt.Printf("  %s  %-14s  %s\n", mark, c.Name(), where)
		}
		if len(detected) == 0 {
			fmt.Println("Ничего не найдено. Установите Codex CLI / VS Code Codex / OpenCode и запустите снова.")
			os.Exit(1)
		}
		fmt.Println()
		ans := prompt("Настроить все найденные? [Y/n]", "Y", false)
		if strings.HasPrefix(strings.ToLower(ans), "n") {
			fmt.Println("Отменено.")
			return
		}
		selected = detected
	} else if clientArg == "all" {
		selected = all
	} else {
		want := map[string]bool{}
		for _, n := range strings.Split(clientArg, ",") {
			want[strings.TrimSpace(n)] = true
		}
		for _, c := range all {
			if want[c.Name()] {
				selected = append(selected, c)
			}
		}
		if len(selected) == 0 {
			fmt.Fprintln(os.Stderr, "Не найдено клиентов под --client=", clientArg)
			os.Exit(1)
		}
	}

	for _, c := range selected {
		fmt.Printf("\n→ %s\n", c.Name())
		if err := c.Install(token, baseURL, model); err != nil {
			fmt.Fprintf(os.Stderr, "  ✗ %s: %v\n", c.Name(), err)
			continue
		}
		fmt.Printf("  ✓ %s настроен\n", c.Name())
	}

	fmt.Println("\n→ Проверяю связь с GPT через", baseURL, "...")
	res := verify.CheckChat(token, baseURL, model)
	res.Print(os.Stdout)
	if !res.OK {
		os.Exit(1)
	}
}

func runDoctor(token, baseURL, model string) {
	if token == "" {
		token = prompt("Введите API-токен для проверки", "", true)
	}
	fmt.Println("→ Проверяю", baseURL)
	res := verify.CheckChat(token, baseURL, model)
	res.Print(os.Stdout)
	if !res.OK {
		os.Exit(1)
	}
}

func runUninstall() {
	for _, c := range client.All() {
		if err := c.Uninstall(); err != nil {
			fmt.Fprintf(os.Stderr, "  ✗ %s: %v\n", c.Name(), err)
			continue
		}
		fmt.Printf("  ✓ %s\n", c.Name())
	}
	fmt.Println("Готово. Бэкапы сохранены в ~/.codex/backups/")
}
