# Forwards a Claude Code hook payload to pet-core. Never prints, never fails.
param([string]$Event)
try {
    $dir = Join-Path $env:APPDATA 'ClaudePet'
    $port = (Get-Content (Join-Path $dir 'port') -Raw).Trim()
    $token = (Get-Content (Join-Path $dir 'token') -Raw).Trim()
    $body = [Console]::In.ReadToEnd()
    Invoke-RestMethod -Method Post -Uri "http://127.0.0.1:$port/event?event=$Event" `
        -Headers @{ 'X-Pet-Token' = $token } -ContentType 'application/json' `
        -Body $body -TimeoutSec 1 | Out-Null
} catch {}
exit 0
