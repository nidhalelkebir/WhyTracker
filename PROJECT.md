# Decision Debt Tracker (DDT) - Project Summary

## 🎯 Project Overview

**Author:** Nidhal Elkebir  
**Version:** 0.1.0  
**License:** MIT  
**Language:** Go 1.21+

DDT is a CLI-first tool for tracking technical decisions, their assumptions, and automatically detecting when those assumptions become outdated.

## 📁 Project Structure

```
DDD/
├── main.go                    # Application entry point
├── go.mod                     # Go module dependencies
├── README.md                  # Full documentation
├── QUICKSTART.md              # Getting started guide
├── BUILD.md                   # Build instructions
├── EXAMPLES.md                # Usage examples
├── LICENSE                    # MIT license
├── .gitignore                 # Git ignore rules
│
├── cmd/                       # CLI commands
│   ├── root.go               # Root command & utilities
│   ├── init.go               # Initialize DDT
│   ├── add.go                # Add new decision
│   ├── list.go               # List decisions
│   ├── show.go               # Show decision details
│   ├── search.go             # Search decisions
│   ├── check.go              # Run debt detection
│   ├── update.go             # Update decision
│   ├── retire.go             # Retire decision
│   ├── export.go             # Export to YAML
│   └── backup.go             # Create backup
│
├── models/                    # Data structures
│   └── decision.go           # Decision, Alert, Update models
│
├── storage/                   # Persistence layer
│   ├── storage.go            # SQLite database operations
│   ├── export.go             # YAML import/export
│   └── utils.go              # Helper functions
│
└── detector/                  # Debt detection engine
    └── detector.go           # Trigger evaluation logic
```

## ✨ Core Features Implemented

### 1. Decision Management
- ✅ Interactive CLI for adding decisions
- ✅ Full CRUD operations (Create, Read, Update, Retire)
- ✅ Rich metadata capture (why, who, when, assumptions, triggers)
- ✅ Tagging and categorization
- ✅ Status tracking (active, updated, retired)

### 2. Query & Search
- ✅ List decisions with filters (status, limit)
- ✅ Full-text search across title, reasoning, tags
- ✅ Detailed decision view with full context
- ✅ Show related alerts

### 3. Debt Detection
- ✅ Automatic scanning of expiration triggers
- ✅ Multiple trigger types:
  - File size thresholds
  - File existence checks
  - Count/numeric thresholds
  - Time-based expiration
  - Keyword presence detection
- ✅ Alert generation with evidence
- ✅ Severity classification (warning, critical)

### 4. Data Management
- ✅ SQLite database for fast queries
- ✅ YAML export for version control
- ✅ Automatic backups
- ✅ Timestamped backup files
- ✅ Plain-text format for human readability

### 5. User Experience
- ✅ Minimal overhead (< 2 minutes to log decision)
- ✅ Interactive prompts with validation
- ✅ Clear visual feedback (emojis, formatting)
- ✅ Helpful error messages
- ✅ Comprehensive help text

## 🛠️ Technical Stack

| Component | Technology | Purpose |
|-----------|------------|---------|
| Language | Go 1.21+ | Performance, CLI tooling |
| CLI Framework | Cobra | Command structure |
| Database | SQLite | Local storage |
| Prompts | promptui | Interactive input |
| Serialization | YAML | Human-readable backups |
| Storage | mattn/go-sqlite3 | SQLite driver |

## 📊 Database Schema

### decisions
- id (TEXT, PRIMARY KEY)
- title (TEXT)
- reasoning (TEXT)
- decided_by (TEXT)
- decided_at (DATETIME)
- assumptions (TEXT, JSON array)
- expiration_triggers (TEXT, JSON array)
- linked_resources (TEXT, JSON array)
- tags (TEXT, JSON array)
- status (TEXT)
- updated_at (DATETIME)
- retired_at (DATETIME)
- retirement_reason (TEXT)

### alerts
- id (TEXT, PRIMARY KEY)
- decision_id (TEXT, FOREIGN KEY)
- trigger_text (TEXT)
- evidence (TEXT)
- severity (TEXT)
- detected_at (DATETIME)
- acknowledged (BOOLEAN)

### decision_updates
- id (TEXT, PRIMARY KEY)
- decision_id (TEXT, FOREIGN KEY)
- updated_by (TEXT)
- updated_at (DATETIME)
- changes (TEXT)
- reason (TEXT)

## 🎬 Usage Flow

### 1. Initialize
```powershell
cd your-project
ddt init
```

### 2. Add Decisions
```powershell
ddt add "Use PostgreSQL for main database"
# Interactive prompts guide you through
```

### 3. Query Decisions
```powershell
ddt list                    # List all
ddt search "database"       # Search
ddt show ddt-2026-042       # Show details
```

### 4. Monitor Debt
```powershell
ddt check                   # Run detection
# Alerts shown for outdated decisions
```

### 5. Manage Decisions
```powershell
ddt update ddt-2026-042     # Update
ddt retire ddt-2026-042     # Retire
```

## 🔍 Detection Engine

The detector evaluates expiration triggers against the current state:

**Supported Trigger Patterns:**
- `"When data size exceeds 1GB"` → Monitors file sizes
- `"When we add authentication"` → Scans for keywords
- `"When more than 10 endpoints"` → Counts occurrences
- `"After 6 months"` → Time-based checks
- `"When configs require nesting"` → File structure analysis

