# Configuration

## CLI Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--dev-listen <addr>` | *(unset)* | Run in dev mode on specified address (e.g., `:8080`). Skips Tailscale. |
| `--hostname <name>` | `notebook` | Tailscale hostname |
| `--state-dir <dir>` | `tsnet-state` | Tailscale state directory |
| `--db <path>` | `notebook.db` | SQLite database file |

### Dev Mode

```bash
notebook --dev-listen :8080 --db notebook.db
```

Access at http://localhost:8080

### Tailscale Mode (Production)

```bash
notebook --hostname notebook --state-dir ./tsnet-state --db notebook.db
```

Access at https://notebook.your-tailnet.ts.net

**Note**: In Tailscale mode, the application is **only accessible via your Tailnet** (not localhost). Requires Tailscale authentication on first run.

## LLM Configuration

Configure via Web UI under "Configuration":

| Setting | Description | Example |
|---------|-------------|---------|
| **Provider URL** | Base URL for LLM API | `https://api.openai.com/v1` |
| **API Key** | Authentication key (masked after saving) | `sk-...` |
| **Model** | Model identifier | `gpt-4o`, `claude-opus-4-6` |
| **Summary Prompt** | Template for meeting summaries | Supports `{{subject}}`, `{{date}}`, `{{participants}}`, `{{notes}}` |
| **Enhancement Prompt** | Template for note enhancement | Supports `{{content}}` |

Configuration is stored in the SQLite database and persists across restarts. Supports OpenAI, Anthropic, Ollama, LM Studio, vLLM, and other OpenAI-compatible providers.

**Provider detection**: URLs containing `anthropic.com` use the Anthropic API; all others use the OpenAI-compatible API.

## LLM Features

### Meeting Summarization

- Click the ✨ icon in the meeting detail view
- Generates a concise summary (3-5 sentences) based on all notes
- Summary is saved to the meeting record
- Undo button (↶) appears after summarization to restore previous summary
- Requires at least one note in the meeting

### Note Enhancement

- Available in both note list view and note edit form
- Click the ✨ icon on any note to improve grammar, clarity, and structure
- Enhanced content replaces the original note
- Undo button (↶) appears to restore previous content
- Edit mode enhancement updates the textarea content without saving (allows further editing before save)

### Customizable Prompts

- Default prompts emphasize language preservation (no translation)
- Prompts instruct LLM to provide final text only (no multiple options)
- Edit prompts in Configuration panel to match your workflow
- Uses template syntax: `{{placeholder}}` for dynamic content
