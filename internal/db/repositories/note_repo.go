package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zorak1103/notebook/internal/db/models"
)

// NoteRepository handles note CRUD operations
type NoteRepository struct {
	db *sql.DB
}

// NewNoteRepository creates a new note repository
func NewNoteRepository(db *sql.DB) *NoteRepository {
	return &NoteRepository{db: db}
}

// Create creates a new note with automatic number assignment
func (r *NoteRepository) Create(n *models.Note) error {
	ctx := context.Background()

	// Get next number for this meeting
	var maxNumber int
	err := r.db.QueryRowContext(ctx, `
		SELECT COALESCE(MAX(note_number), 0) FROM notes WHERE meeting_id = ?
	`, n.MeetingID).Scan(&maxNumber)

	if err != nil {
		return fmt.Errorf("get max note number: %w", err)
	}

	n.NoteNumber = maxNumber + 1

	result, err := r.db.ExecContext(ctx, `
		INSERT INTO notes (meeting_id, note_number, content)
		VALUES (?, ?, ?)
	`, n.MeetingID, n.NoteNumber, n.Content)

	if err != nil {
		return fmt.Errorf("create note: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get last insert id: %w", err)
	}

	n.ID = int(id)
	return nil
}

// GetByID retrieves a note by ID
func (r *NoteRepository) GetByID(id int) (*models.Note, error) {
	ctx := context.Background()
	n := &models.Note{}
	err := r.db.QueryRowContext(ctx, `
		SELECT id, meeting_id, note_number, content, created_at, updated_at
		FROM notes WHERE id = ?
	`, id).Scan(&n.ID, &n.MeetingID, &n.NoteNumber, &n.Content, &n.CreatedAt, &n.UpdatedAt)

	if err == sql.ErrNoRows {
		//nolint:nilnil // Intentional: not found is not an error
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get note: %w", err)
	}

	return n, nil
}

// ListByMeeting lists all notes for a meeting
func (r *NoteRepository) ListByMeeting(meetingID int) ([]*models.Note, error) {
	ctx := context.Background()
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, meeting_id, note_number, content, created_at, updated_at
		FROM notes
		WHERE meeting_id = ?
		ORDER BY note_number ASC
	`, meetingID)

	if err != nil {
		return nil, fmt.Errorf("list notes: %w", err)
	}
	defer rows.Close()

	var notes []*models.Note
	for rows.Next() {
		n := &models.Note{}
		err := rows.Scan(&n.ID, &n.MeetingID, &n.NoteNumber, &n.Content, &n.CreatedAt, &n.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan note: %w", err)
		}
		notes = append(notes, n)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rows: %w", err)
	}

	return notes, nil
}

// Update updates an existing note
func (r *NoteRepository) Update(n *models.Note) error {
	ctx := context.Background()
	result, err := r.db.ExecContext(ctx, `
		UPDATE notes SET content = ? WHERE id = ?
	`, n.Content, n.ID)

	if err != nil {
		return fmt.Errorf("update note: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("note not found")
	}

	return nil
}

// Delete deletes a note
func (r *NoteRepository) Delete(id int) error {
	ctx := context.Background()
	result, err := r.db.ExecContext(ctx, "DELETE FROM notes WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete note: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("note not found")
	}

	return nil
}
