import { useTranslation } from 'react-i18next';
import { useMeetings } from '../hooks/useMeetings';
import { LoadingSpinner } from './LoadingSpinner';
import { ErrorMessage } from './ErrorMessage';
import './MeetingList.css';

interface MeetingListProps {
  onEdit: (id: number) => void;
}

export function MeetingList({ onEdit }: MeetingListProps) {
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
      <table className="meeting-table">
        <thead>
          <tr>
            <th onClick={() => handleSort('subject')} className="sortable">
              {t('meetings.subject')}{getSortIndicator('subject')}
            </th>
            <th onClick={() => handleSort('meeting_date')} className="sortable">
              {t('meetings.date')}{getSortIndicator('meeting_date')}
            </th>
            <th onClick={() => handleSort('start_time')} className="sortable">
              {t('meetings.startTime')}{getSortIndicator('start_time')}
            </th>
            <th onClick={() => handleSort('end_time')} className="sortable">
              {t('meetings.endTime')}{getSortIndicator('end_time')}
            </th>
            <th>{t('meetings.participants')}</th>
            <th>{t('meetings.actions')}</th>
          </tr>
        </thead>
        <tbody>
          {meetings.map((meeting) => (
            <tr key={meeting.id}>
              <td>{meeting.subject}</td>
              <td>{meeting.meeting_date}</td>
              <td>{meeting.start_time}</td>
              <td>{meeting.end_time || '-'}</td>
              <td>{meeting.participants || '-'}</td>
              <td className="actions">
                <button
                  onClick={() => onEdit(meeting.id)}
                  className="btn btn-edit"
                >
                  {t('meetings.edit')}
                </button>
                <button
                  onClick={() => confirmDelete(meeting.id, meeting.subject)}
                  className="btn btn-delete"
                >
                  {t('meetings.delete')}
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
