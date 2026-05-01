# designapi-installer

Кросс-платформенный инсталлер для подключения Codex CLI / VS Code Codex / Codex App / OpenCode к API `https://api.designapi.ink`.

## Использование

```
designapi-installer install --token=YOUR_API_KEY            # настроить найденные клиенты
designapi-installer install --token=... --client=codex-cli  # только один
designapi-installer install --client=all --yes              # неинтерактивно
designapi-installer doctor --token=...                      # только проверить связь
designapi-installer uninstall                               # снести наши конфиги
```

Без аргументов запускается интерактивный режим: спросит токен, покажет найденные клиенты и предложит настроить все сразу.

## Что делает

- Делает бэкап `~/.codex/config.toml`, `~/.codex/auth.json` и `~/.config/opencode/opencode.json` в `~/.codex/backups/designapi-<timestamp>/`.
- Пишет конфиги с `base_url = https://api.designapi.ink/v1`.
- Прописывает env-хук в `~/.zshrc`/`~/.bashrc`/`config.fish` (чтобы `OPENAI_API_KEY` подхватывался).
- На macOS — ставит LaunchAgent, чтобы переменные подхватывал и Codex App.
- В конце делает `POST /v1/chat/completions` с одним токеном и классифицирует ошибку, если не получилось (DNS / TLS / 401 / 429 / 5xx / таймаут).

## Сборка

```
go build -trimpath -ldflags "-s -w" -o designapi-installer .
```

CI собирает релизы под linux/amd64, linux/arm64, darwin/amd64, darwin/arm64, windows/amd64 при push тега `v*`.
