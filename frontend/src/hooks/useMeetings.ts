import { useState, useEffect, useCallback } from 'react';
import type { Meeting } from '../api/types';
import { fetchMeetings, deleteMeeting } from '../api/client';

interface UseMeetingsResult {
  meetings: Meeting[];
  loading: boolean;
  error: string | null;
  sortColumn: string;
  sortOrder: string;
  handleSort: (column: string) => void;
  handleDelete: (id: number) => Promise<void>;
  refresh: () => Promise<void>;
}

export function useMeetings(): UseMeetingsResult {
  const [meetings, setMeetings] = useState<Meeting[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [sortColumn, setSortColumn] = useState('meeting_date');
  const [sortOrder, setSortOrder] = useState('desc');

  useEffect(() => {
    let cancelled = false;
    fetchMeetings(sortColumn, sortOrder)
      .then((data) => {
        if (!cancelled) {
          setMeetings(data);
          setError(null);
        }
      })
      .catch((err) => {
        if (!cancelled) setError(err instanceof Error ? err.message : 'Failed to load meetings');
      })
      .finally(() => {
        if (!cancelled) setLoading(false);
      });
    return () => { cancelled = true; };
  }, [sortColumn, sortOrder]);

  const handleSort = useCallback((column: string) => {
    setLoading(true);
    setSortColumn((prev) => {
      if (prev === column) {
        setSortOrder((order) => (order === 'asc' ? 'desc' : 'asc'));
        return prev;
      }
      setSortOrder('asc');
      return column;
    });
  }, []);

  const refresh = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await fetchMeetings(sortColumn, sortOrder);
      setMeetings(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load meetings');
    } finally {
      setLoading(false);
    }
  }, [sortColumn, sortOrder]);

  const handleDelete = useCallback(async (id: number) => {
    try {
      await deleteMeeting(id);
      await refresh();
    } catch (err) {
      throw new Error(err instanceof Error ? err.message : 'Failed to delete meeting');
    }
  }, [refresh]);

  return {
    meetings,
    loading,
    error,
    sortColumn,
    sortOrder,
    handleSort,
    handleDelete,
    refresh,
  };
}
