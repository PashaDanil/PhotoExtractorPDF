param(
  [string]$EnvFile = ".env.tools"
)

if (!(Test-Path $EnvFile)) {
  throw "Env file not found: $EnvFile"
}

Get-Content $EnvFile | ForEach-Object {
  $line = $_.Trim()
  if ($line -eq "" -or $line.StartsWith("#")) { return }

  $parts = $line.Split("=", 2)
  if ($parts.Count -ne 2) { return }

  $name  = $parts[0].Trim()
  $value = $parts[1].Trim().Trim('"')

  # экспортируем в текущую сессию
  Set-Item -Path "Env:$name" -Value $value

  # записываем в постоянную конфигурацию go env
  if ($name -in @("GOPRIVATE", "GONOSUMDB", "GONOPROXY", "GOPROXY")) {
    go env -w "$name=$value" | Out-Null
  }
}

Write-Host "Applied:"
Write-Host ("GOPRIVATE = " + (go env GOPRIVATE))
Write-Host ("GONOSUMDB = " + (go env GONOSUMDB))
