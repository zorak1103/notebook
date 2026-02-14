package repositories_test

import (
	"testing"

	"github.com/zorak1103/notebook/internal/db/models"
	"github.com/zorak1103/notebook/internal/db/repositories"
)

func TestCascadeDelete_MeetingDeletesNotes(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	meetingRepo := repositories.NewMeetingRepository(database.DB)
	noteRepo := repositories.NewNoteRepository(database.DB)

	// Create a meeting
	meeting := &models.Meeting{
		CreatedBy:   "test@example.com",
		Subject:     "Test Meeting",
		MeetingDate: "2026-02-14",
		StartTime:   "10:00",
	}
	if err := meetingRepo.Create(meeting); err != nil {
		t.Fatalf("failed to create meeting: %v", err)
	}

	// Create notes for the meeting
	note1 := &models.Note{
		MeetingID: meeting.ID,
		Content:   "First note",
	}
	if err := noteRepo.Create(note1); err != nil {
		t.Fatalf("failed to create note 1: %v", err)
	}

	note2 := &models.Note{
		MeetingID: meeting.ID,
		Content:   "Second note",
	}
	if err := noteRepo.Create(note2); err != nil {
		t.Fatalf("failed to create note 2: %v", err)
	}

	// Verify notes exist
	notes, err := noteRepo.ListByMeeting(meeting.ID)
	if err != nil {
		t.Fatalf("failed to list notes: %v", err)
	}
	if len(notes) != 2 {
		t.Errorf("expected 2 notes before delete, got %d", len(notes))
	}

	// Delete the meeting
	if err := meetingRepo.Delete(meeting.ID); err != nil {
		t.Fatalf("failed to delete meeting: %v", err)
	}

	// Verify notes were cascade deleted
	notesAfterDelete, err := noteRepo.ListByMeeting(meeting.ID)
	if err != nil {
		t.Fatalf("failed to list notes after delete: %v", err)
	}
	if len(notesAfterDelete) != 0 {
		t.Errorf("expected 0 notes after cascade delete, got %d", len(notesAfterDelete))
	}

	// Verify notes are actually gone (not just filtered)
	note1Retrieved, err := noteRepo.GetByID(note1.ID)
	if err != nil {
		t.Fatalf("failed to get note 1: %v", err)
	}
	if note1Retrieved != nil {
		t.Error("expected note 1 to be deleted, but it still exists")
	}

	note2Retrieved, err := noteRepo.GetByID(note2.ID)
	if err != nil {
		t.Fatalf("failed to get note 2: %v", err)
	}
	if note2Retrieved != nil {
		t.Error("expected note 2 to be deleted, but it still exists")
	}
}
