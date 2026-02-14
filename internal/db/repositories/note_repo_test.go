package repositories_test

import (
	"testing"

	"github.com/zorak1103/notebook/internal/db/models"
	"github.com/zorak1103/notebook/internal/db/repositories"
)

func TestNoteRepository_Create(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	// Create a meeting first
	meetingRepo := repositories.NewMeetingRepository(database.DB)
	meeting := &models.Meeting{
		CreatedBy:   "test@example.com",
		Subject:     "Test Meeting",
		MeetingDate: "2026-02-14",
		StartTime:   "10:00",
	}
	if err := meetingRepo.Create(meeting); err != nil {
		t.Fatalf("failed to create meeting: %v", err)
	}

	// Create notes
	noteRepo := repositories.NewNoteRepository(database.DB)
	note := &models.Note{
		MeetingID: meeting.ID,
		Content:   "Test note content",
	}

	err := noteRepo.Create(note)
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	if note.ID == 0 {
		t.Error("expected ID to be set")
	}

	if note.NoteNumber != 1 {
		t.Errorf("expected note_number to be 1, got %d", note.NoteNumber)
	}
}

func TestNoteRepository_Create_AutoIncrement(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	// Create a meeting first
	meetingRepo := repositories.NewMeetingRepository(database.DB)
	meeting := &models.Meeting{
		CreatedBy:   "test@example.com",
		Subject:     "Test Meeting",
		MeetingDate: "2026-02-14",
		StartTime:   "10:00",
	}
	if err := meetingRepo.Create(meeting); err != nil {
		t.Fatalf("failed to create meeting: %v", err)
	}

	// Create multiple notes
	noteRepo := repositories.NewNoteRepository(database.DB)
	for i := 1; i <= 3; i++ {
		note := &models.Note{
			MeetingID: meeting.ID,
			Content:   "Note content",
		}
		if err := noteRepo.Create(note); err != nil {
			t.Fatalf("create note %d failed: %v", i, err)
		}
		if note.NoteNumber != i {
			t.Errorf("expected note_number %d, got %d", i, note.NoteNumber)
		}
	}
}

func TestNoteRepository_GetByID(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	// Create meeting and note
	meetingRepo := repositories.NewMeetingRepository(database.DB)
	meeting := &models.Meeting{
		CreatedBy:   "test@example.com",
		Subject:     "Test Meeting",
		MeetingDate: "2026-02-14",
		StartTime:   "10:00",
	}
	meetingRepo.Create(meeting)

	noteRepo := repositories.NewNoteRepository(database.DB)
	note := &models.Note{
		MeetingID: meeting.ID,
		Content:   "Test note",
	}
	noteRepo.Create(note)

	// Get the note
	retrieved, err := noteRepo.GetByID(note.ID)
	if err != nil {
		t.Fatalf("getByID failed: %v", err)
	}

	if retrieved == nil {
		t.Fatal("expected note to be found")
	}

	if retrieved.Content != note.Content {
		t.Errorf("expected content %q, got %q", note.Content, retrieved.Content)
	}
}

func TestNoteRepository_ListByMeeting(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	// Create meeting
	meetingRepo := repositories.NewMeetingRepository(database.DB)
	meeting := &models.Meeting{
		CreatedBy:   "test@example.com",
		Subject:     "Test Meeting",
		MeetingDate: "2026-02-14",
		StartTime:   "10:00",
	}
	meetingRepo.Create(meeting)

	// Create notes
	noteRepo := repositories.NewNoteRepository(database.DB)
	for i := 1; i <= 3; i++ {
		note := &models.Note{
			MeetingID: meeting.ID,
			Content:   "Note content",
		}
		noteRepo.Create(note)
	}

	// List notes
	notes, err := noteRepo.ListByMeeting(meeting.ID)
	if err != nil {
		t.Fatalf("listByMeeting failed: %v", err)
	}

	if len(notes) != 3 {
		t.Errorf("expected 3 notes, got %d", len(notes))
	}

	// Verify ordering (ascending note_number)
	for i, note := range notes {
		expectedNum := i + 1
		if note.NoteNumber != expectedNum {
			t.Errorf("expected note_number %d at index %d, got %d", expectedNum, i, note.NoteNumber)
		}
	}
}

func TestNoteRepository_Update(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	// Create meeting and note
	meetingRepo := repositories.NewMeetingRepository(database.DB)
	meeting := &models.Meeting{
		CreatedBy:   "test@example.com",
		Subject:     "Test Meeting",
		MeetingDate: "2026-02-14",
		StartTime:   "10:00",
	}
	meetingRepo.Create(meeting)

	noteRepo := repositories.NewNoteRepository(database.DB)
	note := &models.Note{
		MeetingID: meeting.ID,
		Content:   "Original content",
	}
	noteRepo.Create(note)

	// Update
	note.Content = "Updated content"
	if err := noteRepo.Update(note); err != nil {
		t.Fatalf("update failed: %v", err)
	}

	// Verify
	retrieved, err := noteRepo.GetByID(note.ID)
	if err != nil {
		t.Fatalf("getByID failed: %v", err)
	}

	if retrieved.Content != "Updated content" {
		t.Errorf("expected content %q, got %q", "Updated content", retrieved.Content)
	}
}

func TestNoteRepository_Delete(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	// Create meeting and note
	meetingRepo := repositories.NewMeetingRepository(database.DB)
	meeting := &models.Meeting{
		CreatedBy:   "test@example.com",
		Subject:     "Test Meeting",
		MeetingDate: "2026-02-14",
		StartTime:   "10:00",
	}
	meetingRepo.Create(meeting)

	noteRepo := repositories.NewNoteRepository(database.DB)
	note := &models.Note{
		MeetingID: meeting.ID,
		Content:   "Test note",
	}
	noteRepo.Create(note)

	// Delete
	if err := noteRepo.Delete(note.ID); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	// Verify
	retrieved, err := noteRepo.GetByID(note.ID)
	if err != nil {
		t.Fatalf("getByID failed: %v", err)
	}

	if retrieved != nil {
		t.Error("expected note to be deleted")
	}
}

func TestNoteRepository_DeleteCascade(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	// Create meeting
	meetingRepo := repositories.NewMeetingRepository(database.DB)
	meeting := &models.Meeting{
		CreatedBy:   "test@example.com",
		Subject:     "Test Meeting",
		MeetingDate: "2026-02-14",
		StartTime:   "10:00",
	}
	meetingRepo.Create(meeting)

	// Create notes
	noteRepo := repositories.NewNoteRepository(database.DB)
	note := &models.Note{
		MeetingID: meeting.ID,
		Content:   "Test note",
	}
	noteRepo.Create(note)

	// Delete meeting (should cascade to notes)
	if err := meetingRepo.Delete(meeting.ID); err != nil {
		t.Fatalf("delete meeting failed: %v", err)
	}

	// Verify note is deleted
	retrieved, err := noteRepo.GetByID(note.ID)
	if err != nil {
		t.Fatalf("getByID failed: %v", err)
	}

	if retrieved != nil {
		t.Error("expected note to be deleted via CASCADE")
	}
}
