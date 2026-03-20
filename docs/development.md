# Development

## Project Structure

```
notebook/
├── cmd/notebook/          # Main entry point
├── internal/
│   ├── db/               # Database layer (SQLite)
│   ├── llm/              # LLM integration
│   ├── tsapp/            # Tailscale wrapper
│   ├── validation/       # Generated validation rules
│   └── web/              # HTTP server & handlers
├── frontend/             # React + Vite frontend
├── docs/                 # Documentation
├── .github/workflows/    # CI/CD pipelines
└── Taskfile.yml          # Build system (Task)
```

## Build Commands

The project uses [Task](https://taskfile.dev) (YAML-based Make alternative).

```bash
task build              # Build frontend + backend (default)
task build:frontend     # Build frontend only
task build:check        # Compile-check Go without building frontend
task dev                # Show dev setup instructions
task dev:backend        # Start Go backend in dev mode
task dev:frontend       # Start Vite dev server
task test               # Run tests
task test:coverage      # Run tests with per-file 80% coverage enforcement
task lint               # Run golangci-lint
task lint:frontend      # Run ESLint
task verify             # Run lint + lint:frontend + test
task clean              # Clean build artifacts
task --list             # List all available tasks
```

Install Task: `go install github.com/go-task/task/v3/cmd/task@latest` (or via Scoop/Homebrew/apt).

**Important**: `go:embed` requires the frontend to be built before any Go command (build, test, lint). Always run `task build:frontend` first. For Go-only work (e.g. unit tests without frontend), `task build:check` and `task test` use `prepare:embed` to create a placeholder directory automatically.

## Running Tests

```bash
go test -v ./...
go test -cover ./...

# Run a specific test
go test -v ./internal/db/repositories -run TestMeetingRepository

# Run handler tests only
go test -v ./internal/web/... -run TestHandle
```

The linter requires a timeout flag to avoid slow linters being skipped:
```bash
golangci-lint run --timeout=5m ./...
```

## Frontend Development

```bash
cd frontend
npm install
npm run dev       # Vite dev server on :5173
```

In another terminal:
```bash
go run ./cmd/notebook --dev-listen :8080
```

Frontend proxies API requests to backend (configured in `vite.config.ts`).

**Build output**: Vite builds to `internal/web/frontend/dist/` (go:embed constraint — NOT `frontend/dist/`).

## Prerequisites

- Go 1.26+
- Node.js 20+
- Task: `go install github.com/go-task/task/v3/cmd/task@latest` (or via Scoop/Homebrew/apt)
- golangci-lint v2

## CI/CD

GitHub Actions workflows in `.github/workflows/`:
- **ci.yml**: Lint, test, build (runs on every push/PR)
- **release.yml**: Multi-platform builds + Docker images (runs on `v*` tags)

Release process:
```bash
git tag v1.0.0
git push origin v1.0.0
# GitHub Actions builds and publishes automatically
```

