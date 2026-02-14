package validation

import "testing"

func TestValidationConstants(t *testing.T) {
	tests := []struct {
		name     string
		value    int
		minValue int
		maxValue int
	}{
		{"MaxSubjectLength", MaxSubjectLength, 1, 1000},
		{"MaxParticipantsLength", MaxParticipantsLength, 1, 10000},
		{"MaxSummaryLength", MaxSummaryLength, 1, 100000},
		{"MaxKeywordsLength", MaxKeywordsLength, 1, 10000},
		{"MaxNoteContentLength", MaxNoteContentLength, 1, 100000},
		{"MaxConfigKeyLength", MaxConfigKeyLength, 1, 500},
		{"MaxConfigValueLength", MaxConfigValueLength, 1, 10000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value < tt.minValue {
				t.Errorf("%s = %d, want >= %d", tt.name, tt.value, tt.minValue)
			}
			if tt.value > tt.maxValue {
				t.Errorf("%s = %d, want <= %d", tt.name, tt.value, tt.maxValue)
			}
		})
	}
}

func TestValidationConstantsArePositive(t *testing.T) {
	constants := map[string]int{
		"MaxSubjectLength":      MaxSubjectLength,
		"MaxParticipantsLength": MaxParticipantsLength,
		"MaxSummaryLength":      MaxSummaryLength,
		"MaxKeywordsLength":     MaxKeywordsLength,
		"MaxNoteContentLength":  MaxNoteContentLength,
		"MaxConfigKeyLength":    MaxConfigKeyLength,
		"MaxConfigValueLength":  MaxConfigValueLength,
	}

	for name, value := range constants {
		if value <= 0 {
			t.Errorf("%s must be positive, got %d", name, value)
		}
	}
}
