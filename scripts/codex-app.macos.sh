#!/usr/bin/env bash
set -euo pipefail

API_KEY="${1:-${OPENAI_API_KEY:-__API_KEY__}}"
BASE_URL="__BASE_URL__/v1"
MODEL="__MODEL__"

CODEX_DIR="${HOME}/.codex"
AUTH_PATH="${CODEX_DIR}/auth.json"
CONFIG_PATH="${CODEX_DIR}/config.toml"
ENV_PATH="${CODEX_DIR}/designapi.env"
LAUNCH_AGENT_DIR="${HOME}/Library/LaunchAgents"
LAUNCH_AGENT="${LAUNCH_AGENT_DIR}/ink.designapi.codex.plist"
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
wire_api = "chat"
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

# LaunchAgent — пробрасываем env в GUI-приложения (включая Codex App)
mkdir -p "$LAUNCH_AGENT_DIR"
cat > "$LAUNCH_AGENT" <<EOF
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

launchctl unload "$LAUNCH_AGENT" 2>/dev/null || true
launchctl load   "$LAUNCH_AGENT"
launchctl setenv OPENAI_API_KEY  "$API_KEY"
launchctl setenv OPENAI_BASE_URL "$BASE_URL"

echo
echo "✅ Codex App configured for designapi.ink"
echo "   Полностью закрой и снова открой Codex App"
