# DDT - Installation & First Steps

## What You Have

A complete **Decision Debt Tracker** system with:
- ✅ Full-featured CLI tool
- ✅ SQLite database storage
- ✅ YAML backup system
- ✅ Automatic debt detection
- ✅ Interactive prompts
- ✅ Search and filtering

## Quick Setup (3 Steps)

### Step 1: Install Go (if not already installed)

**Download:** https://go.dev/dl/

- Get Go 1.21 or newer for Windows
- Run the installer
- Restart PowerShell

**Verify:**
```powershell
go version
```

### Step 2: Build DDT

```powershell
# Go to project directory
cd "c:\Users\Mega Pc\Desktop\DDD"

# Download dependencies
go mod tidy

# Build executable
go build -o ddt.exe main.go
```

**You should now have `ddt.exe` in the directory!**

### Step 3: Try It Out

```powershell
# Show help
.\ddt.exe --help

# Create a test directory
mkdir test-ddt
cd test-ddt

# Initialize DDT
..\ddt.exe init

# Add your first decision
..\ddt.exe add "Use SQLite for local storage"
```

## What to Read First

1. **[QUICKSTART.md](QUICKSTART.md)** - Start here! Complete getting started guide
2. **[README.md](README.md)** - Full documentation
3. **[EXAMPLES.md](EXAMPLES.md)** - Real-world usage examples
4. **[BUILD.md](BUILD.md)** - Detailed build instructions

## Commands Reference

```powershell
ddt init                        # Initialize in project
ddt add "decision title"        # Add new decision (interactive)
ddt list                        # List all decisions
ddt list --status active        # List only active
ddt search "keyword"            # Search decisions
ddt show [id]                   # Show full details
ddt check                       # Check for outdated decisions
ddt update [id]                 # Update a decision
ddt retire [id]                 # Retire a decision
ddt export                      # Export to YAML
ddt backup                      # Create backup
```

## Project Structure at a Glance

```
DDD/
├── main.go              # Entry point
├── cmd/                 # All CLI commands
│   ├── add.go          # Add decisions
│   ├── list.go         # List decisions
│   ├── search.go       # Search decisions
│   ├── check.go        # Debt detection
│   └── ...
├── models/             # Data structures
├── storage/            # Database & YAML
├── detector/           # Detection engine
└── README.md           # Full docs
```

## Common Issues & Solutions

### "go: command not found"
👉 Install Go from https://go.dev/dl/ and restart PowerShell

### "gcc: command not found"
👉 SQLite needs a C compiler. Install:
- MinGW-w64: https://sourceforge.net/projects/mingw-w64/
- OR TDM-GCC: https://jmeubank.github.io/tdm-gcc/

### Build errors
```powershell
go mod download
go mod tidy
go clean -cache
```

## Next Steps

### 1. Make It Permanent
```powershell
# Copy to tools directory
mkdir C:\tools
copy ddt.exe C:\tools\

# Add to PATH (restart PowerShell after)
[Environment]::SetEnvironmentVariable(
    "Path", 
    [Environment]::GetEnvironmentVariable("Path", "User") + ";C:\tools", 
    "User"
)

# Now you can use from anywhere
ddt --help
```

### 2. Use in Real Project
```powershell
cd your-actual-project
ddt init
ddt add "Your first real decision"
```

### 3. Integrate with Git
```powershell
# The .ddt directory is version-controllable
git add .ddt/
git commit -m "Add decision tracking"
```

### 4. Set Up Automation
- Run `ddt check` in your CI/CD pipeline
- Add pre-commit hooks
- Schedule weekly reviews

## Example Workflow

```powershell
# Initialize in project
PS> cd C:\MyProject
PS> ddt init
✅ DDT initialized

# Add a decision
PS> ddt add "Use Redis for caching"

Why? → "Need fast key-value store"
Who? → "Backend Team"
Assumptions:
  1: Cache size < 500MB
  2: No persistence required
  3: [empty line]
Triggers:
  1: When cache exceeds 1GB
  2: When we need persistence
  3: [empty line]

✅ Decision saved [ID: ddt-2026-01-15-123]

# Check decisions
PS> ddt check
🔍 Checking decisions...
✅ All decisions are up to date!

# List decisions
PS> ddt list
🟢 [ddt-2026-01-15-123] Use Redis for caching
   Decided: 2026-01-15 by Backend Team
```

## File Locations

After running `ddt init`, you'll have:

```
your-project/
├── .ddt/
│   ├── decisions.db              # SQLite database
│   └── backups/
│       ├── ddt-2026-*.yaml       # Individual backups
│       └── backup-*.yaml         # Full backups
```

## Documentation Guide

| Document | What It's For |
|----------|---------------|
| **QUICKSTART.md** | 🚀 Getting started (read first!) |
| **README.md** | 📖 Complete guide and philosophy |
| **EXAMPLES.md** | 💡 Real-world usage examples |
| **BUILD.md** | 🔧 Building and installation |
| **PROJECT.md** | 🗺️ Project structure and internals |
| **INSTALL.md** | 📦 This file - quick setup |

## Core Features

✨ **What DDT Does:**
- Logs technical decisions with full context
- Tracks assumptions and constraints
- Defines when decisions should be reconsidered
- Automatically detects when assumptions break
- Alerts you before decisions become technical debt
- Exports everything to version-controllable YAML

🎯 **Use Cases:**
- "Why did we choose this database?"
- "When should we reconsider this architecture?"
- "What were the assumptions behind this choice?"
- "Which decisions are now outdated?"

## Philosophy

> **Technical debt is tracked — decision debt is invisible, yet more costly.**

DDT makes the invisible visible. It's not about documenting everything, but tracking decisions that will matter later.

Every technical choice is a bet on the future. DDT helps you remember:
- **What** you decided
- **Why** you decided it
- **When** to revisit it

## Need Help?

1. **Check EXAMPLES.md** - See real-world usage
2. **Read README.md** - Full documentation
3. **Review BUILD.md** - Build troubleshooting
4. **Check PROJECT.md** - Internal architecture

## Success Indicators

You're using DDT successfully when:
- ✅ Adding a decision takes < 2 minutes
- ✅ You can instantly answer "why did we do this?"
- ✅ Outdated decisions are caught before causing issues
- ✅ New team members can understand past choices
- ✅ Technical discussions have context

## Ready to Start?

```powershell
# 1. Build
go mod tidy
go build -o ddt.exe main.go

# 2. Test
.\ddt.exe --help

# 3. Use
cd your-project
.\ddt.exe init
.\ddt.exe add "Your first decision"
```

---

**Author:** Nidhal Elkebir  
**License:** MIT  
**Version:** 0.1.0

🎉 **Happy Decision Tracking!**
