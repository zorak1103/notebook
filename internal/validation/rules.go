// Package validation provides shared validation constants for the application.
// These constants are used both in Go backend validation and code-generated
// TypeScript frontend validation.
package validation

const (
	// MaxSubjectLength is the maximum length for meeting subject field.
	MaxSubjectLength = 255
	// MaxParticipantsLength is the maximum length for meeting participants field.
	MaxParticipantsLength = 1000
	// MaxSummaryLength is the maximum length for meeting summary field.
	MaxSummaryLength = 10000
	// MaxKeywordsLength is the maximum length for meeting keywords field.
	MaxKeywordsLength = 500

	// MaxNoteContentLength is the maximum length for note content field.
	MaxNoteContentLength = 50000

	// MaxConfigKeyLength is the maximum length for configuration key field.
	MaxConfigKeyLength = 100
	// MaxConfigValueLength is the maximum length for configuration value field.
	MaxConfigValueLength = 1000
)
