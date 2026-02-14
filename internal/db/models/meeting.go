package models

import "time"

// Meeting represents a meeting record
type Meeting struct {
	ID           int       `json:"id"`
	CreatedBy    string    `json:"created_by"`
	Subject      string    `json:"subject"`
	MeetingDate  string    `json:"meeting_date"`  // YYYY-MM-DD
	StartTime    string    `json:"start_time"`    // HH:MM
	EndTime      *string   `json:"end_time"`      // optional
	Participants *string   `json:"participants"`  // optional
	Summary      *string   `json:"summary"`       // optional
	Keywords     *string   `json:"keywords"`      // optional
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
