#!/usr/bin/env pwsh

# Script to generate OpenAPI documentation and update frontend API client

$ErrorActionPreference = "Stop"

# Get the root directory (two levels up from this script)
$rootDir = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
$jobsApiDir = Join-Path $rootDir "apps\jobs-api"
$contractsDir = Join-Path $rootDir "libs\contracts\openapi"
$frontendDir = Join-Path $rootDir "apps\frontend"
$swagEntryPoint = "cmd/jobs-api/main.go"

Write-Host "Generating Swagger documentation..." -ForegroundColor Cyan

# Navigate to jobs-api directory and run swag init
Push-Location $jobsApiDir
try {
    $output = swag init -g $swagEntryPoint --parseInternal --outputTypes go,json,yaml -o docs 2>&1
    $filteredOutput = $output | Where-Object { 
        $line = $_.ToString()
        -not ($line -match "warning: failed to get package name in dir: \./") -and
        -not ($line -match "warning: failed to evaluate const")
    }
    $filteredOutput | ForEach-Object { Write-Host $_ }
    
    if ($LASTEXITCODE -ne 0) {
        throw "Swag init failed with exit code $LASTEXITCODE"
    }
    Write-Host "✓ Swagger documentation generated" -ForegroundColor Green
} finally {
    Pop-Location
}

# Copy the generated YAML to contracts directory
Write-Host "Copying OpenAPI spec to contracts..." -ForegroundColor Cyan
$sourceYaml = Join-Path $jobsApiDir "docs\swagger.yaml"
$targetYaml = Join-Path $contractsDir "imgpdf.yaml"

if (Test-Path $sourceYaml) {
    Copy-Item $sourceYaml $targetYaml -Force
    Write-Host "✓ OpenAPI spec copied to $targetYaml" -ForegroundColor Green
} else {
    throw "Source YAML file not found: $sourceYaml"
}