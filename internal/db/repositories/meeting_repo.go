package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zorak1103/notebook/internal/db/models"
)

// MeetingRepository handles meeting CRUD operations
type MeetingRepository struct {
	db *sql.DB
}

// NewMeetingRepository creates a new meeting repository
func NewMeetingRepository(db *sql.DB) *MeetingRepository {
	return &MeetingRepository{db: db}
}

// Create creates a new meeting
func (r *MeetingRepository) Create(m *models.Meeting) error {
	ctx := context.Background()
	result, err := r.db.ExecContext(ctx, `
		INSERT INTO meetings (created_by, subject, meeting_date, start_time, end_time, participants, summary, keywords)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, m.CreatedBy, m.Subject, m.MeetingDate, m.StartTime, m.EndTime, m.Participants, m.Summary, m.Keywords)

	if err != nil {
		return fmt.Errorf("create meeting: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get last insert id: %w", err)
	}

	m.ID = int(id)
	return nil
}

// GetByID retrieves a meeting by ID
func (r *MeetingRepository) GetByID(id int) (*models.Meeting, error) {
	ctx := context.Background()
	m := &models.Meeting{}
	err := r.db.QueryRowContext(ctx, `
		SELECT id, created_by, subject, meeting_date, start_time, end_time, participants, summary, keywords, created_at, updated_at
		FROM meetings WHERE id = ?
	`, id).Scan(&m.ID, &m.CreatedBy, &m.Subject, &m.MeetingDate, &m.StartTime, &m.EndTime, &m.Participants, &m.Summary, &m.Keywords, &m.CreatedAt, &m.UpdatedAt)

	if err == sql.ErrNoRows {
		//nolint:nilnil // Intentional: not found is not an error
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get meeting: %w", err)
	}

	return m, nil
}

// List lists all meetings with optional sorting
func (r *MeetingRepository) List(orderBy string, ascending bool) ([]*models.Meeting, error) {
	// Whitelist for ORDER BY (SQL injection protection)
	validColumns := map[string]bool{
		"meeting_date": true,
		"start_time":   true,
		"end_time":     true,
		"subject":      true,
		"keywords":     true,
	}

	if !validColumns[orderBy] {
		orderBy = "meeting_date"
	}

	direction := "DESC"
	if ascending {
		direction = "ASC"
	}

	//nolint:gosec // SQL injection protected by whitelist validation above
	query := fmt.Sprintf(`
		SELECT id, created_by, subject, meeting_date, start_time, end_time, participants, summary, keywords, created_at, updated_at
		FROM meetings
		ORDER BY %s COLLATE NOCASE %s
	`, orderBy, direction)

	ctx := context.Background()
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list meetings: %w", err)
	}
	defer rows.Close()

	var meetings []*models.Meeting
	for rows.Next() {
		m := &models.Meeting{}
		err := rows.Scan(&m.ID, &m.CreatedBy, &m.Subject, &m.MeetingDate, &m.StartTime, &m.EndTime, &m.Participants, &m.Summary, &m.Keywords, &m.CreatedAt, &m.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan meeting: %w", err)
		}
		meetings = append(meetings, m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rows: %w", err)
	}

	return meetings, nil
}

// Update updates an existing meeting
func (r *MeetingRepository) Update(m *models.Meeting) error {
	ctx := context.Background()
	result, err := r.db.ExecContext(ctx, `
		UPDATE meetings
		SET subject = ?, meeting_date = ?, start_time = ?, end_time = ?, participants = ?, summary = ?, keywords = ?
		WHERE id = ?
	`, m.Subject, m.MeetingDate, m.StartTime, m.EndTime, m.Participants, m.Summary, m.Keywords, m.ID)

	if err != nil {
		return fmt.Errorf("update meeting: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("meeting not found")
	}

	return nil
}

// Delete deletes a meeting (and all associated notes via CASCADE)
func (r *MeetingRepository) Delete(id int) error {
	ctx := context.Background()
	result, err := r.db.ExecContext(ctx, "DELETE FROM meetings WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete meeting: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("meeting not found")
	}

	return nil
}

// escapeLikePattern escapes special LIKE characters and wraps with wildcards
func escapeLikePattern(s string) string {
	// Escape backslash first to avoid double-escaping
	s = strings.ReplaceAll(s, `\`, `\\`)
	// Escape % and _ wildcards
	s = strings.ReplaceAll(s, "%", `\%`)
	s = strings.ReplaceAll(s, "_", `\_`)
	// Wrap with wildcards
	return "%" + s + "%"
}

// Search searches meetings by subject, summary, participants, and keywords
func (r *MeetingRepository) Search(query string) ([]*models.Meeting, error) {
	ctx := context.Background()
	pattern := escapeLikePattern(query)

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, created_by, subject, meeting_date, start_time, end_time, participants, summary, keywords, created_at, updated_at
		FROM meetings
		WHERE subject LIKE ? ESCAPE '\'
		   OR summary LIKE ? ESCAPE '\'
		   OR participants LIKE ? ESCAPE '\'
		   OR keywords LIKE ? ESCAPE '\'
		ORDER BY meeting_date DESC, start_time DESC
	`, pattern, pattern, pattern, pattern)

	if err != nil {
		return nil, fmt.Errorf("search meetings: %w", err)
	}
	defer rows.Close()

	var meetings []*models.Meeting
	for rows.Next() {
		m := &models.Meeting{}
		err := rows.Scan(&m.ID, &m.CreatedBy, &m.Subject, &m.MeetingDate, &m.StartTime, &m.EndTime, &m.Participants, &m.Summary, &m.Keywords, &m.CreatedAt, &m.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan meeting: %w", err)
		}
		meetings = append(meetings, m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rows: %w", err)
	}

	return meetings, nil
}
