import { useTranslation } from 'react-i18next';
import { useMeetings } from '../hooks/useMeetings';
import { LoadingSpinner } from './LoadingSpinner';
import { ErrorMessage } from './ErrorMessage';
import './MeetingList.css';

interface MeetingListProps {
  onEdit: (id: number) => void;
  onViewDetail: (id: number) => void;
}

export function MeetingList({ onEdit, onViewDetail }: MeetingListProps) {
  const { t } = useTranslation();
  const { meetings, loading, error, sortColumn, sortOrder, handleSort, handleDelete } = useMeetings();

  const getSortIndicator = (column: string) => {
    if (sortColumn !== column) return '';
    return sortOrder === 'asc' ? ' ↑' : ' ↓';
  };

  const confirmDelete = async (id: number, subject: string) => {
    if (window.confirm(t('meetings.confirmDelete', { subject }))) {
      try {
        await handleDelete(id);
      } catch (err) {
        alert(err instanceof Error ? err.message : t('meetings.deleteFailed'));
      }
    }
  };

  if (loading) return <LoadingSpinner />;
  if (error) return <ErrorMessage message={error} />;

  if (meetings.length === 0) {
    return (
      <div className="meeting-list-empty">
        <p>{t('meetings.empty')}</p>
      </div>
    );
  }

  return (
    <div className="meeting-list">
      {/* Sort toolbar */}
      <div className="sort-toolbar">
        <button
          onClick={() => handleSort('subject')}
          className={`sort-pill ${sortColumn === 'subject' ? 'sort-pill--active' : ''}`}
        >
          {t('meetings.subject')}{getSortIndicator('subject')}
        </button>
        <button
          onClick={() => handleSort('meeting_date')}
          className={`sort-pill ${sortColumn === 'meeting_date' ? 'sort-pill--active' : ''}`}
        >
          {t('meetings.date')}{getSortIndicator('meeting_date')}
        </button>
        <button
          onClick={() => handleSort('start_time')}
          className={`sort-pill ${sortColumn === 'start_time' ? 'sort-pill--active' : ''}`}
        >
          {t('meetings.startTime')}{getSortIndicator('start_time')}
        </button>
        <button
          onClick={() => handleSort('end_time')}
          className={`sort-pill ${sortColumn === 'end_time' ? 'sort-pill--active' : ''}`}
        >
          {t('meetings.endTime')}{getSortIndicator('end_time')}
        </button>
      </div>

      {/* Card grid */}
      <div className="meeting-grid">
        {meetings.map((meeting, index) => (
          <div
            key={meeting.id}
            className="meeting-card"
            style={{ animationDelay: `${index * 40}ms` }}
          >
            {/* Card content - clickable */}
            <div
              className="meeting-card-content"
              onClick={() => onViewDetail(meeting.id)}
            >
              <h3 className="meeting-card-subject">{meeting.subject}</h3>
              <div className="meeting-card-meta">
                <span className="meeting-card-date">{meeting.meeting_date}</span>
                <span className="meeting-card-time">
                  {meeting.start_time}
                  {meeting.end_time && ` – ${meeting.end_time}`}
                </span>
              </div>
              {meeting.participants && (
                <div className="meeting-card-participants">
                  {meeting.participants}
                </div>
              )}
              {meeting.keywords && (
                <div className="meeting-card-keywords">
                  {meeting.keywords}
                </div>
              )}
            </div>

            {/* Action buttons - visible on hover */}
            <div className="meeting-card-actions">
              <button
                onClick={(e) => {
                  e.stopPropagation();
                  onEdit(meeting.id);
                }}
                className="btn-icon btn-edit"
                title={t('meetings.edit')}
              >
                ✏
              </button>
              <button
                onClick={(e) => {
                  e.stopPropagation();
                  confirmDelete(meeting.id, meeting.subject);
                }}
                className="btn-icon btn-delete"
                title={t('meetings.delete')}
              >
                🗑
              </button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
