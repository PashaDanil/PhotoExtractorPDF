#!/usr/bin/env pwsh

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

# Repo root = two levels up from tools/scripts
$repoRoot = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
$makefilePath = Join-Path $repoRoot 'Makefile'

Write-Host "Usage: make [target]"
Write-Host ""
Write-Host "Targets:"

Get-Content -LiteralPath $makefilePath | ForEach-Object {
    if ($_ -match '^(?<t>[A-Za-z0-9_-]+):.*?## (?<d>.*)$') {
        '{0,-20} {1}' -f $Matches.t, $Matches.d
    }
}
