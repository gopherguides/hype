# Hype installer for Windows
# Usage: irm https://raw.githubusercontent.com/gopherguides/hype/main/install.ps1 | iex
# Or with version: .\install.ps1 -Version v0.5.0

param(
    [string]$Version = "latest",
    [string]$InstallDir = "$env:LOCALAPPDATA\hype\bin"
)

$ErrorActionPreference = "Stop"

$Repo = "gopherguides/hype"

function Get-Architecture {
    $arch = [System.Runtime.InteropServices.RuntimeInformation]::OSArchitecture
    switch ($arch) {
        "X64" { return "x86_64" }
        "Arm64" { return "arm64" }
        "X86" { return "i386" }
        default {
            Write-Error "Unsupported architecture: $arch"
            exit 1
        }
    }
}

function Main {
    $arch = Get-Architecture

    if ($Version -eq "latest") {
        Write-Host "Fetching latest version..."
        $release = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest"
        $Version = $release.tag_name
    }

    Write-Host "Installing hype $Version for Windows/$arch..."

    $archiveName = "hype_Windows_${arch}.zip"
    $downloadUrl = "https://github.com/$Repo/releases/download/$Version/$archiveName"

    $tempDir = Join-Path $env:TEMP "hype-install"
    if (Test-Path $tempDir) {
        Remove-Item -Recurse -Force $tempDir
    }
    New-Item -ItemType Directory -Path $tempDir | Out-Null

    $archivePath = Join-Path $tempDir $archiveName

    Write-Host "Downloading $downloadUrl..."
    Invoke-WebRequest -Uri $downloadUrl -OutFile $archivePath

    Write-Host "Extracting..."
    Expand-Archive -Path $archivePath -DestinationPath $tempDir -Force

    if (-not (Test-Path $InstallDir)) {
        New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
    }

    $exePath = Join-Path $InstallDir "hype.exe"
    Move-Item -Path (Join-Path $tempDir "hype.exe") -Destination $exePath -Force

    Remove-Item -Recurse -Force $tempDir

    $currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if ($currentPath -notlike "*$InstallDir*") {
        Write-Host "Adding $InstallDir to PATH..."
        [Environment]::SetEnvironmentVariable("Path", "$currentPath;$InstallDir", "User")
        $env:Path = "$env:Path;$InstallDir"
    }

    Write-Host ""
    Write-Host "hype installed successfully to $exePath"
    Write-Host ""
    Write-Host "You may need to restart your terminal for PATH changes to take effect."
    Write-Host ""

    & $exePath version
}

Main
