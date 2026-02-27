# Notebook

![CI](https://github.com/zorak1103/notebook/workflows/CI/badge.svg)
![Release](https://github.com/zorak1103/notebook/workflows/Release/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/zorak1103/notebook)](https://goreportcard.com/report/github.com/zorak1103/notebook)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
![Go Version](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go&logoColor=white)
[![Docker Pulls](https://img.shields.io/docker/pulls/zorak1103/notebook)](https://hub.docker.com/r/zorak1103/notebook)
![TypeScript](https://img.shields.io/badge/TypeScript-5.7-3178C6?logo=typescript&logoColor=white)
![React](https://img.shields.io/badge/React-19-61DAFB?logo=react&logoColor=black)
![Vite](https://img.shields.io/badge/Vite-6-646CFF?logo=vite&logoColor=white)
![i18n](https://img.shields.io/badge/i18n-en%20%7C%20de%20%7C%20fr%20%7C%20es-brightgreen)

Meeting and Notes Management Application with Tailscale Integration and optional LLM-powered summaries.

## Features

- **Internationalization**: Support for German, English, French, Spanish (react-i18next)
- **Meeting Management**: Create, edit, and organize meetings with metadata
- **Notes System**: Attach numbered notes to meetings with automatic numbering and manual reordering (▲/▼)
- **Full-Text Search**: Search across meeting subjects, summaries, participants, and keywords
- **Configuration Management**: Web UI for LLM provider settings with masked API keys and customizable prompts
- **Tailscale Integration**: Seamless authentication and secure network access via tsnet with user information display
- **LLM Integration**: Optional AI-powered meeting summaries and note enhancement with undo functionality (OpenAI, Anthropic, Ollama, LM Studio, vLLM)
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

**Dev Mode**:
```bash
docker run -d -p 8080:8080 -v $(pwd)/data:/data --name notebook zorak1103/notebook:latest
```
Open http://localhost:8080

**Tailscale Mode**:
```bash
docker run -d --network host -v $(pwd)/data:/data --name notebook \
  zorak1103/notebook:latest --hostname notebook --state-dir /data/tsnet-state --db /data/notebook.db
```
Access at https://notebook.your-tailnet.ts.net

**Docker Compose**: `docker compose up -d notebook-dev` (see `docker-compose.yml`)

### Binary Download

Download the latest release from [Releases](https://github.com/zorak1103/notebook/releases):

```bash
# Linux / macOS
tar xzf notebook_Linux_x86_64.tar.gz
./notebook --dev-listen :8080
```

Windows: download `notebook_Windows_x86_64.zip`.

### From Source

**Prerequisites**: Go 1.26+, Node.js 20+

```bash
git clone https://github.com/zorak1103/notebook.git
cd notebook
mage build  # or: go run github.com/magefile/mage@latest build
./notebook --dev-listen :8080
```

## Usage

```bash
# Dev mode (no Tailscale)
notebook --dev-listen :8080 --db notebook.db

# Tailscale mode (production)
notebook --hostname notebook --state-dir ./tsnet-state --db notebook.db
```

For CLI flags, LLM configuration, and AI feature details see [`.docs/configuration.md`](.docs/configuration.md).

## Development

| Command | Description |
|---------|-------------|
| `mage build` | Build frontend + backend |
| `mage test` | Run tests |
| `mage lint` | Run linter |
| `mage verify` | Run lint + test |
| `mage -l` | List all targets |

For project structure, build system details, and frontend dev setup see [`.docs/development.md`](.docs/development.md).

## Documentation

- [Configuration & LLM Setup](.docs/configuration.md)
- [Development Guide](.docs/development.md)
- [API Reference & DB Schema](.docs/api.md)
- [Product Requirements](.docs/prd.md)
- [Project Plan](.docs/projektplan.md)
- [Internationalization Guide](.docs/internationalization.md)
- [CI/CD Setup](.docs/cicd-setup.md)

## Architecture

See [`.docs/api.md`](.docs/api.md) for the full API reference and database schema.

Request flow: `Browser → HTTP → Server → Handler → Repository → SQLite`

## Multi-Platform Support

Binaries for Linux/macOS/Windows (amd64 + arm64) and Docker image `zorak1103/notebook:latest` (multi-arch).

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

## Acknowledgments

- [Tailscale](https://tailscale.com/) for tsnet
- [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite) for Pure Go SQLite
- [Vite](https://vitejs.dev/) for blazing fast frontend builds
- [GoReleaser](https://goreleaser.com/) for multi-platform releases
