import { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { fetchMeeting } from '../api/client';
import type { Meeting } from '../api/types';
import { LoadingSpinner } from './LoadingSpinner';
import { ErrorMessage } from './ErrorMessage';
import { NoteList } from './NoteList';
import { NoteForm } from './NoteForm';
import './MeetingDetail.css';

interface MeetingDetailProps {
  meetingId: number;
  onBack: () => void;
  onEdit: () => void;
}

type NoteView = 'list' | 'create' | 'edit';

export function MeetingDetail({ meetingId, onBack, onEdit }: MeetingDetailProps) {
  const { t } = useTranslation();
  const [meeting, setMeeting] = useState<Meeting | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [noteView, setNoteView] = useState<NoteView>('list');
  const [editingNoteId, setEditingNoteId] = useState<number | undefined>(undefined);

  useEffect(() => {
    setLoading(true);
    fetchMeeting(meetingId)
      .then((data) => {
        setMeeting(data);
      })
      .catch((err) => {
        setError(err instanceof Error ? err.message : 'Failed to load meeting');
      })
      .finally(() => {
        setLoading(false);
      });
  }, [meetingId]);

  const handleNoteSuccess = () => {
    setNoteView('list');
    setEditingNoteId(undefined);
  };

  const handleNoteCancel = () => {
    setNoteView('list');
    setEditingNoteId(undefined);
  };

  const handleAddNote = () => {
    setNoteView('create');
    setEditingNoteId(undefined);
  };

  const handleEditNote = (noteId: number) => {
    setNoteView('edit');
    setEditingNoteId(noteId);
  };

  if (loading) return <LoadingSpinner />;
  if (error) return <ErrorMessage message={error} />;
  if (!meeting) return <ErrorMessage message={t('meetingDetail.notFound')} />;

  return (
    <div className="meeting-detail">
      <div className="detail-header">
        <button onClick={onBack} className="btn btn-back">
          ← {t('meetingDetail.back')}
        </button>
        <button onClick={onEdit} className="btn btn-edit">
          {t('meetingDetail.editMeeting')}
        </button>
      </div>

      <div className="meeting-info">
        <h1>{meeting.subject}</h1>
        <div className="meeting-metadata">
          <div className="metadata-row">
            <span className="metadata-label">{t('meetingDetail.date')}:</span>
            <span className="metadata-value">{meeting.meeting_date}</span>
          </div>
          <div className="metadata-row">
            <span className="metadata-label">{t('meetingDetail.time')}:</span>
            <span className="metadata-value">
              {meeting.start_time}
              {meeting.end_time && ` - ${meeting.end_time}`}
            </span>
          </div>
          {meeting.participants && (
            <div className="metadata-row">
              <span className="metadata-label">{t('meetingDetail.participants')}:</span>
              <span className="metadata-value">{meeting.participants}</span>
            </div>
          )}
          {meeting.summary && (
            <div className="metadata-row">
              <span className="metadata-label">{t('meetingDetail.summary')}:</span>
              <span className="metadata-value">{meeting.summary}</span>
            </div>
          )}
          {meeting.keywords && (
            <div className="metadata-row">
              <span className="metadata-label">{t('meetingDetail.keywords')}:</span>
              <span className="metadata-value">{meeting.keywords}</span>
            </div>
          )}
        </div>
      </div>

      <div className="notes-section">
        {noteView === 'list' && (
          <NoteList
            meetingId={meetingId}
            onEdit={handleEditNote}
            onAdd={handleAddNote}
          />
        )}
        {noteView === 'create' && (
          <NoteForm
            meetingId={meetingId}
            onSuccess={handleNoteSuccess}
            onCancel={handleNoteCancel}
          />
        )}
        {noteView === 'edit' && editingNoteId && (
          <NoteForm
            meetingId={meetingId}
            noteId={editingNoteId}
            onSuccess={handleNoteSuccess}
            onCancel={handleNoteCancel}
          />
        )}
      </div>
    </div>
  );
}