**Detection Process:**
1. Retrieve all active decisions
2. For each decision, evaluate expiration triggers
3. Scan codebase for evidence
4. Generate alerts with context
5. Store alerts in database
6. Display summary with actionable insights

## 📦 Dependencies

```go
require (
    github.com/mattn/go-sqlite3 v1.14.18      // SQLite driver
    github.com/spf13/cobra v1.8.0             // CLI framework
    github.com/manifoldco/promptui v0.9.0     // Interactive prompts
    gopkg.in/yaml.v3 v3.0.1                   // YAML serialization
)
```

## 🚀 Building & Running

### Build
```powershell
go mod tidy
go build -o ddt.exe main.go
```

### Run
```powershell
.\ddt.exe --help
.\ddt.exe init
.\ddt.exe add "Your decision"
```

### Install (Optional)
```powershell
copy ddt.exe C:\tools\
# Add C:\tools to PATH
```

## 📝 File Organization

### Runtime Files
- `.ddt/decisions.db` - SQLite database (auto-created)
- `.ddt/backups/*.yaml` - Individual decision backups
- `.ddt/backups/backup-*.yaml` - Full database backups

### Source Files
- `cmd/*.go` - CLI command implementations
- `models/*.go` - Data structure definitions
- `storage/*.go` - Database operations
- `detector/*.go` - Detection engine logic

## 🎯 Design Principles

1. **Minimal Overhead** - 2 minutes or less to log a decision
2. **Frictionless Querying** - Fast search and retrieval
3. **Proactive Alerts** - Detect issues before they cause problems
4. **Machine-Readable** - Plain text, queryable, scriptable
5. **Ownership Tracking** - Link decisions to people and code
6. **Version Control Friendly** - YAML backups commit-able

## 🔮 Future Enhancements

### Planned Features
- [ ] Git integration (`ddt why <commit>`)
- [ ] VS Code extension
- [ ] Web dashboard
- [ ] Encrypted team sync
- [ ] Slack/Discord notifications
- [ ] Custom detection rules (DSL)
- [ ] Decision templates
- [ ] Impact analysis and visualization
- [ ] Markdown export for documentation
- [ ] Jira/GitHub issue integration

### Potential Improvements
- [ ] Unit tests and integration tests
- [ ] CI/CD examples (GitHub Actions, GitLab CI)
- [ ] Docker image
- [ ] Homebrew formula
- [ ] Chocolatey package
- [ ] Performance optimization for large codebases
- [ ] Parallel detection
- [ ] Incremental updates
- [ ] Plugin system

## 📚 Documentation Files

| File | Purpose |
|------|---------|
| README.md | Complete documentation and philosophy |
| QUICKSTART.md | Getting started guide |
| BUILD.md | Build and installation instructions |
| EXAMPLES.md | Real-world usage examples |
| LICENSE | MIT license terms |
| PROJECT.md | This file - project overview |

## 🎓 Learning Resources

### Understanding the Code
1. Start with `main.go` - entry point
2. Read `cmd/root.go` - CLI structure
3. Explore `cmd/add.go` - interactive prompts
4. Review `storage/storage.go` - database operations
5. Study `detector/detector.go` - detection logic

### Extending the Project
- **Add new commands:** Create new file in `cmd/`
- **New trigger types:** Extend `detector/detector.go`
- **Custom storage:** Implement interface in `storage/`
- **Export formats:** Add exporters to `storage/export.go`

## 🐛 Known Limitations

1. **CGO Requirement** - SQLite requires C compiler
2. **Local Only** - No cloud sync yet (planned)
3. **Single User** - No team features yet
4. **English Only** - No i18n support
5. **Basic Detection** - Heuristic-based, not AI-powered

## 💡 Philosophy

> "Every decision is a bet on the future. DDT tracks those bets and tells you when it's time to cash out."

DDT is NOT:
- ❌ A project management tool
- ❌ A wiki or knowledge base
- ❌ A comprehensive documentation system

DDT IS:
- ✅ A decision memory system
- ✅ An assumption tracker
- ✅ A technical debt detector
- ✅ Version control for choices

## 🤝 Contributing

Areas where contributions are welcome:
1. **Detection Rules** - Smarter trigger evaluation
2. **Integrations** - GitHub, Jira, Slack
3. **Visualization** - Decision graphs, timelines
4. **Testing** - Unit and integration tests
5. **Documentation** - Tutorials, videos
6. **Platform Support** - Package managers, installers

## 📬 Contact

- **GitHub:** github.com/nidhalelbkir/ddt (example)
- **Email:** nidhal@example.com
- **Issues:** Report bugs via GitHub Issues
- **Discussions:** Feature requests via GitHub Discussions

## 📊 Project Status

**Current Version:** 0.1.0 - MVP Complete

**Features Complete:**
- ✅ Core CLI commands
- ✅ SQLite storage
- ✅ YAML export
- ✅ Basic debt detection
- ✅ Interactive prompts
- ✅ Search and filtering

**In Progress:**
- 🔄 Unit tests
- 🔄 CI/CD examples
- 🔄 Advanced detection rules

**Planned:**
- 📅 Git integration
- 📅 VS Code extension
- 📅 Cloud sync

---

**Built with ❤️ by Nidhal Elkebir**  
**January 2026**
