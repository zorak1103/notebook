import { useTranslation } from 'react-i18next';
import { useNotes } from '../hooks/useNotes';
import { LoadingSpinner } from './LoadingSpinner';
import { ErrorMessage } from './ErrorMessage';
import './NoteList.css';

interface NoteListProps {
  meetingId: number;
  onEdit: (noteId: number) => void;
  onAdd: () => void;
  onEnhance?: (noteId: number) => void;
}

export function NoteList({ meetingId, onEdit, onAdd, onEnhance }: NoteListProps) {
  const { t } = useTranslation();
  const { notes, loading, error, handleDelete } = useNotes(meetingId);

  const confirmDelete = async (id: number, noteNumber: number) => {
    if (window.confirm(t('notes.confirmDelete', { number: noteNumber }))) {
      try {
        await handleDelete(id);
      } catch (err) {
        alert(err instanceof Error ? err.message : t('notes.deleteFailed'));
      }
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
                  {onEnhance && (
                    <button
                      onClick={() => onEnhance(note.id)}
                      className="btn btn-icon btn-ai"
                      title={t('notes.enhance')}
                    >
                      ✨
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
