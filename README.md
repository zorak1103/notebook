# Notebook

![CI](https://github.com/zorak1103/notebook/workflows/CI/badge.svg)
![Release](https://github.com/zorak1103/notebook/workflows/Release/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/zorak1103/notebook)](https://goreportcard.com/report/github.com/zorak1103/notebook)

Meeting and Notes Management Application with Tailscale Integration and optional LLM-powered summaries.

## Features

- **Internationalization**: Support for German, English, French, Spanish (react-i18next)
- **Meeting Management**: Create, edit, and organize meetings with metadata
- **Notes System**: Attach numbered notes to meetings with automatic numbering
- **Full-Text Search**: Search across meeting subjects and summaries
- **Configuration Management**: Web UI for LLM provider settings with masked API keys
- **Tailscale Integration**: Seamless authentication and secure network access via tsnet
- **LLM Integration**: Optional AI-powered meeting summaries (OpenAI, Anthropic)
- **Single Binary**: Frontend embedded using go:embed
- **Dev Mode**: Run without Tailscale for local development
- **Responsive UI**: React + Vite + TypeScript frontend

## Tech Stack

- **Backend**: Go 1.26+ with stdlib HTTP routing
- **Frontend**: React + Vite + TypeScript
- **Internationalization**: react-i18next (German, English, French, Spanish)
- **Database**: SQLite with modernc.org/sqlite (Pure Go, no CGO)
- **Auth**: Tailscale tsnet integration
- **Embedding**: go:embed for single binary deployment

## Quick Start

### Docker (Recommended)

```bash
docker pull zorak1103/notebook:latest

docker run -d \
  -p 8080:8080 \
  -v $(pwd)/data:/data \
  --name notebook \
  zorak1103/notebook:latest
```

Open http://localhost:8080

### Binary Download

Download the latest release for your platform:

**Linux**:
```bash
wget https://github.com/zorak1103/notebook/releases/latest/download/notebook_Linux_x86_64.tar.gz
tar xzf notebook_Linux_x86_64.tar.gz
./notebook --dev-listen :8080
```

**macOS**:
```bash
wget https://github.com/zorak1103/notebook/releases/latest/download/notebook_Darwin_x86_64.tar.gz
tar xzf notebook_Darwin_x86_64.tar.gz
./notebook --dev-listen :8080
```

**Windows**:
Download `notebook_Windows_x86_64.zip` from [Releases](https://github.com/zorak1103/notebook/releases)

### From Source

**Prerequisites**:
- Go 1.26+
- Node.js 20+
- Mage (optional: `go install github.com/magefile/mage@latest`)

**Build**:
```bash
git clone https://github.com/zorak1103/notebook.git
cd notebook
mage build  # or: go run github.com/magefile/mage@latest build
./notebook --dev-listen :8080
```

## Usage

### Development Mode (without Tailscale)

```bash
notebook --dev-listen :8080 --db notebook.db
```

Access at http://localhost:8080

### Tailscale Mode (Production)

```bash
notebook --hostname notebook --state-dir ./tsnet-state --db notebook.db
```

Access at https://notebook.your-tailnet.ts.net

### Configuration

**CLI Flags**:
- `--dev-listen <addr>` - Run in dev mode on specified address (e.g., `:8080`)
- `--hostname <name>` - Tailscale hostname (default: `notebook`)
- `--state-dir <dir>` - Tailscale state directory (default: `tsnet-state`)
- `--db <path>` - SQLite database file (default: `notebook.db`)

**LLM Configuration**:

Configure via Web UI under "Configuration":
- **Provider URL**: Base URL for LLM API (e.g., `https://api.openai.com/v1`)
- **API Key**: Authentication key (automatically masked after saving)
- **Model**: Model identifier (e.g., `gpt-4o`, `claude-opus-4-6`)

Configuration is stored in the SQLite database and persists across restarts.

## Development

### Project Structure

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
├── .docs/                # Documentation
├── .github/workflows/    # CI/CD pipelines
└── magefile.go           # Build system (Mage)
```

### Build Commands

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

### Running Tests

```bash
go test -v ./...
go test -cover ./...
```

### Frontend Development

```bash
cd frontend
npm install
npm run dev       # Vite dev server on :5173
```

In another terminal:
```bash
go run ./cmd/notebook --dev-listen :8080
```

Frontend proxies API requests to backend (configured in `vite.config.ts`)

## Documentation

- [Product Requirements](/.docs/prd.md)
- [Project Plan](/.docs/projektplan.md)
- [Internationalization Guide](/.docs/internationalization.md)
- [CI/CD Setup](/.docs/cicd-setup.md)
- [Phase Documentation](/.docs/)

## Architecture

**Database**: SQLite with automatic migrations

**Tables**:
- `meetings` - Meeting metadata with full-text search indices
- `notes` - Notes with auto-incrementing note numbers per meeting
- `config` - Key-value configuration store

**API Endpoints**:

*Meetings*:
- `GET /api/meetings` - List all meetings (supports `?sort=` and `?order=`)
- `GET /api/meetings/{id}` - Get meeting by ID
- `POST /api/meetings` - Create meeting
- `PUT /api/meetings/{id}` - Update meeting
- `DELETE /api/meetings/{id}` - Delete meeting

*Notes*:
- `GET /api/meetings/{meetingId}/notes` - List notes for meeting
- `GET /api/notes/{id}` - Get note by ID
- `POST /api/notes` - Create note (auto-assigns note_number)
- `PUT /api/notes/{id}` - Update note
- `DELETE /api/notes/{id}` - Delete note

*Search*:
- `GET /api/search?q=<query>` - Search meetings by subject/summary

*Configuration*:
- `GET /api/config` - Get configuration (API keys masked)
- `POST /api/config` - Update configuration

*Authentication*:
- `GET /api/whoami` - Get Tailscale user identity (WhoIs API)

**Frontend**:
- Single-page application (SPA)
- Vite build output embedded in Go binary via `go:embed`
- Served from `/` with SPA routing fallback

## Multi-Platform Support

**Supported Platforms**:
- Linux: amd64, arm64
- macOS: amd64, arm64 (Apple Silicon)
- Windows: amd64, arm64

**Docker Images**:
- `zorak1103/notebook:latest` (multi-arch: amd64, arm64)
- `zorak1103/notebook:v1.0.0` (version tags)

**Package Formats**:
- DEB (Debian/Ubuntu)
- RPM (RedHat/Fedora)
- APK (Alpine)
- Arch Linux packages

## Security

- **Authentication**: Automatic via Tailscale WhoIs API
- **Network Isolation**: Only accessible on your Tailnet (Tailscale mode)
- **Database**: SQLite file-based, no network exposure
- **API Keys**: Masked in UI, stored in database (encryption optional)
- **Non-root Docker**: Runs as user `notebook` (UID 1000)

## License

MIT License - see [LICENSE](LICENSE) for details

## Contributing

Contributions welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) first.

## Support

- **Issues**: https://github.com/zorak1103/notebook/issues
- **Discussions**: https://github.com/zorak1103/notebook/discussions

## Roadmap

**Phase 1-2 (Completed)**:
- [x] Project scaffolding with Tailscale integration
- [x] Database layer (SQLite with modernc.org/sqlite)
- [x] Schema migrations system
- [x] Repository pattern (Meeting, Note, Config)

**Phase 3-6 (Completed)**:
- [x] Meeting CRUD API endpoints
- [x] Notes CRUD API endpoints with auto-numbering
- [x] Full-text search API
- [x] Configuration management API with masked API keys
- [x] React frontend with i18n (DE, EN, FR, ES)

**Phase 7 (In Progress)**:
- [ ] LLM integration (OpenAI, Anthropic)
- [ ] Meeting summarization endpoints

**Future Enhancements**:
- [ ] Meeting templates
- [ ] Markdown support in notes
- [ ] Export to PDF
- [ ] Mobile app
- [ ] Real-time collaboration

## Acknowledgments

- [Tailscale](https://tailscale.com/) for tsnet
- [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite) for Pure Go SQLite
- [Vite](https://vitejs.dev/) for blazing fast frontend builds
- [GoReleaser](https://goreleaser.com/) for multi-platform releases
