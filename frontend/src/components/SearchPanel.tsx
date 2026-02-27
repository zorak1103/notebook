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
      {parts.map((part, i) => {
        // When split with capturing group, odd indices are matches
        const isMatch = i % 2 === 1;
        return isMatch ? <mark key={i}>{part}</mark> : part;
      })}
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
      return;
    }

    async function fetchResults() {
      setLoading(true);
      setError(null);
      try {
        const data = await searchMeetings(debouncedQuery);
        setResults(data);
      } catch (err) {
        setError((err as Error).message || t('search.error'));
      } finally {
        setLoading(false);
      }
    }

    fetchResults();
  }, [debouncedQuery, t]);

  const hasQuery = debouncedQuery.trim().length > 0;
  const visibleResults = hasQuery ? results : [];
  const visibleError = hasQuery ? error : null;

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

      {visibleError && <div className="search-error">{visibleError}</div>}

      {!loading && !visibleError && query.trim() && visibleResults.length === 0 && (
        <div className="search-no-results">{t('search.noResults')}</div>
      )}

      {!loading && !visibleError && visibleResults.length > 0 && (
        <>
          <div className="search-result-count">
            {t('search.resultCount', { count: visibleResults.length })}
          </div>
          <div className="search-results">
            {visibleResults.map((meeting) => (
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
                {meeting.participants && (
                  <div className="search-result-participants">
                    👥 {highlightMatch(meeting.participants, query)}
                  </div>
                )}
                {meeting.keywords && (
                  <div className="search-result-keywords">
                    🏷 {highlightMatch(meeting.keywords, query)}
                  </div>
                )}
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
