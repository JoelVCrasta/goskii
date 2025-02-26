$AppName = "goskii"
$Repo = "JoelVCrasta/goskii"
$BinDir = "$env:ProgramFiles\$AppName"

# Check permissions
If (-Not ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] "Administrator")) {
    Write-Host "This script must be run as an administrator."
    exit 1
}

# Determine OS
$OS = $env:OS
if (-not $OS -or $OS -notlike "*Windows*") {
    Write-Host "Unsupported OS: $OS"
    exit 1
}

# Determine Architecture
$Arch = [System.Environment]::Is64BitOperatingSystem
if (-not $Arch) {
    Write-Host "Unsupported Architecture: $Arch"
    exit 1
}

# Fetch Latest Version
Write-Host "Fetching latest version of $AppName..."
$LatestVersion = (Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest").tag_name

if (-not $LatestVersion) {
    Write-Host "Failed to fetch latest version of $AppName"
    exit 1
}

# Download dependencies
Write-Host "Downloading dependencies..."

# Check if FFmpeg is installed
if (-not (Get-Command ffmpeg -ErrorAction SilentlyContinue)) {
    Write-Host "FFmpeg not found, downloading latest FFmpeg version..."
    $FFmpegUrl = "https://www.gyan.dev/ffmpeg/builds/ffmpeg-release-essentials.zip"
    $FFmpegZip = "$env:TEMP\ffmpeg.zip"

    try {
        Invoke-WebRequest -Uri $FFmpegUrl -OutFile $FFmpegZip
        Expand-Archive -Path $FFmpegZip -DestinationPath "$env:TEMP\ffmpeg" -Force
        $FFmpegBin = Get-ChildItem -Path "$env:TEMP\ffmpeg" -Recurse -Filter "ffmpeg.exe" | Select-Object -First 1

        if ($FFmpegBin) {
            $FFmpegSystemPath = "$env:SystemRoot\System32\ffmpeg.exe"
            Copy-Item $FFmpegBin.FullName -Destination $FFmpegSystemPath -Force
            Write-Host "FFmpeg installed successfully."
        } else {
            Write-Host "Failed to locate ffmpeg.exe in the downloaded archive."
            exit 1
        }
    } catch {
        Write-Host "Failed to download or extract FFmpeg: $_"
        exit 1
    } finally {
        Remove-Item -Path "$env:TEMP\ffmpeg" -Recurse -Force -ErrorAction SilentlyContinue
        Remove-Item -Path $FFmpegZip -Force -ErrorAction SilentlyContinue
    }
} else {
    Write-Host "FFmpeg already installed."
}

# Check if yt-dlp is installed
if (-not (Get-Command yt-dlp -ErrorAction SilentlyContinue)) {
    Write-Host "yt-dlp not found, downloading latest yt-dlp version..."
    
    try {
        $YtDlpUrl = (Invoke-RestMethod -Uri "https://api.github.com/repos/yt-dlp/yt-dlp/releases/latest").assets | Where-Object { $_.name -eq "yt-dlp.exe" } | Select-Object -ExpandProperty browser_download_url
    } catch {
        Write-Host "Failed to fetch the latest version of yt-dlp."
        exit 1
    }

    $YtDlpSystemPath = "$env:SystemRoot\System32\yt-dlp.exe"
    Invoke-WebRequest -Uri $YtDlpUrl -OutFile $YtDlpSystemPath
    Write-Host "yt-dlp installed successfully."
} else {
    Write-Host "yt-dlp already installed."
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

    if ($env:Path -notlike "*$BinDir*") {
        $env:Path += ";$BinDir"
    }

    Write-Host "Restart your terminal or system to apply the PATH changes."
    Write-Host "Run `$AppName --help` or `$AppName -h` to get started."
} else {
    Write-Host "Failed to install $AppName."
    exit 1
}
