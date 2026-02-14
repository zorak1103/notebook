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

func TestMeetingRepository_Search_EscapesSpecialChars(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	repo := repositories.NewMeetingRepository(database.DB)

	// Create meetings with literal % in subject
	meetings := []*models.Meeting{
		{
			CreatedBy:   "test@example.com",
			Subject:     "100% Coverage",
			MeetingDate: "2026-02-14",
			StartTime:   "10:00",
		},
		{
			CreatedBy:   "test@example.com",
			Subject:     "Design Review",
			MeetingDate: "2026-02-15",
			StartTime:   "14:00",
		},
	}

	for _, m := range meetings {
		if err := repo.Create(m); err != nil {
			t.Fatalf("create failed: %v", err)
		}
	}

	// Search for literal % should only return "100% Coverage"
	results, err := repo.Search("%")
	if err != nil {
		t.Fatalf("search failed: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("expected 1 result for %% search, got %d", len(results))
	}

	if len(results) > 0 && results[0].Subject != "100% Coverage" {
		t.Errorf("expected '100%% Coverage', got '%s'", results[0].Subject)
	}
}

func TestMeetingRepository_Search_ByParticipants(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	repo := repositories.NewMeetingRepository(database.DB)

	participants1 := "Alice, Bob, Charlie"
	participants2 := "David, Eve"
	// Create meetings with different participants
	meetings := []*models.Meeting{
		{
			CreatedBy:    "test@example.com",
			Subject:      "Meeting 1",
			MeetingDate:  "2026-02-14",
			StartTime:    "10:00",
			Participants: &participants1,
		},
		{
			CreatedBy:    "test@example.com",
			Subject:      "Meeting 2",
			MeetingDate:  "2026-02-15",
			StartTime:    "14:00",
			Participants: &participants2,
		},
	}

	for _, m := range meetings {
		if err := repo.Create(m); err != nil {
			t.Fatalf("create failed: %v", err)
		}
	}

	// Search by participant name
	results, err := repo.Search("Alice")
	if err != nil {
		t.Fatalf("search failed: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("expected 1 result for participant search, got %d", len(results))
	}

	if len(results) > 0 && results[0].Subject != "Meeting 1" {
		t.Errorf("expected 'Meeting 1', got %q", results[0].Subject)
	}
}

func TestMeetingRepository_Search_ByKeywords(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	repo := repositories.NewMeetingRepository(database.DB)

	keywords1 := "planning, quarterly, review"
	keywords2 := "retrospective, sprint"
	// Create meetings with different keywords
	meetings := []*models.Meeting{
		{
			CreatedBy:   "test@example.com",
			Subject:     "Q1 Planning",
			MeetingDate: "2026-02-14",
			StartTime:   "10:00",
			Keywords:    &keywords1,
		},
		{
			CreatedBy:   "test@example.com",
			Subject:     "Sprint Review",
			MeetingDate: "2026-02-15",
			StartTime:   "14:00",
			Keywords:    &keywords2,
		},
	}

	for _, m := range meetings {
		if err := repo.Create(m); err != nil {
			t.Fatalf("create failed: %v", err)
		}
	}

	// Search by keyword
	results, err := repo.Search("quarterly")
	if err != nil {
		t.Fatalf("search failed: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("expected 1 result for keyword search, got %d", len(results))
	}

	if len(results) > 0 && results[0].Subject != "Q1 Planning" {
		t.Errorf("expected 'Q1 Planning', got %q", results[0].Subject)
	}
}

func TestMeetingRepository_Search_EscapesUnderscore(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	repo := repositories.NewMeetingRepository(database.DB)

	// Create meetings with underscore patterns
	meetings := []*models.Meeting{
		{
			CreatedBy:   "test@example.com",
			Subject:     "A_B Test",
			MeetingDate: "2026-02-14",
			StartTime:   "10:00",
		},
		{
			CreatedBy:   "test@example.com",
			Subject:     "AXB Test",
			MeetingDate: "2026-02-15",
			StartTime:   "14:00",
		},
	}

	for _, m := range meetings {
		if err := repo.Create(m); err != nil {
			t.Fatalf("create failed: %v", err)
		}
	}

	// Search for literal _ should only return "A_B Test"
	results, err := repo.Search("_")
	if err != nil {
		t.Fatalf("search failed: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("expected 1 result for _ search, got %d", len(results))
	}

	if len(results) > 0 && results[0].Subject != "A_B Test" {
		t.Errorf("expected 'A_B Test', got '%s'", results[0].Subject)
	}
}
