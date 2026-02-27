# API Reference

## Database Schema

**SQLite** with automatic migrations (tracked via `schema_version` table).

### Tables

**`meetings`** — Meeting metadata

| Column | Type | Notes |
|--------|------|-------|
| id | INTEGER | Primary key |
| created_by | TEXT | Tailscale user identity |
| subject | TEXT | Meeting title |
| meeting_date | TEXT | Date (YYYY-MM-DD) |
| start_time | TEXT | Start time (HH:MM) |
| end_time | TEXT | End time (HH:MM) |
| participants | TEXT | Participant list |
| summary | TEXT | LLM-generated or manual summary |
| keywords | TEXT | Searchable keywords |
| created_at | DATETIME | Auto-set on insert |
| updated_at | DATETIME | Auto-set on update |

**`notes`** — Notes attached to meetings

| Column | Type | Notes |
|--------|------|-------|
| id | INTEGER | Primary key |
| meeting_id | INTEGER | FK → meetings(id) ON DELETE CASCADE |
| note_number | INTEGER | Auto-incremented per meeting |
| content | TEXT | Note body |
| created_at | DATETIME | Auto-set on insert |
| updated_at | DATETIME | Auto-set on update |

Unique constraint: `(meeting_id, note_number)`

**`config`** — Key-value configuration store

| Key | Description |
|-----|-------------|
| `llm_provider_url` | Base URL for LLM API |
| `llm_api_key` | API key (masked in responses) |
| `llm_model` | Model identifier |
| `language` | UI language code (en, de, fr, es) |
| `llm_prompt_summary` | Customizable summary prompt template |
| `llm_prompt_enhance` | Customizable note enhancement prompt template |

## API Endpoints

All endpoints are prefixed with `/api/`.

### Meetings

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/meetings` | List all meetings. Supports `?sort=meeting_date&order=desc` |
| `GET` | `/api/meetings/{id}` | Get meeting by ID |
| `POST` | `/api/meetings` | Create meeting |
| `PUT` | `/api/meetings/{id}` | Update meeting |
| `DELETE` | `/api/meetings/{id}` | Delete meeting |
| `POST` | `/api/meetings/{id}/summarize` | Generate AI summary from notes |

### Notes

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/meetings/{meetingId}/notes` | List notes for a meeting |
| `GET` | `/api/notes/{id}` | Get note by ID |
| `POST` | `/api/notes` | Create note (auto-assigns `note_number`) |
| `PUT` | `/api/notes/{id}` | Update note |
| `PUT` | `/api/notes/{id}/reorder` | Swap note order. Body: `{"direction": "up"\|"down"}`. Returns full updated note list. |
| `DELETE` | `/api/notes/{id}` | Delete note |
| `POST` | `/api/notes/{id}/enhance` | Enhance note content with AI |

### Search

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/search?q=<query>` | Search meetings across subject, summary, participants, and keywords |

### Configuration

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/config` | Get all configuration (API keys masked) |
| `POST` | `/api/config` | Update configuration |

### Authentication

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/whoami` | Get Tailscale user identity (WhoIs API) |

## Frontend Architecture

- Single-page application (SPA) — React + Vite + TypeScript
- Vite build output embedded in Go binary via `go:embed` (path: `internal/web/frontend/dist/`)
- Served from `/` with SPA routing fallback
- API calls proxied through Vite dev server in development (`vite.config.ts`)
