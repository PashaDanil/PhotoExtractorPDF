#!/usr/bin/env pwsh

# Script to generate OpenAPI documentation and update frontend API client

$ErrorActionPreference = "Stop"

# Get the root directory (two levels up from this script)
$rootDir = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
$backendDir = Join-Path $rootDir "apps\backend"
$contractsDir = Join-Path $rootDir "libs\contracts\openapi"
$frontendDir = Join-Path $rootDir "apps\frontend"

Write-Host "Generating Swagger documentation..." -ForegroundColor Cyan

# Navigate to backend directory and run swag init
Push-Location $backendDir
try {
    swag init -g main.go -d ./cmd/go-api,./internal/http/handler -o ./docs --outputTypes yaml,json
    if ($LASTEXITCODE -ne 0) {
        throw "Swag init failed with exit code $LASTEXITCODE"
    }
    Write-Host "✓ Swagger documentation generated" -ForegroundColor Green
} finally {
    Pop-Location
}

# Copy the generated YAML to contracts directory
Write-Host "Copying OpenAPI spec to contracts..." -ForegroundColor Cyan
$sourceYaml = Join-Path $backendDir "docs\swagger.yaml"
$targetYaml = Join-Path $contractsDir "backend.yaml"

if (Test-Path $sourceYaml) {
    Copy-Item $sourceYaml $targetYaml -Force
    Write-Host "✓ OpenAPI spec copied to $targetYaml" -ForegroundColor Green
} else {
    throw "Source YAML file not found: $sourceYaml"
}

# Generate frontend API client
Write-Host "Generating frontend API client..." -ForegroundColor Cyan
Push-Location $frontendDir
try {
    npm run api:gen
    if ($LASTEXITCODE -ne 0) {
        throw "Frontend API generation failed with exit code $LASTEXITCODE"
    }
    Write-Host "✓ Frontend API client generated" -ForegroundColor Green
} finally {
    Pop-Location
}

Write-Host "`n✓ All done!" -ForegroundColor Green
