package repositories_test

import (
	"testing"

	"github.com/zorak1103/notebook/internal/db"
	"github.com/zorak1103/notebook/internal/db/models"
	"github.com/zorak1103/notebook/internal/db/repositories"
)

func setupTestDB(t *testing.T) *db.DB {
	t.Helper()

	database, err := db.Open(":memory:")
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	if err := database.Migrate(); err != nil {
		t.Fatalf("failed to migrate test database: %v", err)
	}

	return database
}

func TestMeetingRepository_Create(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	repo := repositories.NewMeetingRepository(database.DB)

	meeting := &models.Meeting{
		CreatedBy:   "test@example.com",
		Subject:     "Test Meeting",
		MeetingDate: "2026-02-14",
		StartTime:   "10:00",
	}

	err := repo.Create(meeting)
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	if meeting.ID == 0 {
		t.Error("expected ID to be set")
	}
}

func TestMeetingRepository_GetByID(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	repo := repositories.NewMeetingRepository(database.DB)

	// Create a meeting
	meeting := &models.Meeting{
		CreatedBy:   "test@example.com",
		Subject:     "Test Meeting",
		MeetingDate: "2026-02-14",
		StartTime:   "10:00",
	}
	if err := repo.Create(meeting); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	// Get the meeting
	retrieved, err := repo.GetByID(meeting.ID)
	if err != nil {
		t.Fatalf("getByID failed: %v", err)
	}

	if retrieved == nil {
		t.Fatal("expected meeting to be found")
	}

	if retrieved.Subject != meeting.Subject {
		t.Errorf("expected subject %q, got %q", meeting.Subject, retrieved.Subject)
	}
}

func TestMeetingRepository_GetByID_NotFound(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	repo := repositories.NewMeetingRepository(database.DB)

	retrieved, err := repo.GetByID(999)
	if err != nil {
		t.Fatalf("getByID failed: %v", err)
	}

	if retrieved != nil {
		t.Error("expected meeting to be nil for non-existent ID")
	}
}

func TestMeetingRepository_List(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	repo := repositories.NewMeetingRepository(database.DB)

	// Create test meetings
	meetings := []*models.Meeting{
		{
			CreatedBy:   "test@example.com",
			Subject:     "Meeting A",
			MeetingDate: "2026-02-14",
			StartTime:   "10:00",
		},
		{
			CreatedBy:   "test@example.com",
			Subject:     "Meeting B",
			MeetingDate: "2026-02-15",
			StartTime:   "14:00",
		},
	}

	for _, m := range meetings {
		if err := repo.Create(m); err != nil {
			t.Fatalf("create failed: %v", err)
		}
	}

	// List meetings
	list, err := repo.List("meeting_date", false)
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}

	if len(list) != 2 {
		t.Errorf("expected 2 meetings, got %d", len(list))
	}

	// Check descending order (most recent first)
	if list[0].MeetingDate < list[1].MeetingDate {
		t.Error("expected meetings in descending order")
	}
}

func TestMeetingRepository_Update(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	repo := repositories.NewMeetingRepository(database.DB)

	// Create
	meeting := &models.Meeting{
		CreatedBy:   "test@example.com",
		Subject:     "Original Subject",
		MeetingDate: "2026-02-14",
		StartTime:   "10:00",
	}
	if err := repo.Create(meeting); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	// Update
	meeting.Subject = "Updated Subject"
	if err := repo.Update(meeting); err != nil {
		t.Fatalf("update failed: %v", err)
	}

	// Verify
	retrieved, err := repo.GetByID(meeting.ID)
	if err != nil {
		t.Fatalf("getByID failed: %v", err)
	}

	if retrieved.Subject != "Updated Subject" {
		t.Errorf("expected subject %q, got %q", "Updated Subject", retrieved.Subject)
	}
}

func TestMeetingRepository_Delete(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	repo := repositories.NewMeetingRepository(database.DB)

	// Create
	meeting := &models.Meeting{
		CreatedBy:   "test@example.com",
		Subject:     "Test Meeting",
		MeetingDate: "2026-02-14",
		StartTime:   "10:00",
	}
	if err := repo.Create(meeting); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	// Delete
	if err := repo.Delete(meeting.ID); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	// Verify
	retrieved, err := repo.GetByID(meeting.ID)
	if err != nil {
		t.Fatalf("getByID failed: %v", err)
	}

	if retrieved != nil {
		t.Error("expected meeting to be deleted")
	}
}

func TestMeetingRepository_Search(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	repo := repositories.NewMeetingRepository(database.DB)

	summary := "Important discussion"
	// Create test meetings
	meetings := []*models.Meeting{
		{
			CreatedBy:   "test@example.com",
			Subject:     "Team Meeting",
			MeetingDate: "2026-02-14",
			StartTime:   "10:00",
			Summary:     &summary,
		},
		{
			CreatedBy:   "test@example.com",
			Subject:     "Project Review",
			MeetingDate: "2026-02-15",
			StartTime:   "14:00",
		},
	}

	for _, m := range meetings {
		if err := repo.Create(m); err != nil {
			t.Fatalf("create failed: %v", err)
		}
	}

	// Search by subject
	results, err := repo.Search("Team")
	if err != nil {
		t.Fatalf("search failed: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}

	if len(results) > 0 && results[0].Subject != "Team Meeting" {
		t.Errorf("expected subject %q, got %q", "Team Meeting", results[0].Subject)
	}

	// Search by summary
	results, err = repo.Search("discussion")
	if err != nil {
		t.Fatalf("search failed: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
}
