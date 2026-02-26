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

func TestNoteRepository_SwapNoteOrder_Success(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

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

	noteRepo := repositories.NewNoteRepository(database.DB)
	note1 := &models.Note{MeetingID: meeting.ID, Content: "First"}
	note2 := &models.Note{MeetingID: meeting.ID, Content: "Second"}
	if err := noteRepo.Create(note1); err != nil {
		t.Fatalf("create note1: %v", err)
	}
	if err := noteRepo.Create(note2); err != nil {
		t.Fatalf("create note2: %v", err)
	}

	if err := noteRepo.SwapNoteOrder(note1.ID, note2.ID); err != nil {
		t.Fatalf("swap failed: %v", err)
	}

	// note1 should now have number 2, note2 should have number 1
	updated1, _ := noteRepo.GetByID(note1.ID)
	updated2, _ := noteRepo.GetByID(note2.ID)

	if updated1.NoteNumber != 2 {
		t.Errorf("expected note1 number=2, got %d", updated1.NoteNumber)
	}
	if updated2.NoteNumber != 1 {
		t.Errorf("expected note2 number=1, got %d", updated2.NoteNumber)
	}
}

func TestNoteRepository_SwapNoteOrder_NonAdjacentNumbers(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

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

	noteRepo := repositories.NewNoteRepository(database.DB)
	note1 := &models.Note{MeetingID: meeting.ID, Content: "First"}
	note2 := &models.Note{MeetingID: meeting.ID, Content: "Second"}
	note3 := &models.Note{MeetingID: meeting.ID, Content: "Third"}
	if err := noteRepo.Create(note1); err != nil {
		t.Fatalf("create note1: %v", err)
	}
	if err := noteRepo.Create(note2); err != nil {
		t.Fatalf("create note2: %v", err)
	}
	if err := noteRepo.Create(note3); err != nil {
		t.Fatalf("create note3: %v", err)
	}

	// Swap first (1) and last (3), skipping middle
	if err := noteRepo.SwapNoteOrder(note1.ID, note3.ID); err != nil {
		t.Fatalf("swap failed: %v", err)
	}

	updated1, _ := noteRepo.GetByID(note1.ID)
	updated3, _ := noteRepo.GetByID(note3.ID)

	if updated1.NoteNumber != 3 {
		t.Errorf("expected note1 number=3, got %d", updated1.NoteNumber)
	}
	if updated3.NoteNumber != 1 {
		t.Errorf("expected note3 number=1, got %d", updated3.NoteNumber)
	}
}

func TestNoteRepository_SwapNoteOrder_DifferentMeetings(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	meetingRepo := repositories.NewMeetingRepository(database.DB)
	meeting1 := &models.Meeting{
		CreatedBy:   "test@example.com",
		Subject:     "Meeting 1",
		MeetingDate: "2026-02-14",
		StartTime:   "10:00",
	}
	meeting2 := &models.Meeting{
		CreatedBy:   "test@example.com",
		Subject:     "Meeting 2",
		MeetingDate: "2026-02-14",
		StartTime:   "11:00",
	}
	if err := meetingRepo.Create(meeting1); err != nil {
		t.Fatalf("create meeting1: %v", err)
	}
	if err := meetingRepo.Create(meeting2); err != nil {
		t.Fatalf("create meeting2: %v", err)
	}

	noteRepo := repositories.NewNoteRepository(database.DB)
	note1 := &models.Note{MeetingID: meeting1.ID, Content: "Note in meeting 1"}
	note2 := &models.Note{MeetingID: meeting2.ID, Content: "Note in meeting 2"}
	if err := noteRepo.Create(note1); err != nil {
		t.Fatalf("create note1: %v", err)
	}
	if err := noteRepo.Create(note2); err != nil {
		t.Fatalf("create note2: %v", err)
	}

	err := noteRepo.SwapNoteOrder(note1.ID, note2.ID)
	if err == nil {
		t.Error("expected error swapping notes from different meetings, got nil")
	}
}

func TestNoteRepository_SwapNoteOrder_NoteNotFound(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	meetingRepo := repositories.NewMeetingRepository(database.DB)
	meeting := &models.Meeting{
		CreatedBy:   "test@example.com",
		Subject:     "Test Meeting",
		MeetingDate: "2026-02-14",
		StartTime:   "10:00",
	}
	if err := meetingRepo.Create(meeting); err != nil {
		t.Fatalf("create meeting: %v", err)
	}

	noteRepo := repositories.NewNoteRepository(database.DB)
	note := &models.Note{MeetingID: meeting.ID, Content: "Real note"}
	if err := noteRepo.Create(note); err != nil {
		t.Fatalf("create note: %v", err)
	}

	err := noteRepo.SwapNoteOrder(note.ID, 9999)
	if err == nil {
		t.Error("expected error for non-existent note, got nil")
	}
}

func TestNoteRepository_SwapNoteOrder_SameNote(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	meetingRepo := repositories.NewMeetingRepository(database.DB)
	meeting := &models.Meeting{
		CreatedBy:   "test@example.com",
		Subject:     "Test Meeting",
		MeetingDate: "2026-02-14",
		StartTime:   "10:00",
	}
	if err := meetingRepo.Create(meeting); err != nil {
		t.Fatalf("create meeting: %v", err)
	}

	noteRepo := repositories.NewNoteRepository(database.DB)
	note := &models.Note{MeetingID: meeting.ID, Content: "Test note"}
	if err := noteRepo.Create(note); err != nil {
		t.Fatalf("create note: %v", err)
	}

	err := noteRepo.SwapNoteOrder(note.ID, note.ID)
	if err == nil {
		t.Error("expected error swapping note with itself, got nil")
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
