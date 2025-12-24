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
	start := time.Date(2023, 7, 29, 14, 0, 0, 0, time.UTC)
	end := time.Date(2023, 7, 29, 16, 0, 0, 0, time.UTC)

	tests := []struct {
		name           string
		session        *domain.Session
		now            time.Time
		expectedOutput []string
	}{
		{
			name: "Future Session (Long term)",
			session: &domain.Session{
				DateStart: start,
				DateEnd:   end,
			},
			now: start.Add(-52 * time.Hour).Add(-30 * time.Minute),
			expectedOutput: []string{
				"Status: Future",
				"Starts in: 2d 4h 30m",
			},
		},
		{
			name: "Future Session (Short term)",
			session: &domain.Session{
				DateStart: start,
				DateEnd:   end,
			},
			now: start.Add(-2 * time.Hour).Add(-14 * time.Minute),
			expectedOutput: []string{
				"Status: Future",
				"Starts in: 2h 14m",
			},
		},
		{
			name: "Live Session",
			session: &domain.Session{
				DateStart: start,
				DateEnd:   end,
			},
			now: start.Add(1 * time.Hour).Add(30 * time.Minute),
			expectedOutput: []string{
				"Status: Live",
				"Ends in: 0h 30m",
			},
		},
		{
			name: "Finished Session",
			session: &domain.Session{
				DateStart: start,
				DateEnd:   end,
			},
			// 1 day and 2 hours after end
			now: end.Add(26 * time.Hour),
			expectedOutput: []string{
				"Status: Finished",
				"Ended: 1d 2h 0m ago",
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
					t.Errorf("[%s] Expected to contain %q, but got:\n%s", tt.name, expected, output)
				}
			}
		})
	}
}
