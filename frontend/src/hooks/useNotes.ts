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

  const loadNotes = useCallback(async () => {
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

  useEffect(() => {
    loadNotes();
  }, [loadNotes]);

  const handleDelete = useCallback(async (id: number) => {
    try {
      await deleteNote(id);
      await loadNotes();
    } catch (err) {
      throw new Error(err instanceof Error ? err.message : 'Failed to delete note');
    }
  }, [loadNotes]);

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
    refresh: loadNotes,
  };
}
