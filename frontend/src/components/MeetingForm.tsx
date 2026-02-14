import { useState, useEffect, FormEvent, useRef } from 'react';
import { useTranslation } from 'react-i18next';
import { fetchMeeting, createMeeting, updateMeeting } from '../api/client';
import type { CreateMeetingRequest } from '../api/types';
import { LoadingSpinner } from './LoadingSpinner';
import { ErrorMessage } from './ErrorMessage';
import {
  MaxSubjectLength,
  MaxParticipantsLength,
  MaxSummaryLength,
  MaxKeywordsLength,
} from '../generated/validationRules';
import './MeetingForm.css';

interface MeetingFormProps {
  meetingId?: number;
  onSuccess: () => void;
  onCancel: () => void;
}

export function MeetingForm({ meetingId, onSuccess, onCancel }: MeetingFormProps) {
  const { t } = useTranslation();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const subjectInputRef = useRef<HTMLInputElement>(null);

  // Get current date and time for default values
  const now = new Date();
  const currentDate = now.toISOString().split('T')[0]; // YYYY-MM-DD
  const currentTime = now.toTimeString().slice(0, 5); // HH:MM

  const [formData, setFormData] = useState<CreateMeetingRequest>({
    subject: '',
    meeting_date: meetingId ? '' : currentDate,
    start_time: meetingId ? '' : currentTime,
    end_time: null,
    participants: null,
    summary: null,
    keywords: null,
  });

  // Load meeting data if editing
  useEffect(() => {
    if (meetingId) {
      setLoading(true);
      fetchMeeting(meetingId)
        .then((meeting) => {
          setFormData({
            subject: meeting.subject,
            meeting_date: meeting.meeting_date,
            start_time: meeting.start_time,
            end_time: meeting.end_time,
            participants: meeting.participants,
            summary: meeting.summary,
            keywords: meeting.keywords,
          });
        })
        .catch((err) => {
          setError(err instanceof Error ? err.message : 'Failed to load meeting');
        })
        .finally(() => {
          setLoading(false);
        });
    } else {
      // Focus subject field when creating new meeting
      subjectInputRef.current?.focus();
    }
  }, [meetingId]);

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    try {
      if (meetingId) {
        await updateMeeting(meetingId, formData);
      } else {
        await createMeeting(formData);
      }
      onSuccess();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to save meeting');
    } finally {
      setLoading(false);
    }
  };

  const handleChange = (field: keyof CreateMeetingRequest, value: string) => {
    setFormData((prev) => ({
      ...prev,
      [field]: value === '' ? null : value,
    }));
  };

  if (loading && meetingId) return <LoadingSpinner />;

  return (
    <div className="meeting-form">
      <h2>{meetingId ? t('meetingForm.editTitle') : t('meetingForm.createTitle')}</h2>

      {error && <ErrorMessage message={error} />}

      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label htmlFor="subject">
            {t('meetingForm.subject')} <span className="required">*</span>
          </label>
          <input
            ref={subjectInputRef}
            type="text"
            id="subject"
            value={formData.subject}
            onChange={(e) => handleChange('subject', e.target.value)}
            placeholder={t('meetingForm.subjectPlaceholder')}
            maxLength={MaxSubjectLength}
            required
          />
          <small className="char-count">
            {t('meetingForm.charCount', { current: formData.subject.length, max: MaxSubjectLength })}
          </small>
        </div>

        <div className="form-row">
          <div className="form-group">
            <label htmlFor="meeting_date">
              {t('meetingForm.date')} <span className="required">*</span>
            </label>
            <input
              type="date"
              id="meeting_date"
              value={formData.meeting_date}
              onChange={(e) => handleChange('meeting_date', e.target.value)}
              required
            />
          </div>

          <div className="form-group">
            <label htmlFor="start_time">
              {t('meetingForm.startTime')} <span className="required">*</span>
            </label>
            <input
              type="time"
              id="start_time"
              value={formData.start_time}
              onChange={(e) => handleChange('start_time', e.target.value)}
              required
            />
          </div>

          <div className="form-group">
            <label htmlFor="end_time">{t('meetingForm.endTime')}</label>
            <input
              type="time"
              id="end_time"
              value={formData.end_time || ''}
              onChange={(e) => handleChange('end_time', e.target.value)}
            />
          </div>
        </div>

        <div className="form-group">
          <label htmlFor="participants">{t('meetingForm.participants')}</label>
          <input
            type="text"
            id="participants"
            value={formData.participants || ''}
            onChange={(e) => handleChange('participants', e.target.value)}
            placeholder={t('meetingForm.participantsPlaceholder')}
            maxLength={MaxParticipantsLength}
          />
          <small className="char-count">
            {t('meetingForm.charCount', { current: (formData.participants || '').length, max: MaxParticipantsLength })}
          </small>
        </div>

        <div className="form-group">
          <label htmlFor="summary">{t('meetingForm.summary')}</label>
          <textarea
            id="summary"
            rows={4}
            value={formData.summary || ''}
            onChange={(e) => handleChange('summary', e.target.value)}
            placeholder={t('meetingForm.summaryPlaceholder')}
            maxLength={MaxSummaryLength}
          />
          <small className="char-count">
            {t('meetingForm.charCount', { current: (formData.summary || '').length, max: MaxSummaryLength })}
          </small>
        </div>

        <div className="form-group">
          <label htmlFor="keywords">{t('meetingForm.keywords')}</label>
          <input
            type="text"
            id="keywords"
            value={formData.keywords || ''}
            onChange={(e) => handleChange('keywords', e.target.value)}
            placeholder={t('meetingForm.keywordsPlaceholder')}
            maxLength={MaxKeywordsLength}
          />
          <small className="char-count">
            {t('meetingForm.charCount', { current: (formData.keywords || '').length, max: MaxKeywordsLength })}
          </small>
        </div>

        <div className="form-actions">
          <button type="button" onClick={onCancel} className="btn btn-cancel">
            {t('meetingForm.cancel')}
          </button>
          <button type="submit" className="btn btn-submit" disabled={loading}>
            {loading ? t('meetingForm.saving') : t('meetingForm.save')}
          </button>
        </div>
      </form>
    </div>
  );
}
