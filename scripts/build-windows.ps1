$ErrorActionPreference = "Stop"

$repoRoot = Split-Path -Parent $PSScriptRoot
Set-Location $repoRoot

New-Item -ItemType Directory -Force -Path "build" | Out-Null

& windres -O coff -F pe-x86-64 -o "cmd/gofinder_windows.syso" "cmd/gofinder.rc"
if ($LASTEXITCODE -ne 0) {
	exit $LASTEXITCODE
}

$env:GOOS = "windows"
$env:GOARCH = "amd64"

& go build -ldflags="-H=windowsgui -s -w" -o "build/goFinder.exe" "./cmd"
exit $LASTEXITCODE
