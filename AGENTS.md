# Sprayer Development Guidelines

## Build Commands

```bash
# Build the main application
go build -o sprayer .

# Build with specific entry points
go build -o sprayer-cli cmd/cli/main.go
go build -o sprayer-tui cmd/tui/main.go
go build -o sprayer-server cmd/server/main.go
```

## Test Commands

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/job/...
go test ./internal/parse/...
go test ./internal/profile/...

# Run a single test function
go test -run TestPipeline_Composition ./internal/job/

# Run end-to-end tests
go test ./tests/e2e/...
```

## Code Style Guidelines

### Imports
- Standard library imports first
- Third-party imports second
- Local imports last with project prefix "sprayer/"
- Use blank lines to separate import groups

Example:
```go
import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"

	"sprayer/internal/api"
	"sprayer/internal/job"
)
```

### Naming Conventions
- **Packages**: Single word, lowercase (e.g., `job`, `parse`, `profile`)
- **Types**: PascalCase (e.g., `Job`, `Pipeline`, `Scraper`)
- **Functions**: PascalCase for exported, camelCase for unexported
- **Variables**: camelCase for unexported, PascalCase for exported
- **Constants**: PascalCase or camelCase depending on export
- **Interfaces**: End with "-er" suffix when appropriate (e.g., `Scraper`, `Filter`)

### Error Handling
- Return errors explicitly, don't panic
- Use descriptive error messages
- Wrap errors with context using `fmt.Errorf`
- Check errors immediately after function calls

Example:
```go
jobs, err := scraper.Scrape()
if err != nil {
	return fmt.Errorf("failed to scrape jobs: %w", err)
}
```

### Types and Structs
- Use JSON tags for all exported struct fields
- Include omitempty for optional fields
- Keep structs focused and cohesive

Example:
```go
type Job struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Company     string    `json:"company"`
	Salary      string    `json:"salary,omitempty"`
	HasTraps    bool      `json:"has_traps"`
}
```

### Function Design
- Keep functions small and focused
- Use functional patterns where appropriate (filters, mappers)
- Document exported functions
- Return concrete types when possible

### Concurrency
- Use goroutines for parallel operations
- Channel-based communication preferred
- Handle context cancellation properly

### Database
- SQLite3 for data persistence
- Use prepared statements for queries
- Handle connection lifecycle properly

### Dependencies
- Key libraries: Bubble Tea (TUI), Chi (HTTP router), Rod (browser automation)
- Keep dependencies minimal and well-maintained
- Use Go 1.24+ features appropriately

### Testing
- Write unit tests for core logic
- Use table-driven tests for multiple scenarios
- Mock external dependencies
- Include integration tests for critical paths

### Project Structure
```
sprayer/
├── cmd/           # Application entry points
│   ├── cli/       # CLI version
│   ├── tui/       # TUI version
│   └── server/    # HTTP server version
├── internal/      # Private application code
│   ├── api/       # HTTP API handlers
│   ├── job/       # Job processing logic
│   ├── parse/     # Content parsing
│   ├── profile/   # User profiles
│   └── ui/        # User interface components
└── tests/e2e/     # End-to-end tests
```

## Environment
- Use `.env` files for configuration
- Load environment variables with `godotenv`
- Support both development and production configurations