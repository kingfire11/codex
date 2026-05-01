#!/usr/bin/env bash
set -euo pipefail

API_KEY="${1:-${OPENAI_API_KEY:-__API_KEY__}}"
BASE_URL="__BASE_URL__/v1"
MODEL="__MODEL__"

OC_DIR="${HOME}/.config/opencode"
OC_CFG="${OC_DIR}/opencode.json"
BACKUP_DIR="${OC_DIR}/backups/designapi-$(date +%Y%m%d-%H%M%S)"

mkdir -p "$OC_DIR"
if [[ -f "$OC_CFG" ]]; then
  mkdir -p "$BACKUP_DIR"
  cp "$OC_CFG" "$BACKUP_DIR/"
  echo "Backed up opencode.json -> $BACKUP_DIR"
fi

cat > "$OC_CFG" <<EOF
{
  "\$schema": "https://opencode.ai/config.json",
  "provider": {
    "designapi": {
      "npm": "@ai-sdk/openai-compatible",
      "name": "DesignAPI",
      "options": {
        "baseURL": "$BASE_URL",
        "apiKey": "$API_KEY"
      },
      "models": {
        "$MODEL": { "name": "$MODEL" }
      }
    }
  }
}
EOF
chmod 600 "$OC_CFG"

echo
echo "✅ OpenCode configured for designapi.ink"
echo "   Перезапусти OpenCode"
