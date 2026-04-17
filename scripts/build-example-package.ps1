# Build example/package.tgz from example/payload/ (repo root).
$ErrorActionPreference = "Stop"
$root = Split-Path -Parent $PSScriptRoot
if (-not (Test-Path (Join-Path $root "example\payload"))) {
    throw "example\payload not found. Run from repo root context."
}
Set-Location $root
$out = Join-Path $root "example\package.tgz"
$payload = Join-Path $root "example\payload"
if (Test-Path $out) { Remove-Item $out -Force }
# Windows 10+ ships tar.exe
& tar.exe -czf $out -C $payload .
Write-Host "Wrote $out"
