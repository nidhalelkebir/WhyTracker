# Building and Running DDT

## Prerequisites

- Go 1.21 or higher
- Git (optional, for version control)

## Build Instructions

### Windows

```powershell
# Navigate to project directory
cd "c:\Users\Mega Pc\Desktop\DDD"

# Download dependencies
go mod tidy

# Build the binary
go build -o ddt.exe main.go

# Test the binary
.\ddt.exe --help
```

### macOS/Linux

```bash
# Navigate to project directory
cd ~/ddt

# Download dependencies
go mod tidy

# Build the binary
go build -o ddt main.go

# Make executable
chmod +x ddt

# Test the binary
./ddt --help
```

## Installation

### Windows - Add to PATH

```powershell
# Create tools directory
mkdir C:\tools

# Copy binary
copy ddt.exe C:\tools\

# Add to PATH (current session)
$env:Path += ";C:\tools"

# Add to PATH (permanent)
[Environment]::SetEnvironmentVariable("Path", $env:Path, [System.EnvironmentVariableTarget]::User)

# Verify installation
ddt --help
```

### macOS/Linux - Add to PATH

```bash
# Copy to system path
sudo cp ddt /usr/local/bin/

# Verify installation
ddt --help
```

## First Run

```bash
# Initialize DDT in a project
cd your-project-directory
ddt init

# Add your first decision
ddt add "Your first technical decision"

# List decisions
ddt list

# Get help
ddt --help
```

## Development

### Run without building

```bash
go run main.go [command]
```

### Run tests (when added)

```bash
go test ./...
```

### Format code

```bash
go fmt ./...
```

## Troubleshooting

### "module not found" errors

```bash
go mod tidy
go mod download
```

### CGO errors (sqlite3)

**Windows:**
- Install MinGW-w64 or TDM-GCC
- Ensure gcc is in PATH

**macOS:**
```bash
xcode-select --install
```

**Linux:**
```bash
sudo apt-get install build-essential
# or
sudo yum install gcc
```

### Permission errors

**Windows:** Run PowerShell as Administrator

**macOS/Linux:**
```bash
chmod +x ddt
```

## Cross-Compilation

### Build for Linux from Windows

```powershell
$env:GOOS="linux"
$env:GOARCH="amd64"
$env:CGO_ENABLED="1"
go build -o ddt main.go
```

### Build for Windows from Linux

```bash
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -o ddt.exe main.go
```

### Build for macOS from Linux

```bash
GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -o ddt main.go
```

## Docker (Optional)

```dockerfile
FROM golang:1.21-alpine

RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o ddt main.go

ENTRYPOINT ["./ddt"]
```

Build and run:

```bash
docker build -t ddt .
docker run -v $(pwd):/workspace -w /workspace ddt list
```
