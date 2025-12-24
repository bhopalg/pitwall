package utils

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/bhopalg/pitwall/domain"
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

func TestPrintSessionStatus(t *testing.T) {
	// Baseline: Saturday 29 July 2023, 14:00 UTC
	start := time.Date(2023, 7, 29, 14, 0, 0, 0, time.UTC)
	end := time.Date(2023, 7, 29, 16, 0, 0, 0, time.UTC)

	tests := []struct {
		name           string
		session        *domain.Session
		now            time.Time
		expectedOutput []string
	}{
		{
			name: "Future Session",
			session: &domain.Session{
				DateStart: start,
				DateEnd:   end,
			},
			now: start.Add(-2 * time.Hour).Add(-14 * time.Minute),
			expectedOutput: []string{
				"Status: Future",
				"Starts at: Sat 29 Jul 2023, 15:00 (UK)",
				"Starts in: 2h 14m",
			},
		},
		{
			name: "Live Session",
			session: &domain.Session{
				DateStart: start,
				DateEnd:   end,
			},
			now: start.Add(1 * time.Hour).Add(42 * time.Minute),
			expectedOutput: []string{
				"Status: Live",
				"Ends at: 17:00 (UK)",
				"Ends in: 0h 18m",
			},
		},
		{
			name: "Finished Session",
			session: &domain.Session{
				DateStart: start,
				DateEnd:   end,
			},
			now: end.Add(1 * time.Hour).Add(42 * time.Minute),
			expectedOutput: []string{
				"Status: Finished",
				"Ended at: 17:00 (UK)",
				"Ended: 1h 42m ago",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			PrintSessionStatus(tt.session, tt.now)

			w.Close()
			os.Stdout = oldStdout

			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()

			for _, expected := range tt.expectedOutput {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected output to contain %q, but got:\n%s", expected, output)
				}
			}
		})
	}
}
