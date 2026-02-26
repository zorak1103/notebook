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

// SwapNoteOrder swaps the note_number values of two notes within the same meeting.
// Uses a transaction with sentinel value 0 to work around the UNIQUE(meeting_id, note_number) constraint.
func (r *NoteRepository) SwapNoteOrder(noteID1, noteID2 int) error {
	if noteID1 == noteID2 {
		return fmt.Errorf("cannot swap note with itself")
	}

	ctx := context.Background()

	note1, err := r.GetByID(noteID1)
	if err != nil {
		return fmt.Errorf("get note1: %w", err)
	}
	if note1 == nil {
		return fmt.Errorf("note not found: %d", noteID1)
	}

	note2, err := r.GetByID(noteID2)
	if err != nil {
		return fmt.Errorf("get note2: %w", err)
	}
	if note2 == nil {
		return fmt.Errorf("note not found: %d", noteID2)
	}

	if note1.MeetingID != note2.MeetingID {
		return fmt.Errorf("notes belong to different meetings")
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	// Step 1: Move note1 to sentinel value 0 to free its slot
	_, err = tx.ExecContext(ctx, "UPDATE notes SET note_number = 0 WHERE id = ?", noteID1)
	if err != nil {
		return fmt.Errorf("set temp value: %w", err)
	}

	// Step 2: Move note2 into note1's original position
	_, err = tx.ExecContext(ctx, "UPDATE notes SET note_number = ? WHERE id = ?", note1.NoteNumber, noteID2)
	if err != nil {
		return fmt.Errorf("set note2 number: %w", err)
	}

	// Step 3: Move note1 into note2's original position
	_, err = tx.ExecContext(ctx, "UPDATE notes SET note_number = ? WHERE id = ?", note2.NoteNumber, noteID1)
	if err != nil {
		return fmt.Errorf("set note1 number: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
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
