# Decision Debt Tracker (DDT)

**Author:** Nidhal Elkebir

## Core Philosophy

Technical debt is tracked — decision debt is invisible, yet more costly. DDT is a lightweight system to log, contextualize, and monitor key technical decisions until they become obsolete.

## What is Decision Debt?

Every technical decision carries assumptions. When those assumptions become invalid, the decision becomes "debt" — often invisible until it causes problems. DDT makes these decisions visible and alerts you when they need revisiting.

## Features

✅ **Decision Logging** - Quick 2-minute interactive CLI to capture decisions  
✅ **Context Tracking** - Record why, who, when, assumptions, and constraints  
✅ **Expiration Triggers** - Define when a decision should be reconsidered  
✅ **Smart Detection** - Automatic scanning for outdated decisions  
✅ **Version Control** - Plain-text YAML backups in `.ddt/` directory  
✅ **Fast Queries** - Search by keyword, view history, understand "why"  
✅ **Minimal Overhead** - CLI-first, locally run, no external dependencies  

## Installation

### From Source

```bash
git clone https://github.com/nidhalelbkir/ddt.git
cd ddt
go build -o ddt main.go
```

### Add to PATH (Windows)

```powershell
# Move binary to a permanent location
move ddt.exe C:\tools\ddt.exe

# Add to PATH
$env:Path += ";C:\tools"
# Make permanent:
[Environment]::SetEnvironmentVariable("Path", $env:Path, [System.EnvironmentVariableTarget]::User)
```

### Add to PATH (macOS/Linux)

```bash
sudo mv ddt /usr/local/bin/
chmod +x /usr/local/bin/ddt
```

## Quick Start

### 1. Initialize DDT in your project

```bash
cd your-project
ddt init
```

This creates a `.ddt/` directory with:
- `decisions.db` - SQLite database
- `backups/` - YAML exports of decisions

### 2. Log your first decision

```bash
ddt add "Use SQLite for local cache"
```

You'll be prompted for:
- **Why?** → Reasoning behind the decision
- **Who decided?** → Person or team name
- **Assumptions** → What must remain true (e.g., "Data size < 1GB")
- **Expiration triggers** → When to reconsider (e.g., "When data exceeds 1GB")
- **Linked resources** → Git commits, tickets, docs

### 3. View decisions

```bash
# List recent decisions
ddt list

# Show full details
ddt show ddt-2025-042

# Search decisions
ddt search "cache"
```

### 4. Check for outdated decisions

```bash
ddt check
```

The detector scans your codebase and checks if assumptions are still valid.

## Example Workflow

```bash
# Add a decision
$ ddt add "Use JSON for config files"
Why? → "Team familiarity, no need for schema yet"
Assumptions:
  1: Configs under 10 keys
  2: No nesting needed
Expiration trigger:
  1: When configs require nested sections
Linked commit → abc123
✅ Decision saved [ID: ddt-2025-042]

# Weeks later, check decisions
$ ddt check
⚠️  Decision ddt-2025-042 may be outdated
   Trigger: When configs require nested sections
   Evidence: Found keyword 'nested' in codebase
   
# Review the decision
$ ddt show ddt-2025-042

# Update or retire it
$ ddt update ddt-2025-042
$ ddt retire ddt-2025-042
```

## Commands

### Core Commands

| Command | Description |
|---------|-------------|
| `ddt init` | Initialize DDT in current directory |
| `ddt add "title"` | Log a new decision (interactive) |
| `ddt list` | List recent decisions |
| `ddt show [id]` | Show full decision details |
| `ddt search [query]` | Search decisions by keyword |
| `ddt check` | Run debt detection engine |
| `ddt update [id]` | Update an existing decision |
| `ddt retire [id]` | Mark a decision as obsolete |

### Data Management

| Command | Description |
|---------|-------------|
| `ddt export` | Export all decisions to YAML |
| `ddt backup` | Create timestamped backup |

### Options

```bash
# List only active decisions
ddt list --status active

# List all (including retired)
ddt list --all

# Limit results
ddt list --limit 5
```

## How Detection Works

The debt detection engine analyzes expiration triggers and checks your codebase:

### Trigger Types

1. **File Size** - "When data exceeds 1GB"
   - Monitors database, config, and data files

2. **File Existence** - "When we add multi-user access"
   - Searches for new files/features

3. **Count Thresholds** - "When we have more than 10 endpoints"
   - Counts occurrences in codebase

4. **Time-Based** - "After 6 months"
   - Checks decision age

5. **Keyword Presence** - "When we implement authentication"
   - Scans source code for keywords

### Example Triggers

```
"When data size exceeds 500MB"
"When we add more than 5 microservices"
"When configs require nested sections"
"After 1 year"
"When we implement user authentication"
"When team size exceeds 10 developers"
```

## Storage & Backup

### Local Storage

