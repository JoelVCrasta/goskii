$AppName = "goskii"
$Repo = "JoelVCrasta/goskii"
$BinDir = "$env:ProgramFiles\$AppName"


If (-Not ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] "Administrator")) {
    Write-Host "This script must be run as an administrator."
    exit 1
}


$OS = $env:OS
$Arch = [System.Environment]::Is64BitOperatingSystem

if (-not $OS -or $OS -notlike "*Windows*") {
    Write-Host "Unsupported OS: $OS"
    exit 1
}

if (-not $Arch) {
    Write-Host "Unsupported Architecture: $Arch"
    exit 1
}


Write-Host "Fetching latest version of $AppName..."
$LatestVersion = (Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest").tag_name

if (-not $LatestVersion) {
    Write-Host "Failed to fetch latest version of $AppName"
    exit 1
}


$Binary = "$AppName-windows-amd64"
$DownloadUrl = "https://github.com/$Repo/releases/download/$LatestVersion/$Binary"

if (-not (Test-Path -Path $BinDir)) {
    Write-Host "Creating directory $BinDir..."
    New-Item -ItemType Directory -Path $BinDir | Out-Null
}

Write-Host "Downloading $Binary..."
$BinaryPath = "$BinDir\$AppName.exe"

try {
    Invoke-WebRequest -Uri $DownloadUrl -OutFile $BinaryPath
} catch {
    Write-Host "Failed to download $AppName from $DownloadUrl."
    exit 1
}


if (-not ($env:Path -like "*$BinDir*")) {
    Write-Host "Adding $BinDir to PATH..."
    [System.Environment]::SetEnvironmentVariable("Path", "$env:Path;$BinDir", [EnvironmentVariableTarget]::Machine)
}

if (Test-Path -Path $BinaryPath) {
    Write-Host "Successfully installed $AppName $LatestVersion to $BinDir."
    Write-Host "Restart your terminal or system to apply the PATH changes."
    Write-Host "Run `$AppName --help` or `$AppName -h` to get started."
} else {
    Write-Host "Failed to install $AppName."
    exit 1
}
