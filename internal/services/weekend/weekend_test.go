package weekend

import (
	"context"
	"errors"
	"testing"

	"github.com/bhopalg/pitwall/internal/openf1"
)

// MockOpenF1 to satisfy the WeekendProvider interface
type MockOpenF1 struct {
	sessions *[]openf1.Session
	err      error
}

func (m *MockOpenF1) GetSessions(ctx context.Context, country, year string) (*[]openf1.Session, error) {
	return m.sessions, m.err
}

func TestWeekendService_Weekend(t *testing.T) {
	testcases := []struct {
		name          string
		mockResp      []openf1.Session
		mockErr       error
		expectedError bool
		expectedLen   int
	}{
		{
			name: "Success - Multiple Sessions",
			mockResp: []openf1.Session{
				{SessionName: "Practice 1", DateStart: "2023-07-28T11:30:00Z", DateEnd: "2023-07-28T12:30:00Z"},
				{SessionName: "Practice 2", DateStart: "2023-07-28T15:00:00Z", DateEnd: "2023-07-28T16:00:00Z"},
			},
			mockErr:       nil,
			expectedError: false,
			expectedLen:   2,
		},
		{
			name:          "API Error",
			mockResp:      nil,
			mockErr:       errors.New("network failure"),
			expectedError: true,
			expectedLen:   0,
		},
		{
			name: "Invalid Date Format",
			mockResp: []openf1.Session{
				{SessionName: "Broken Date", DateStart: "invalid-date"},
			},
			mockErr:       nil,
			expectedError: true,
			expectedLen:   0,
		},
		{
			name:          "No Sessions Found",
			mockResp:      []openf1.Session{},
			mockErr:       nil,
			expectedError: false,
			expectedLen:   0,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := &MockOpenF1{
				sessions: &tc.mockResp,
				err:      tc.mockErr,
			}
			service := New(mockClient)

			sessions, err := service.Weekend(context.Background(), "Belgium", "2023")

			if (err != nil) != tc.expectedError {
				t.Fatalf("expected error: %v, got: %v", tc.expectedError, err)
			}

			if !tc.expectedError {
				if tc.expectedLen == 0 && sessions != nil {
					t.Errorf("expected nil sessions for empty list, got %d", len(*sessions))
				} else if tc.expectedLen > 0 {
					if sessions == nil || len(*sessions) != tc.expectedLen {
						t.Errorf("expected %d sessions, got %v", tc.expectedLen, sessions)
					}
				}
			}
		})
	}
}