- **Database:** `.ddt/decisions.db` (SQLite)
- **Backups:** `.ddt/backups/*.yaml`

### YAML Format

Each decision is also exported as YAML for easy version control:

```yaml
id: ddt-2025-042
title: Use SQLite for local cache
reasoning: Simple, embedded, no server needed
decided_by: Engineering Team
decided_at: 2025-01-15T10:30:00Z
assumptions:
  - Data size < 1GB
  - No concurrent writes
expiration_triggers:
  - When data exceeds 1GB
  - When we add multi-user access
linked_resources:
  - abc123
  - https://ticket-tracker.com/ISSUE-42
tags:
  - database
  - caching
status: active
updated_at: 2025-01-15T10:30:00Z
```

### Version Control

Add `.ddt/` to your repository:

```bash
git add .ddt/
git commit -m "Add decision: Use SQLite for caching"
```

## Best Practices

### ✅ Do

- **Be specific** - Write clear, concrete triggers
- **Link everything** - Commits, tickets, docs
- **Review regularly** - Run `ddt check` in CI/CD
- **Keep it quick** - 2 minutes or less per decision
- **Track ownership** - Record who decided

### ❌ Don't

- **Over-document** - DDT is not a wiki
- **Skip triggers** - Always define expiration conditions
- **Ignore alerts** - Review flagged decisions promptly
- **Track everything** - Focus on impactful decisions

## CI/CD Integration

### GitHub Actions

```yaml
name: Check Decision Debt

on: [push, pull_request]

jobs:
  check-decisions:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      
      - name: Setup Go
        uses: actions/setup-go@v2
        with:
          go-version: '1.21'
      
      - name: Install DDT
        run: |
          git clone https://github.com/nidhalelbkir/ddt
          cd ddt && go build -o ddt main.go
          sudo mv ddt /usr/local/bin/
      
      - name: Check Decisions
        run: ddt check
```

### Pre-commit Hook

```bash
#!/bin/bash
# .git/hooks/pre-commit

echo "Checking decision debt..."
ddt check

if [ $? -ne 0 ]; then
  echo "⚠️  Some decisions may be outdated. Review with 'ddt check'"
  echo "Continue anyway? (y/n)"
  read -r response
  if [ "$response" != "y" ]; then
    exit 1
  fi
fi
```

## Architecture

```
ddt/
├── main.go              # Entry point
├── cmd/                 # CLI commands
│   ├── root.go
│   ├── add.go           # Add decisions
│   ├── list.go          # List decisions
│   ├── search.go        # Search decisions
│   ├── show.go          # Show details
│   ├── check.go         # Run detection
│   ├── update.go        # Update decisions
│   ├── retire.go        # Retire decisions
│   ├── export.go        # Export to YAML
│   ├── backup.go        # Create backups
│   └── init.go          # Initialize DDT
├── models/              # Data structures
│   └── decision.go      # Decision, Alert models
├── storage/             # Database layer
│   ├── storage.go       # SQLite operations
│   ├── export.go        # YAML export/import
│   └── utils.go         # Helpers
└── detector/            # Debt detection
    └── detector.go      # Trigger evaluation
```

## Tech Stack

- **Language:** Go 1.21+
- **Database:** SQLite (via mattn/go-sqlite3)
- **CLI:** Cobra framework
- **Prompts:** promptui
- **Format:** YAML (gopkg.in/yaml.v3)

## Roadmap

### Current Version (v0.1)

- ✅ Core CLI commands
- ✅ SQLite storage
- ✅ YAML export
- ✅ Basic debt detection
- ✅ Interactive prompts

### Planned Features

- [ ] Git integration (`ddt why <commit>`)
- [ ] VS Code extension
- [ ] Web dashboard
- [ ] Team sync (encrypted)
- [ ] Slack/email notifications
- [ ] Custom detection rules
- [ ] Decision templates
- [ ] Impact analysis

## Contributing

Contributions welcome! Areas of interest:

1. **Detection Rules** - More sophisticated trigger types
2. **Integrations** - Jira, GitHub, GitLab
3. **Visualization** - Decision timeline, dependency graph
4. **Testing** - Unit tests, integration tests

## License

MIT License - See LICENSE file for details

## Philosophy

> "Every decision is a bet on the future. DDT tracks those bets and tells you when it's time to cash out."

DDT is not about documenting everything. It's about tracking the decisions that matter — the ones that will bite you when assumptions change. Use it to:

- Remember *why* you made choices
- Know *when* to revisit them
- Prevent *invisible debt* from accumulating

## Support

- **Issues:** [GitHub Issues](https://github.com/nidhalelbkir/ddt/issues)
- **Discussions:** [GitHub Discussions](https://github.com/nidhalelbkir/ddt/discussions)
- **Email:** nidhal@example.com

---

Built with ❤️ by Nidhal Elkebir
