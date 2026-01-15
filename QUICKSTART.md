# Quick Start Guide - Decision Debt Tracker

## Setup Instructions

### Step 1: Install Go

DDT is written in Go, so you'll need Go installed first.

**Download Go:**
- Visit: https://go.dev/dl/
- Download Go 1.21 or higher for Windows
- Run the installer
- Verify installation: `go version`

### Step 2: Build DDT

```powershell
# Navigate to DDT directory
cd "c:\Users\Mega Pc\Desktop\DDD"

# Download dependencies
go mod tidy

# Build the executable
go build -o ddt.exe main.go
```

### Step 3: Test DDT

```powershell
# Show help
.\ddt.exe --help

# Initialize in a test directory
mkdir test-project
cd test-project
..\ddt.exe init
```

### Step 4: Use DDT

```powershell
# Add a decision
.\ddt.exe add "Use SQLite for local database"

# List decisions
.\ddt.exe list

# Show details
.\ddt.exe show [decision-id]

# Check for outdated decisions
.\ddt.exe check
```

## Optional: Add to PATH

To use `ddt` from anywhere:

```powershell
# Copy to a permanent location
mkdir C:\tools
copy ddt.exe C:\tools\ddt.exe

# Add to PATH (run as Administrator)
$oldPath = [Environment]::GetEnvironmentVariable("Path", [System.EnvironmentVariableTarget]::User)
$newPath = $oldPath + ";C:\tools"
[Environment]::SetEnvironmentVariable("Path", $newPath, [System.EnvironmentVariableTarget]::User)

# Restart PowerShell, then test
ddt --help
```

## Available Commands

| Command | Description |
|---------|-------------|
| `ddt init` | Initialize DDT in current directory |
| `ddt add "title"` | Add a new decision (interactive) |
| `ddt list` | List all decisions |
| `ddt show [id]` | Show decision details |
| `ddt search [query]` | Search decisions |
| `ddt check` | Check for outdated decisions |
| `ddt update [id]` | Update a decision |
| `ddt retire [id]` | Retire a decision |
| `ddt export` | Export decisions to YAML |
| `ddt backup` | Create a backup |

## Example Session

```powershell
PS> cd my-project
PS> ddt init
✅ DDT initialized in C:\my-project
   Database: C:\my-project\.ddt\decisions.db
   Backups: C:\my-project\.ddt\backups\

PS> ddt add "Use JSON for config files"

📝 Adding decision: Use JSON for config files

Why was this decision made?: Team familiar with JSON, simple structure
Who decided? (person/team): Backend Team

Assumptions (enter each assumption, empty line to finish):
  1: Config files under 10 keys
  2: No nesting required
  3: 

Expiration triggers (when should this be reconsidered?):
  1: When configs require nested sections
  2: When we need schema validation
  3: 

Linked resources (commit hash, URL, doc path - optional):
  1: abc123
  2: 

Tags (optional, comma-separated): configuration, json

✅ Decision saved [ID: ddt-2026-01-15-042]

PS> ddt list

📋 Recent Decisions (1)

🟢 [ddt-2026-01-15-042] Use JSON for config files
   Decided: 2026-01-15 by Backend Team
   Why: Team familiar with JSON, simple structure
   Assumptions: 2 | Triggers: 2
   Tags: configuration, json

PS> ddt check

🔍 Checking decisions for expiration triggers...
✅ All decisions are up to date! ✨
```

## Troubleshooting

### "go: command not found"
- Go is not installed or not in PATH
- Download from: https://go.dev/dl/
- Restart PowerShell after installation

### "gcc: command not found"
- SQLite requires CGO (C compiler)
- Install MinGW-w64: https://sourceforge.net/projects/mingw-w64/
- Or use TDM-GCC: https://jmeubank.github.io/tdm-gcc/

### "module not found"
```powershell
go mod download
go mod tidy
```

### Permission errors
- Run PowerShell as Administrator
- Check antivirus isn't blocking

## Next Steps

1. **Read the full README** - [README.md](README.md)
2. **See examples** - [EXAMPLES.md](EXAMPLES.md)
3. **Integrate with Git** - Add `.ddt/` to version control
4. **Set up CI/CD** - Run `ddt check` in your pipeline
5. **Customize** - Modify detection rules in `detector/detector.go`

## Support

- **Issues**: Create an issue on GitHub
- **Questions**: Check EXAMPLES.md
- **Feature requests**: Open a discussion

---

**Author:** Nidhal Elkebir  
**License:** MIT  
**Version:** 0.1.0
