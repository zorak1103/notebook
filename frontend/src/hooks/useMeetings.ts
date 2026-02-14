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

  const loadMeetings = useCallback(async () => {
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

  useEffect(() => {
    loadMeetings();
  }, [loadMeetings]);

  const handleSort = useCallback((column: string) => {
    setSortColumn((prev) => {
      if (prev === column) {
        // Toggle order if same column
        setSortOrder((order) => (order === 'asc' ? 'desc' : 'asc'));
        return prev;
      }
      // New column, default to ascending
      setSortOrder('asc');
      return column;
    });
  }, []);

  const handleDelete = useCallback(async (id: number) => {
    try {
      await deleteMeeting(id);
      await loadMeetings();
    } catch (err) {
      throw new Error(err instanceof Error ? err.message : 'Failed to delete meeting');
    }
  }, [loadMeetings]);

  return {
    meetings,
    loading,
    error,
    sortColumn,
    sortOrder,
    handleSort,
    handleDelete,
    refresh: loadMeetings,
  };
}
