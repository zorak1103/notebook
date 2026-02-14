package models

import "time"

// Note represents a note within a meeting
type Note struct {
	ID         int       `json:"id"`
	MeetingID  int       `json:"meeting_id"`
	NoteNumber int       `json:"note_number"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
