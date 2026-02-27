# Internationalization (i18n)

## Supported Languages

| Code | Language |
|------|----------|
| `en` | English (default fallback) |
| `de` | German |
| `fr` | French |
| `es` | Spanish |

Language preference is persisted in the backend database and restored on startup. Browser language is detected automatically on first visit.

## Implementation

**Libraries**: `react-i18next@^15`, `i18next@^23`, `i18next-http-backend@^2.5`

**Translation files**: `frontend/public/locales/{en,de,fr,es}/translation.json`

**i18next config**: `frontend/src/i18n.ts`

### Usage in Components

```typescript
import { useTranslation } from 'react-i18next';

function MyComponent() {
  const { t } = useTranslation();
  return <h2>{t('meetings.title')}</h2>;
}
```

Interpolation: `t('search.resultCount', { count: 5 })`
Plurals: use `key_other` suffix — `"resultCount_other": "{{count}} results"`

### Translation Key Structure

Keys are grouped by feature:

```json
{
  "app": { "title": "Notebook", "loading": "Loading..." },
  "meetings": { "title": "Meetings", "create": "New Meeting", ... },
  "notes": { "title": "Notes", "add": "Add Note", ... },
  "search": { "placeholder": "Search meetings...", ... },
  "config": { "title": "Configuration", ... },
  "info": { "title": "User Info", ... }
}
```

Use semantic keys (`meetings.create`, not `btn1`). Never concatenate translated strings.

### Language Switcher

The `LanguageSwitcher` component owns language persistence — it saves the selected language to the backend config. Other components must not write the `language` field to avoid overwriting the user's selection.

Language comparison uses `i18n.language.startsWith(code)` (not strict equality) because i18next resolves to BCP 47 tags like `"en-US"` while the backend stores short codes like `"en"`.

## Adding New Strings

1. Add key to `frontend/public/locales/en/translation.json`
2. Add the same key to `de`, `fr`, and `es` — missing keys cause UI errors
3. Use `t()` in components — never hardcode UI strings

## Backend

All backend error messages and API responses are in English. User-generated content (meeting subjects, notes) is stored as entered and never translated.
