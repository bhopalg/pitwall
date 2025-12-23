package next

import (
	"context"
	"testing"

	"github.com/bhopalg/pitwall/internal/openf1"
)

type mockClient struct {
	fn func(ctx context.Context) (*openf1.Session, error)
}

func (m *mockClient) Next(ctx context.Context) (*openf1.Session, error) {
	return m.fn(ctx)
}

func TestNext(t *testing.T) {
	testcases := []struct {
		name          string
		mockResp      *openf1.Session
		mockErr       error
		expectedError bool
	}{
		{
			name: "successful session retrieval",
			mockResp: &openf1.Session{
				SessionKey:  1,
				SessionName: "Race",
				DateStart:   "2024-03-02T15:00:00Z",
				DateEnd:     "2024-03-02T17:00:00Z",
			},
			mockErr:       nil,
			expectedError: false,
		},
		{
			name: "fail on invalid date format",
			mockResp: &openf1.Session{
				DateStart: "02-03-2024",
				DateEnd:   "2024-03-02T17:00:00Z",
			},
			mockErr:       nil,
			expectedError: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			m := &mockClient{
				fn: func(ctx context.Context) (*openf1.Session, error) {
					return tc.mockResp, tc.mockErr
				},
			}

			s := New(m)

			res, err := s.Next(context.Background())

			if (err != nil) != tc.expectedError {
				t.Fatalf("expected error: %v, got: %v", tc.expectedError, err)
			}

			if !tc.expectedError && res.SessionName != "Race" {
				t.Errorf("expected session Race, got %s", res.SessionName)
			}
		})
	}
}
