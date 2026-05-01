$ErrorActionPreference = "Stop"

$ApiKey  = if ($args.Count -ge 1) { $args[0] } elseif ($env:OPENAI_API_KEY) { $env:OPENAI_API_KEY } else { "__API_KEY__" }
$BaseUrl = "__BASE_URL__/v1"
$Model   = "__MODEL__"

$OcDir    = Join-Path $HOME ".config\opencode"
$OcCfg    = Join-Path $OcDir "opencode.json"
$BackupDir = Join-Path $OcDir ("backups\designapi-" + (Get-Date -Format "yyyyMMdd-HHmmss"))

New-Item -ItemType Directory -Force -Path $OcDir | Out-Null
if (Test-Path $OcCfg) {
  New-Item -ItemType Directory -Force -Path $BackupDir | Out-Null
  Copy-Item $OcCfg $BackupDir
  Write-Host "Backed up opencode.json -> $BackupDir"
}

@"
{
  "`$schema": "https://opencode.ai/config.json",
  "provider": {
    "designapi": {
      "npm": "@ai-sdk/openai-compatible",
      "name": "DesignAPI",
      "options": {
        "baseURL": "$BaseUrl",
        "apiKey": "$ApiKey"
      },
      "models": {
        "$Model": { "name": "$Model" }
      }
    }
  }
}
"@ | Set-Content -Path $OcCfg -Encoding UTF8

Write-Host ""
Write-Host "✅ OpenCode configured for designapi.ink"
Write-Host "   Перезапусти OpenCode"
