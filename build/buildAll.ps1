param (
    [string]$Version
)

# Validate Version parameter
if (-not $Version) {
    Write-Host "Usage: ./build.ps1 -Version <version>"
    exit 1
}

# Create dist directory if it doesn't exist
$distDir = "dist"
if (-not (Test-Path $distDir)) {
    New-Item -ItemType Directory -Path $distDir | Out-Null
}

# Function to build for specific architecture
function BuildForArchitecture {
    param (
        [string]$arch
    )

    Write-Host "Building for architecture: $arch"

    $outputName = "uptime_${Version}_windows_${arch}.exe"

    # Set environment variables inline and build for Windows with specific architecture
    $env:GOOS = "windows"
    $env:GOARCH = $arch
    & go build -o "$distDir\$outputName" ./cmd/uptime

    if ($LASTEXITCODE -ne 0) {
        Write-Host "Failed to build for architecture: $arch" -ForegroundColor Red
        exit 1
    }

    Write-Host "Build successful for architecture: $arch"
}

# Build for amd64 (64-bit)
BuildForArchitecture "amd64"

# Build for arm (ARM)
BuildForArchitecture "arm"

# Build for arm64 (ARM64)
BuildForArchitecture "arm64"

# Build for 386 (32-bit)
BuildForArchitecture "386"

Write-Host "All builds completed successfully."
