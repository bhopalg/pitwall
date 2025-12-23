package utils

import (
	"testing"
)

func TestParseDate(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "Valid RFC3339 Date",
			input:   "2023-11-26T13:00:00Z",
			wantErr: false,
		},
		{
			name:    "Valid RFC3339 with Offset",
			input:   "2023-11-26T13:00:00+01:00",
			wantErr: false,
		},
		{
			name:    "Invalid Date Format",
			input:   "26-11-2023 13:00",
			wantErr: true,
		},
		{
			name:    "Empty String",
			input:   "",
			wantErr: true,
		},
		{
			name:    "Random Text",
			input:   "not-a-date",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDate(tt.input)

			// Check if we got an error when we expected one
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// If no error was expected, ensure the pointer isn't nil
			if !tt.wantErr && got == nil {
				t.Error("ParseDate() returned nil for a valid input")
			}
		})
	}
}
