import { useState, useEffect, FormEvent, useRef } from 'react';
import { useTranslation } from 'react-i18next';
import { fetchNote, createNote, updateNote, enhanceNote } from '../api/client';
import type { CreateNoteRequest, UpdateNoteRequest } from '../api/types';
import { LoadingSpinner } from './LoadingSpinner';
import { ErrorMessage } from './ErrorMessage';
import { MaxNoteContentLength } from '../generated/validationRules';
import './NoteForm.css';

interface NoteFormProps {
  meetingId: number;
  noteId?: number;
  onSuccess: () => void;
  onCancel: () => void;
}

export function NoteForm({ meetingId, noteId, onSuccess, onCancel }: NoteFormProps) {
  const { t } = useTranslation();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [content, setContent] = useState('');
  const contentTextareaRef = useRef<HTMLTextAreaElement>(null);
  const [enhancing, setEnhancing] = useState(false);
  const [enhanceError, setEnhanceError] = useState<string | null>(null);
  const [previousContent, setPreviousContent] = useState<string | null>(null);

  // Load note data if editing
  useEffect(() => {
    if (noteId) {
      setLoading(true);
      fetchNote(noteId)
        .then((note) => {
          setContent(note.content);
        })
        .catch((err) => {
          setError(err instanceof Error ? err.message : 'Failed to load note');
        })
        .finally(() => {
          setLoading(false);
        });
    } else {
      // Focus textarea when creating new note
      contentTextareaRef.current?.focus();
    }
  }, [noteId]);

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    try {
      if (noteId) {
        const updateData: UpdateNoteRequest = { content };
        await updateNote(noteId, updateData);
      } else {
        const createData: CreateNoteRequest = {
          meeting_id: meetingId,
          content,
        };
        await createNote(createData);
      }
      onSuccess();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to save note');
    } finally {
      setLoading(false);
    }
  };

  const handleEnhance = async () => {
    if (!noteId) return; // Only available in edit mode

    try {
      setEnhancing(true);
      setEnhanceError(null);
      setPreviousContent(content);

      const result = await enhanceNote(noteId, content);
      setContent(result.content);
    } catch (err) {
      setEnhanceError(err instanceof Error ? err.message : t('notes.enhanceError'));
      setPreviousContent(null);
    } finally {
      setEnhancing(false);
    }
  };

  const handleUndoEnhance = () => {
    if (previousContent !== null) {
      setContent(previousContent);
      setPreviousContent(null);
    }
  };

  if (loading && noteId) return <LoadingSpinner />;

  return (
    <div className="note-form">
      <div className="note-form-header">
        <h2 className="page-heading">{noteId ? t('noteForm.editTitle') : t('noteForm.createTitle')}</h2>
        {noteId && (
          <div className="note-form-actions">
            <button
              type="button"
              onClick={handleEnhance}
              className="btn btn-icon btn-ai"
              title={t('notes.enhance')}
              disabled={enhancing || !content.trim()}
            >
              {enhancing ? '⏳' : '✨'}
            </button>
            {previousContent !== null && (
              <button
                type="button"
                onClick={handleUndoEnhance}
                className="btn btn-icon btn-undo"
                title={t('notes.undoEnhance')}
                disabled={enhancing}
              >
                ↶
              </button>
            )}
          </div>
        )}
      </div>

      {error && <ErrorMessage message={error} />}
      {enhanceError && <ErrorMessage message={enhanceError} />}

      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label htmlFor="content">
            {t('noteForm.content')} <span className="required">*</span>
          </label>
          <textarea
            ref={contentTextareaRef}
            id="content"
            rows={10}
            value={content}
            onChange={(e) => setContent(e.target.value)}
            placeholder={t('noteForm.contentPlaceholder')}
            maxLength={MaxNoteContentLength}
            required
          />
          <small className="char-count">
            {t('noteForm.charCount', { current: content.length, max: MaxNoteContentLength })}
          </small>
        </div>

        <div className="form-actions">
          <button
            type="button"
            onClick={onCancel}
            className="btn btn-icon btn-cancel"
            title={t('noteForm.cancel')}
          >
            ×
          </button>
          <button
            type="submit"
            className="btn btn-icon btn-submit"
            disabled={loading}
            title={loading ? t('noteForm.saving') : t('noteForm.save')}
          >
            {loading ? '⏳' : '✓'}
          </button>
        </div>
      </form>
    </div>
  );
}
