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
└── magefile.go           # Build system (Mage)
```

## Build Commands

The project uses [Mage](https://magefile.org) (Go-based Make alternative).

```bash
mage build        # Build frontend + backend (or just: mage)
mage frontend     # Build frontend only
mage backend      # Build backend only
mage dev          # Show dev setup instructions
mage test         # Run tests
mage lint         # Run linter
mage verify       # Run lint + test
mage clean        # Clean build artifacts
mage -l           # List all available targets
```

Mage without installation (useful for CI/CD or first-time setup):
```bash
go run github.com/magefile/mage@latest build
go run github.com/magefile/mage@latest -l
```

**Important**: `go:embed` requires the frontend to be built before any Go command (build, test, lint). Always run `mage frontend` first.

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
- Mage: `go install github.com/magefile/mage@latest` (or via Scoop/Homebrew)
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

