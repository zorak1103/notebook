import { useState, useEffect, useCallback } from 'react';
import type { Note } from '../api/types';
import { fetchNotes, deleteNote, reorderNote } from '../api/client';

interface UseNotesResult {
  notes: Note[];
  loading: boolean;
  error: string | null;
  handleDelete: (id: number) => Promise<void>;
  handleReorder: (id: number, direction: 'up' | 'down') => Promise<void>;
  refresh: () => Promise<void>;
}

export function useNotes(meetingId: number): UseNotesResult {
  const [notes, setNotes] = useState<Note[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;
    fetchNotes(meetingId)
      .then((data) => {
        if (!cancelled) {
          setNotes(data);
          setError(null);
        }
      })
      .catch((err) => {
        if (!cancelled) setError(err instanceof Error ? err.message : 'Failed to load notes');
      })
      .finally(() => {
        if (!cancelled) setLoading(false);
      });
    return () => { cancelled = true; };
  }, [meetingId]);

  const refresh = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await fetchNotes(meetingId);
      setNotes(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load notes');
    } finally {
      setLoading(false);
    }
  }, [meetingId]);

  const handleDelete = useCallback(async (id: number) => {
    try {
      await deleteNote(id);
      await refresh();
    } catch (err) {
      throw new Error(err instanceof Error ? err.message : 'Failed to delete note');
    }
  }, [refresh]);

  const handleReorder = useCallback(async (id: number, direction: 'up' | 'down') => {
    try {
      const updated = await reorderNote(id, direction);
      setNotes(updated);
    } catch (err) {
      throw new Error(err instanceof Error ? err.message : 'Failed to reorder note');
    }
  }, []);

  return {
    notes,
    loading,
    error,
    handleDelete,
    handleReorder,
    refresh,
  };
}
