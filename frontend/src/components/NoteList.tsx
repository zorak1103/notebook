import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useNotes } from '../hooks/useNotes';
import { enhanceNote, updateNote } from '../api/client';
import { LoadingSpinner } from './LoadingSpinner';
import { ErrorMessage } from './ErrorMessage';
import './NoteList.css';

interface NoteListProps {
  meetingId: number;
  onEdit: (noteId: number) => void;
  onAdd: () => void;
}

export function NoteList({ meetingId, onEdit, onAdd }: NoteListProps) {
  const { t } = useTranslation();
  const { notes, loading, error, handleDelete, refresh } = useNotes(meetingId);
  const [enhancingId, setEnhancingId] = useState<number | null>(null);
  const [enhanceError, setEnhanceError] = useState<string | null>(null);
  const [previousContent, setPreviousContent] = useState<{ noteId: number; content: string } | null>(null);

  const confirmDelete = async (id: number, noteNumber: number) => {
    if (window.confirm(t('notes.confirmDelete', { number: noteNumber }))) {
      try {
        await handleDelete(id);
      } catch (err) {
        alert(err instanceof Error ? err.message : t('notes.deleteFailed'));
      }
    }
  };

  const handleEnhance = async (noteId: number) => {
    const note = notes.find((n) => n.id === noteId);
    if (!note) return;

    try {
      setEnhancingId(noteId);
      setEnhanceError(null);
      setPreviousContent({ noteId, content: note.content });

      await enhanceNote(noteId);
      await refresh();
    } catch (err) {
      setEnhanceError(err instanceof Error ? err.message : t('notes.enhanceError'));
      setPreviousContent(null);
    } finally {
      setEnhancingId(null);
    }
  };

  const handleUndoEnhance = async () => {
    if (!previousContent) return;

    try {
      setEnhancingId(previousContent.noteId);
      setEnhanceError(null);

      await updateNote(previousContent.noteId, { content: previousContent.content });
      await refresh();
      setPreviousContent(null);
    } catch (err) {
      setEnhanceError(err instanceof Error ? err.message : t('notes.enhanceError'));
    } finally {
      setEnhancingId(null);
    }
  };

  if (loading) return <LoadingSpinner />;
  if (error) return <ErrorMessage message={error} />;

  return (
    <div className="note-list">
      <div className="note-list-header">
        <h3>{t('notes.title')}</h3>
        <button onClick={onAdd} className="btn btn-icon btn-add" title={t('notes.add')}>
          +
        </button>
      </div>

      {enhanceError && <ErrorMessage message={enhanceError} />}

      {notes.length === 0 ? (
        <div className="note-list-empty">
          <p>{t('notes.empty')}</p>
        </div>
      ) : (
        <div className="notes">
          {notes.map((note) => (
            <div key={note.id} className="note-card">
              <div className="note-header">
                <span className="note-number">#{note.note_number}</span>
                <div className="note-actions">
                  <button
                    onClick={() => handleEnhance(note.id)}
                    className="btn btn-icon btn-ai"
                    title={t('notes.enhance')}
                    disabled={enhancingId === note.id}
                  >
                    {enhancingId === note.id ? '⏳' : '✨'}
                  </button>
                  {previousContent?.noteId === note.id && (
                    <button
                      onClick={handleUndoEnhance}
                      className="btn btn-icon btn-undo"
                      title={t('notes.undoEnhance')}
                      disabled={enhancingId !== null}
                    >
                      ↶
                    </button>
                  )}
                  <button
                    onClick={() => onEdit(note.id)}
                    className="btn btn-icon btn-edit"
                    title={t('notes.edit')}
                  >
                    ✏
                  </button>
                  <button
                    onClick={() => confirmDelete(note.id, note.note_number)}
                    className="btn btn-icon btn-delete"
                    title={t('notes.delete')}
                  >
                    🗑
                  </button>
                </div>
              </div>
              <div className="note-content">
                {note.content}
              </div>
              <div className="note-footer">
                <small>{t('notes.updated', { date: new Date(note.updated_at).toLocaleString() })}</small>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
