#!/usr/bin/env bash
set -euo pipefail

API_KEY="${1:-${OPENAI_API_KEY:-__API_KEY__}}"
BASE_URL="__BASE_URL__/v1"
MODEL="__MODEL__"

CODEX_DIR="${HOME}/.codex"
AUTH_PATH="${CODEX_DIR}/auth.json"
CONFIG_PATH="${CODEX_DIR}/config.toml"
ENV_PATH="${CODEX_DIR}/designapi.env"
VSCODE_SERVER_ENV="${HOME}/.vscode-server/server-env-setup"
BACKUP_DIR="${CODEX_DIR}/backups/designapi-$(date +%Y%m%d-%H%M%S)"

mkdir -p "$CODEX_DIR"
if [[ -f "$CONFIG_PATH" || -f "$AUTH_PATH" ]]; then
  mkdir -p "$BACKUP_DIR"
  [[ -f "$CONFIG_PATH" ]] && cp "$CONFIG_PATH" "$BACKUP_DIR/"
  [[ -f "$AUTH_PATH"   ]] && cp "$AUTH_PATH"   "$BACKUP_DIR/"
  echo "Backups saved -> $BACKUP_DIR"
fi

cat > "$CONFIG_PATH" <<EOF
model = "$MODEL"
model_provider = "designapi"

[model_providers.designapi]
name = "DesignAPI"
base_url = "$BASE_URL"
wire_api = "responses"
env_key = "OPENAI_API_KEY"
EOF

cat > "$AUTH_PATH" <<EOF
{"OPENAI_API_KEY":"$API_KEY"}
EOF
chmod 600 "$AUTH_PATH"

cat > "$ENV_PATH" <<EOF
export OPENAI_API_KEY="$API_KEY"
export OPENAI_BASE_URL="$BASE_URL"
EOF
chmod 600 "$ENV_PATH"

# macOS: GUI-приложения (включая VS Code из Dock) не читают ~/.zshrc.
# Прокидываем env через launchctl + LaunchAgent, чтобы Codex-расширение увидело OPENAI_API_KEY.
if [[ "$(uname)" == "Darwin" ]]; then
  LA_DIR="${HOME}/Library/LaunchAgents"
  LA_PLIST="${LA_DIR}/ink.designapi.codex.plist"
  mkdir -p "$LA_DIR"
  cat > "$LA_PLIST" <<EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0"><dict>
  <key>Label</key><string>ink.designapi.codex</string>
  <key>ProgramArguments</key>
  <array>
    <string>/bin/sh</string><string>-c</string>
    <string>launchctl setenv OPENAI_API_KEY "$API_KEY"; launchctl setenv OPENAI_BASE_URL "$BASE_URL"</string>
  </array>
  <key>RunAtLoad</key><true/>
</dict></plist>
EOF
  launchctl unload "$LA_PLIST" 2>/dev/null || true
  launchctl load   "$LA_PLIST"
  launchctl setenv OPENAI_API_KEY  "$API_KEY"
  launchctl setenv OPENAI_BASE_URL "$BASE_URL"
fi

# VS Code remote (server) — пробрасываем env при подключении
if [[ -d "${HOME}/.vscode-server" ]]; then
  if [[ -f "$VSCODE_SERVER_ENV" ]]; then
    cp "$VSCODE_SERVER_ENV" "$BACKUP_DIR/" 2>/dev/null || true
  fi
  if ! grep -q "designapi.ink" "$VSCODE_SERVER_ENV" 2>/dev/null; then
    {
      echo "# designapi.ink"
      echo ". \"$ENV_PATH\""
    } >> "$VSCODE_SERVER_ENV"
    echo "VS Code Server env hook installed: $VSCODE_SERVER_ENV"
  fi
fi

# Локальный shell hook
SHELL_NAME="$(basename "${SHELL:-bash}")"
case "$SHELL_NAME" in
  zsh)  RC="$HOME/.zshrc" ;;
  bash) RC="$HOME/.bashrc" ;;
  fish) RC="$HOME/.config/fish/config.fish" ;;
  *)    RC="" ;;
esac
if [[ -n "$RC" ]] && ! grep -q "designapi.ink" "$RC" 2>/dev/null; then
  mkdir -p "$(dirname "$RC")"
  if [[ "$SHELL_NAME" == "fish" ]]; then
    printf '\n# designapi.ink\ntest -f %s; and source %s\n' "$ENV_PATH" "$ENV_PATH" >> "$RC"
  else
    printf '\n# designapi.ink\n[ -f %s ] && . %s\n' "$ENV_PATH" "$ENV_PATH" >> "$RC"
  fi
  echo "Added env hook to $RC"
fi

echo
echo "✅ VS Code Codex extension configured for designapi.ink"
echo
echo "ВАЖНО — иначе VS Code не увидит ключ:"
echo "  1) Полностью закройте VS Code:    osascript -e 'quit app \"Visual Studio Code\"'"
echo "     (или Cmd+Q, дождитесь чтобы пропал из Dock)"
echo "  2) Откройте VS Code этой командой (она передаст env в процесс):"
echo "     source \"$ENV_PATH\" && open -a 'Visual Studio Code'"
echo "  3) В панели Codex должно быть «Logged in with API key»."
