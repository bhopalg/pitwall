package latest

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bhopalg/pitwall/domain"
	"github.com/bhopalg/pitwall/internal/openf1"
)

type mockCache struct {
	storage map[string]interface{}
	found   bool
	isStale bool
}

func (m *mockCache) Get(key string, target interface{}) (bool, bool, error) {
	if !m.found {
		return false, false, nil
	}
	if data, ok := m.storage[key]; ok {
		if sessions, ok := data.([]domain.Session); ok {
			*(target.(*[]domain.Session)) = sessions
		}
	}
	return m.found, m.isStale, nil
}

func (m *mockCache) Set(key string, value interface{}, ttl time.Duration) error {
	return nil
}

type mockClient struct {
	fn     func(ctx context.Context) (*openf1.Session, error)
	called bool
}

func (m *mockClient) Next(ctx context.Context) (*openf1.Session, error) {
	m.called = true
	return m.fn(ctx)
}

func TestNext(t *testing.T) {
	testcases := []struct {
		name            string
		mockResp        *openf1.Session
		mockErr         error
		cacheFound      bool
		cacheStale      bool
		expectedError   bool
		expectedWarning string
		expectRepoCall  bool
	}{
		{
			name: "Cache Hit - Fresh (Repo not called)",
			mockResp: &openf1.Session{
				SessionName: "Race",
				DateStart:   "2024-03-02T15:00:00Z",
				DateEnd:     "2024-03-02T17:00:00Z",
			},
			cacheFound:     true,
			cacheStale:     false,
			expectedError:  false,
			expectRepoCall: false,
		},
		{
			name: "Cache Miss - Call Repo Success",
			mockResp: &openf1.Session{
				SessionName: "Race",
				DateStart:   "2024-03-02T15:00:00Z",
				DateEnd:     "2024-03-02T17:00:00Z",
			},
			cacheFound:     false,
			expectedError:  false,
			expectRepoCall: true,
		},
		{
			name:            "Stale Cache + Repo Failure (Fallback)",
			mockErr:         errors.New("API failure"),
			cacheFound:      true,
			cacheStale:      true,
			expectedError:   false,
			expectedWarning: "⚠️ API unavailable. Showing stale cached data.",
			expectRepoCall:  true,
		},
		{
			name: "Invalid Date Format",
			mockResp: &openf1.Session{
				DateStart: "wrong-format",
			},
			cacheFound:     false,
			expectedError:  true,
			expectRepoCall: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			mClient := &mockClient{
				fn: func(ctx context.Context) (*openf1.Session, error) {
					return tc.mockResp, tc.mockErr
				},
			}

			mCache := &mockCache{
				storage: make(map[string]interface{}),
				found:   tc.cacheFound,
				isStale: tc.cacheStale,
			}

			if tc.cacheFound {
				mCache.storage["latest"] = []domain.Session{{SessionName: "Cached Race"}}
			}

			s := New(mClient, mCache)
			res, err := s.Next(context.Background())

			if (err != nil) != tc.expectedError {
				t.Fatalf("expected error: %v, got: %v", tc.expectedError, err)
			}

			if mClient.called != tc.expectRepoCall {
				t.Errorf("expected repo call: %v, but was: %v", tc.expectRepoCall, mClient.called)
			}

			if res.Warning != tc.expectedWarning {
				t.Errorf("expected warning %q, got %q", tc.expectedWarning, res.Warning)
			}

			if !tc.expectedError && res.Session == nil {
				t.Error("expected session, got nil")
			}
		})
	}
}
