#!/usr/bin/env bash
set -euo pipefail

API_KEY="${1:-${OPENAI_API_KEY:-__API_KEY__}}"
BASE_URL="__BASE_URL__/v1"
MODEL="__MODEL__"

CODEX_DIR="${HOME}/.codex"
AUTH_PATH="${CODEX_DIR}/auth.json"
CONFIG_PATH="${CODEX_DIR}/config.toml"
ENV_PATH="${CODEX_DIR}/designapi.env"
BACKUP_DIR="${CODEX_DIR}/backups/designapi-$(date +%Y%m%d-%H%M%S)"

mkdir -p "$CODEX_DIR"
if [[ -f "$CONFIG_PATH" || -f "$AUTH_PATH" ]]; then
  mkdir -p "$BACKUP_DIR"
  [[ -f "$CONFIG_PATH" ]] && cp "$CONFIG_PATH" "$BACKUP_DIR/" && echo "Backed up config.toml -> $BACKUP_DIR"
  [[ -f "$AUTH_PATH"   ]] && cp "$AUTH_PATH"   "$BACKUP_DIR/" && echo "Backed up auth.json   -> $BACKUP_DIR"
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

SHELL_NAME="$(basename "${SHELL:-bash}")"
case "$SHELL_NAME" in
  zsh)  RC="$HOME/.zshrc" ;;
  bash) RC="$HOME/.bashrc" ;;
  fish) RC="$HOME/.config/fish/config.fish" ;;
  *)    RC="" ;;
esac

if [[ -n "$RC" ]]; then
  HOOK="# designapi.ink
[ -f \"$ENV_PATH\" ] && . \"$ENV_PATH\""
  if [[ "$SHELL_NAME" == "fish" ]]; then
    HOOK="# designapi.ink
test -f $ENV_PATH; and source $ENV_PATH"
  fi
  if ! grep -q "designapi.ink" "$RC" 2>/dev/null; then
    mkdir -p "$(dirname "$RC")"
    printf '\n%s\n' "$HOOK" >> "$RC"
    echo "Added env hook to $RC"
  fi
fi

echo
echo "✅ Codex CLI configured for designapi.ink"
echo "   base_url: $BASE_URL"
echo "   open a new terminal or: source $ENV_PATH"
