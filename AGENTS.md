# AGENTS.md

Time logger CLI that syncs to Jira worklogs via SQLite database.

## Commands

Always use mise:
- `mise run build` - Build binary to `bin/tl`
- `mise run test` - Run tests with coverage
- `mise run lint` - Run golangci-lint

## Architecture

```
cmd/           # CLI commands (Cobra), one package per command
internal/
  api/         # Jira REST client
  db/          # SQLite repository (GORM)
  model/       # Domain models (TimeEntry, TimerEntry, Issue, etc.)
  service/     # Business logic (timer, fetch)
  alias/       # Issue key aliases (~/.config/tl/aliases.yml)
  util/        # Shared utilities
```

## Patterns

- **Dependency Injection**: Commands receive factory functions for lazy init
- **Interfaces**: All repos/services have interfaces with mocks in `*/mocks/`
- **Command Structure**: Each command package exports `NewCommand() *cobra.Command`

## Config & Data

- Config: `~/.config/tl/config.yml`
- Database: `~/.config/tl/db.sqlite`
- Aliases: `~/.config/tl/aliases.yml`
