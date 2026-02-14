import { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import type { Meeting } from '../api/types';
import { searchMeetings } from '../api/client';
import { useDebounce } from '../hooks/useDebounce';
import './SearchPanel.css';

interface SearchPanelProps {
  onSelectMeeting: (id: number) => void;
}

function escapeRegex(str: string): string {
  return str.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
}

function highlightMatch(text: string, query: string): React.JSX.Element {
  if (!query.trim()) {
    return <>{text}</>;
  }

  const escapedQuery = escapeRegex(query);
  const regex = new RegExp(`(${escapedQuery})`, 'gi');
  const parts = text.split(regex);

  return (
    <>
      {parts.map((part, i) =>
        regex.test(part) ? <mark key={i}>{part}</mark> : part
      )}
    </>
  );
}

function truncate(text: string | null, maxLength: number): string {
  if (!text) return '';
  if (text.length <= maxLength) return text;
  return text.slice(0, maxLength) + '...';
}

export function SearchPanel({ onSelectMeeting }: SearchPanelProps) {
  const { t } = useTranslation();
  const [query, setQuery] = useState('');
  const [results, setResults] = useState<Meeting[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const debouncedQuery = useDebounce(query, 300);

  useEffect(() => {
    if (!debouncedQuery.trim()) {
      setResults([]);
      setError(null);
      return;
    }

    setLoading(true);
    setError(null);

    searchMeetings(debouncedQuery)
      .then((data) => {
        setResults(data);
        setLoading(false);
      })
      .catch((err) => {
        setError(err.message || t('search.error'));
        setLoading(false);
      });
  }, [debouncedQuery, t]);

  return (
    <div className="search-panel">
      <h2>{t('search.title')}</h2>

      <div className="search-input-container">
        <input
          type="text"
          className="search-input"
          placeholder={t('search.placeholder')}
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          autoFocus
        />
      </div>

      {loading && <div className="search-loading">{t('search.loading')}</div>}

      {error && <div className="search-error">{error}</div>}

      {!loading && !error && query.trim() && results.length === 0 && (
        <div className="search-no-results">{t('search.noResults')}</div>
      )}

      {!loading && !error && results.length > 0 && (
        <>
          <div className="search-result-count">
            {t('search.resultCount', { count: results.length })}
          </div>
          <div className="search-results">
            {results.map((meeting) => (
              <div
                key={meeting.id}
                className="search-result-item"
                onClick={() => onSelectMeeting(meeting.id)}
              >
                <div className="search-result-subject">
                  {highlightMatch(meeting.subject, query)}
                </div>
                <div className="search-result-meta">
                  {meeting.meeting_date} {meeting.start_time}
                  {meeting.end_time && ` - ${meeting.end_time}`}
                </div>
                {meeting.summary && (
                  <div className="search-result-summary">
                    {highlightMatch(truncate(meeting.summary, 150), query)}
                  </div>
                )}
              </div>
            ))}
          </div>
        </>
      )}
    </div>
  );
}
