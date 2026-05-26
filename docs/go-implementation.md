# SATOS - Go Implementation

**S**warm **A**gentic **T**eam **O**perating **S**ystem - Built with Go

## Architecture

```
SATOS/
├── cmd/
│   ├── director/     # Director entry point
│   └── tl/           # Team Leader dispatcher
├── internal/
│   ├── state/        # SQLite state management
│   ├── tasks/        # Task parsing & execution
│   ├── verifier/     # Verification system
│   └── recovery/     # Crash recovery
├── scripts/          # Build scripts
├── go.mod
├── go.sum
└── Makefile
```

## Building

### Prerequisites
- Go 1.22+
- SQLite3

### Build
```bash
make build
```

### Test
```bash
make test
```

## Commands

### Director
```bash
# Initialize SATOS
./bin/satos-director init

# Show status
./bin/satos-director status

# List tasks
./bin/satos-director tasks
```

### Team Leader (TL)
```bash
# Dispatch next task
./bin/satos-tl dispatch

# Run continuously
./bin/satos-tl run
```

## Project Structure

### Director
- `init` - Initialize SATOS state directory and database
- `status` - Show task statistics
- `tasks` - List all tasks

### Team Leader
- `dispatch` - Process next queued task
- `run` - Continuous processing mode

## State Management

SATOS uses SQLite for persistent state:
- **Tasks** - Track task status, assignments, metadata
- **Verification Logs** - Record verification results
- **Reliability Tracker** - Monitor verifier health
- **Recovery Log** - Track crash recovery events

## Testing

```bash
# Run all tests with coverage
make test-coverage

# Open coverage report
open coverage.html
```

## Development

```bash
# Format code
make fmt

# Lint code
make lint

# Install to PATH
make install
```