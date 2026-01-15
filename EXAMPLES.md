# Decision Debt Tracker - Examples

This directory contains example workflows and use cases for DDT.

## Example 1: Database Technology Choice

```bash
ddt add "Use PostgreSQL for primary database"
```

**Interactive prompts:**
- Why? → "Need ACID compliance, JSON support, and mature ecosystem"
- Who decided? → "Backend Team"
- Assumptions:
  - Single datacenter deployment
  - Less than 10M rows
  - Team has PostgreSQL expertise
- Expiration triggers:
  - When we need multi-region deployment
  - When data exceeds 50M rows
  - When we need real-time sync
- Linked resources:
  - commit: a1b2c3d
  - doc: docs/architecture/database.md

## Example 2: Configuration Format

```bash
ddt add "Use JSON for configuration files"
```

**Details:**
- Why? → "Simple, widely supported, team familiar with JSON"
- Assumptions:
  - Configs remain flat (< 3 levels deep)
  - No need for comments
  - Under 100 config keys
- Expiration triggers:
  - When configs require nested sections
  - When we need schema validation
  - When configs exceed 200 lines

## Example 3: Authentication Strategy

```bash
ddt add "Use JWT for authentication"
```

**Details:**
- Why? → "Stateless, works with microservices, industry standard"
- Assumptions:
  - Token lifespan under 24 hours
  - No need for instant revocation
  - Single sign-on not required
- Expiration triggers:
  - When we need instant token revocation
  - When we implement SSO
  - When we have security audit requirements

## Checking Decisions

After development progresses:

```bash
$ ddt check

🔍 Checking decisions for expiration triggers...

⚠️  Found 1 potential issue:

⚠️ Decision [ddt-2025-042] may be outdated
   Title: Use JSON for configuration files
   Trigger: When configs require nested sections
   Evidence: Found keyword 'nested' in codebase
   Detected: 2025-01-15 14:30

💡 Review these decisions with: ddt show [decision-id]
```

## Updating a Decision

```bash
$ ddt update ddt-2025-042

📝 Updating decision: Use JSON for configuration files

What would you like to update?
> Assumptions
  Expiration Triggers
  All fields

Reason for update: → "Added nested config support, updated assumptions"

✅ Decision ddt-2025-042 updated
   Changes: Updated assumptions
   Reason: Added nested config support
```

## Retiring a Decision

```bash
$ ddt retire ddt-2025-042

Retirement reason: → "Migrated to YAML for better structure support"

✅ Decision ddt-2025-042 retired
   Reason: Migrated to YAML for better structure support
```

## Searching Decisions

```bash
# Find all database-related decisions
$ ddt search "database"

# Find all decisions by a specific person
$ ddt search "Backend Team"

# Find decisions with specific tags
$ ddt search "authentication"
```

## Listing Decisions

```bash
# Show recent active decisions
$ ddt list

# Show all decisions including retired
$ ddt list --all

# Show only retired decisions
$ ddt list --status retired

# Limit results
$ ddt list --limit 5
```

## Backing Up Decisions

```bash
# Export all to YAML
$ ddt export

# Create timestamped backup
$ ddt backup
```

## Git Integration Example

```bash
# Log a decision with commit reference
git commit -m "Add user authentication [DDT: Use JWT tokens]"

# Later, find decisions related to a feature
git log --grep="DDT" | grep -o "ddt-[0-9-]*"
```

## Real-World Scenarios

### Scenario 1: Tech Stack Choice

**Decision:** "Use React for frontend"
- **Assumption:** Team has React expertise
- **Trigger:** When team composition changes > 50%
- **Alert:** Triggered when 3+ new frontend devs join

### Scenario 2: Caching Strategy

**Decision:** "Use Redis for session cache"
- **Assumption:** Session data < 100MB
- **Trigger:** When session data exceeds 500MB
- **Alert:** Triggered when cache size > 400MB

### Scenario 3: API Design

**Decision:** "Use REST for public API"
- **Assumption:** No real-time requirements
- **Trigger:** When we need real-time updates
- **Alert:** Triggered when WebSocket code is added

## Tips

1. **Be Specific** - Vague triggers like "when needed" won't help
2. **Link Everything** - Commits, tickets, docs make context richer
3. **Review Weekly** - Run `ddt check` regularly
4. **Update, Don't Delete** - Retired decisions are still valuable history
5. **Tag Wisely** - Use consistent tags for easier searching

## CI/CD Example

```yaml
# .github/workflows/decision-check.yml
name: Decision Debt Check

on: [push]

jobs:
  check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Check Decision Debt
        run: |
          go install github.com/nidhalelbkir/ddt@latest
          ddt check || echo "⚠️ Review outdated decisions"
```
