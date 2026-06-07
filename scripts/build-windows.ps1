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

$version = $env:GOFINDER_VERSION
if ([string]::IsNullOrWhiteSpace($version)) {
	$version = (& git describe --tags --always --dirty 2>$null)
	if ([string]::IsNullOrWhiteSpace($version)) {
		$version = "dev"
	}
}

$versionPackage = "github.com/adelylria/GoFinder/core/version.Version"
$commit = (& git rev-parse --short HEAD 2>$null)
if ([string]::IsNullOrWhiteSpace($commit)) {
	$commit = "unknown"
}
$buildDate = (Get-Date).ToUniversalTime().ToString("yyyy-MM-ddTHH:mm:ssZ")
$commitPackage = "github.com/adelylria/GoFinder/core/version.Commit"
$buildDatePackage = "github.com/adelylria/GoFinder/core/version.BuildDate"

& go build -ldflags="-H=windowsgui -s -w -X '$versionPackage=$version' -X '$commitPackage=$commit' -X '$buildDatePackage=$buildDate'" -o "build/goFinder.exe" "./cmd"
exit $LASTEXITCODE
