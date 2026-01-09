# AGENTS.md - Coding Guidelines for euro21

This document provides guidelines for AI coding agents working on the euro21 codebase.

## Project Overview

euro21 is a Go-based CLI tool for displaying Euro 2021 tournament information (tables, matches, standings) by fetching data from the Livescore API. Built with Cobra for CLI commands and custom table rendering.

**Tech Stack**: Go 1.16+, Cobra, promptui, emoji utilities

## Build, Test & Lint Commands

### Building
```bash
# Build the binary
make build
# or
go build -o euro21

# Run without building
make run
# or
go run main.go

# Download dependencies
make prepare
# or
go mod download
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests for a specific package
go test ./cmd
go test ./adapters
go test ./utils/stringutil

# Run a single test function
go test ./cmd -run TestFunctionName

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Linting
```bash
# No official linter configured in project, but recommended:
golangci-lint run
go vet ./...
gofmt -s -w .
```

### CI/CD
- GitHub Actions workflow: `.github/workflows/go.yml`
- Runs on: `master`, `develop`, `fix/*` branches
- CI command: `go build -v .`

## Project Structure

```
euro21/
├── cmd/          # Cobra commands (root, table, match)
├── adapters/     # External API client (Livescore)
├── constants/    # Color codes, stage definitions
├── dtos/         # Data transfer objects (API responses)
├── tablewriter/  # Custom table rendering logic
└── utils/        # Helper functions (slice, string utilities)
```

## Code Style Guidelines

### Package Organization
- **cmd/**: CLI command definitions and handlers
- **adapters/**: External API integrations
- **dtos/**: JSON response structures
- **constants/**: Shared constants
- **utils/**: Pure utility functions (no side effects)

### Imports
- Group imports: stdlib → third-party → local
- Use blank imports with comments: `_ "embed"`
- Import aliases for clarity: `emojiflags "github.com/yudgnahk/go-emoji-flags"`

**Example**:
```go
import (
    "encoding/json"
    "fmt"
    "os"
    
    _ "embed"
    
    "github.com/spf13/cobra"
    "github.com/yudgnahk/euro21/dtos"
)
```

### Naming Conventions
- **Packages**: lowercase, single word (e.g., `cmd`, `adapters`, `dtos`)
- **Files**: lowercase with underscores (e.g., `root.go`, `multi_table.go`)
- **Exported**: PascalCase (e.g., `GetTable`, `RootCmd`, `TableData`)
- **Unexported**: camelCase (e.g., `countriesMap`, `getColor`, `renderHeaders`)
- **Constants**: PascalCase for exported, camelCase for unexported
- **Acronyms**: Keep uppercase in names (e.g., `MatchID`, `TBD`)

### Types & Structs
- Define types for clarity (e.g., `type Color string`)
- Use JSON struct tags for API responses: `` `json:"fieldName"` ``
- Embed structs where appropriate
- Use pointer receivers for methods that modify state

**Example**:
```go
type Team struct {
    Nm   string `json:"Nm"`
    ID   string `json:"ID"`
    Tbd  int    `json:"tbd"`
}
```

### Functions
- Exported functions: Start with capital letter, include docstring
- Keep functions focused and small
- Return errors explicitly, don't ignore them
- Use descriptive parameter names

**Example**:
```go
// GetTables fetches Euro 2021 group stage tables from Livescore API
func GetTables() (*dtos.TableData, error) {
    request, _ := newGetRequest(fmt.Sprintf("%v/%v", Host, TablePath))
    setHeaders(request)
    
    var response dtos.TableData
    err := execute(request, &response)
    return &response, err
}
```

### Error Handling
- **Current pattern**: Some errors are ignored with `_`, some logged with fmt.Println
- **Best practice**: Always handle errors, return them up the stack
- Use `logrus` for logging errors in cmd handlers
- Don't use `panic()` unless truly unrecoverable

**Current pattern in codebase**:
```go
data, _ := adapters.GetTables()  // Error ignored
```

**Preferred pattern**:
```go
data, err := adapters.GetTables()
if err != nil {
    logrus.Error(err)
    return err
}
```

### Variables
- Use `:=` for short variable declarations
- Declare and initialize separately when type clarity needed
- Use `make()` for slices/maps with known capacity
- Zero values: prefer `var foo []string` over `foo := []string{}`

### Constants
- Group related constants in const blocks
- Use typed constants: `const Host = "https://..."`
- Define in `constants/` package for shared values

### Comments
- Exported functions: Always add docstrings
- Unexported functions: Add comments for complex logic
- Inline comments: Explain "why", not "what"
- No redundant comments (e.g., don't comment obvious code)

### Formatting
- Use `gofmt` (standard Go formatting)
- Line length: No hard limit, but keep readable
- Indentation: Tabs (Go standard)
- Blank lines: Separate logical sections

## API Integration

### HTTP Requests
- Base URL: `https://prod-public-api.livescore.com/v1/api/react`
- Always set required headers (User-Agent, Referer, etc.)
- Use helper functions: `newGetRequest()`, `setHeaders()`, `execute()`
- Close response bodies with `defer res.Body.Close()`

### Data Flow
1. `adapters/` fetches raw API data
2. Unmarshal into `dtos/` structs
3. `cmd/` processes and displays data
4. `tablewriter/` handles terminal rendering

## UI/UX Patterns

### Terminal Colors
- Use `constants.Color` type for ANSI codes
- Green: Priority/qualified teams
- Yellow: Conditional (e.g., best third place)
- White: Default/eliminated teams

### Table Rendering
- Use custom `tablewriter` package (not external library)
- Supports emoji flags, colored rows, multi-table layouts
- Handle emoji width correctly in `utils/stringutil`

### Interactive Prompts
- Use `promptui` for menus and selections
- Include emoji icons (⚽️, 👉🏻) for visual appeal
- Make prompts searchable

## Common Patterns

### Embedded Resources
```go
//go:embed countries.json
var countriesData string
```

### Slice Utilities
```go
sliceutil.ToStringSlice(flagAndName, team.Points, team.Win)
```

### Time Parsing
- API returns timestamps as int64 (format: YYYYMMDDHHmmss)
- Use custom `getTime()` function in `cmd/match.go`

## Testing Guidelines

- No tests currently exist in the project
- When adding tests:
  - Place in same package: `*_test.go`
  - Use table-driven tests for multiple cases
  - Mock API responses for adapter tests
  - Test edge cases (empty responses, errors)

## Git Workflow

- Main branch: `master`
- Development: `develop`
- Bug fixes: `fix/*` branches
- CI runs on push to these branches

## Dependencies

- `github.com/spf13/cobra`: CLI framework
- `github.com/manifoldco/promptui`: Interactive prompts
- `github.com/sirupsen/logrus`: Logging
- `github.com/tmdvs/Go-Emoji-Utils`: Emoji parsing
- `github.com/yudgnahk/go-emoji-flags`: Flag emojis

## Notes for Agents

- This is a legacy project (Euro 2021), API may not return live data
- Error handling is minimal; improve it when modifying code
- No tests exist; consider adding them for new features
- Custom table rendering supports emojis; test terminal compatibility
- API headers are required for Livescore to respond properly
