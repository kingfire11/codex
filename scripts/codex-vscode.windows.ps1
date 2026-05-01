$ErrorActionPreference = "Stop"

$ApiKey  = if ($args.Count -ge 1) { $args[0] } elseif ($env:OPENAI_API_KEY) { $env:OPENAI_API_KEY } else { "__API_KEY__" }
$BaseUrl = "__BASE_URL__/v1"
$Model   = "__MODEL__"

$CodexDir   = Join-Path $HOME ".codex"
$AuthPath   = Join-Path $CodexDir "auth.json"
$ConfigPath = Join-Path $CodexDir "config.toml"
$EnvPath    = Join-Path $CodexDir "designapi.env"
$BackupDir  = Join-Path $CodexDir ("backups\designapi-" + (Get-Date -Format "yyyyMMdd-HHmmss"))

New-Item -ItemType Directory -Force -Path $CodexDir | Out-Null
if ((Test-Path $ConfigPath) -or (Test-Path $AuthPath)) {
  New-Item -ItemType Directory -Force -Path $BackupDir | Out-Null
  if (Test-Path $ConfigPath) { Copy-Item $ConfigPath $BackupDir }
  if (Test-Path $AuthPath)   { Copy-Item $AuthPath   $BackupDir }
  Write-Host "Backups saved -> $BackupDir"
}

@"
model = "$Model"
model_provider = "designapi"

[model_providers.designapi]
name = "DesignAPI"
base_url = "$BaseUrl"
wire_api = "chat"
env_key = "OPENAI_API_KEY"
"@ | Set-Content -Path $ConfigPath -Encoding UTF8

@"
{"OPENAI_API_KEY":"$ApiKey"}
"@ | Set-Content -Path $AuthPath -Encoding UTF8

@"
`$env:OPENAI_API_KEY = "$ApiKey"
`$env:OPENAI_BASE_URL = "$BaseUrl"
"@ | Set-Content -Path $EnvPath -Encoding UTF8

[Environment]::SetEnvironmentVariable("OPENAI_API_KEY",  $ApiKey,  "User")
[Environment]::SetEnvironmentVariable("OPENAI_BASE_URL", $BaseUrl, "User")

Write-Host ""
Write-Host "✅ VS Code Codex extension configured for designapi.ink"
Write-Host "   1) Полностью перезапусти VS Code"
Write-Host "   2) Открой панель Codex — должно быть «Logged in with API key»"
