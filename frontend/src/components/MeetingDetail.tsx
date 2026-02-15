import { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { fetchMeeting, summarizeMeeting, updateMeeting } from '../api/client';
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
  const [summarizing, setSummarizing] = useState(false);
  const [summaryError, setSummaryError] = useState<string | null>(null);
  const [previousSummary, setPreviousSummary] = useState<string | null>(null);

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

  const handleSummarize = async () => {
    if (!meeting) return;

    try {
      setSummarizing(true);
      setSummaryError(null);
      setPreviousSummary(meeting.summary);

      const updatedMeeting = await summarizeMeeting(meetingId);
      setMeeting(updatedMeeting);
    } catch (err) {
      setSummaryError(err instanceof Error ? err.message : t('meetingDetail.summarizeError'));
      setPreviousSummary(null);
    } finally {
      setSummarizing(false);
    }
  };

  const handleUndoSummary = async () => {
    if (!meeting || previousSummary === null) return;

    try {
      setSummarizing(true);
      setSummaryError(null);

      const updatedMeeting = await updateMeeting(meetingId, {
        subject: meeting.subject,
        meeting_date: meeting.meeting_date,
        start_time: meeting.start_time,
        end_time: meeting.end_time,
        participants: meeting.participants,
        summary: previousSummary,
        keywords: meeting.keywords,
      });

      setMeeting(updatedMeeting);
      setPreviousSummary(null);
    } catch (err) {
      setSummaryError(err instanceof Error ? err.message : t('meetingDetail.summarizeError'));
    } finally {
      setSummarizing(false);
    }
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
        <div className="detail-actions">
          <button
            onClick={handleSummarize}
            className="btn btn-icon btn-ai"
            title={t('meetingDetail.summarize')}
            disabled={summarizing}
          >
            {summarizing ? '⏳' : '✨'}
          </button>
          {previousSummary !== null && (
            <button
              onClick={handleUndoSummary}
              className="btn btn-icon btn-undo"
              title={t('meetingDetail.undoSummary')}
              disabled={summarizing}
            >
              ↶
            </button>
          )}
          <button onClick={onEdit} className="btn btn-icon btn-edit" title={t('meetingDetail.editMeeting')}>
            ✏
          </button>
        </div>
      </div>

      {summaryError && <ErrorMessage message={summaryError} />}

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
            <div className="meeting-summary">
              <span className="metadata-label">{t('meetingDetail.summary')}:</span>
              <p className="summary-text">{meeting.summary}</p>
            </div>
          )}
          {meeting.keywords && (
            <div className="meeting-keywords">
              <span className="metadata-label">{t('meetingDetail.keywords')}:</span>
              <div className="keywords-list">
                {meeting.keywords.split(',').map((keyword, idx) => (
                  <span key={idx} className="keyword-pill">
                    {keyword.trim()}
                  </span>
                ))}
              </div>
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
