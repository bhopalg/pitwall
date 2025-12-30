package getsession

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bhopalg/pitwall/domain"
	"github.com/bhopalg/pitwall/internal/cache"
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

func (m *mockCache) Clear() (int, error) {
	m.storage = make(map[string]interface{})
	return 0, nil
}

func (m *mockCache) Info() ([]cache.InfoEntry, string, error) {
	return []cache.InfoEntry{}, "", nil
}

type mockClient struct {
	fn     func(ctx context.Context, country, session, year string) (*openf1.Session, error)
	called bool
}

func (m *mockClient) GetSession(ctx context.Context, country, session, year string) (*openf1.Session, error) {
	m.called = true
	return m.fn(ctx, country, session, year)
}

func TestGetSession(t *testing.T) {
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
			name: "Cache Hit - Fresh (No Repo Call)",
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
			name: "Cache Miss - Successful API Call",
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
			name:            "Stale Cache + API Failure (Fallback)",
			mockErr:         errors.New("API down"),
			cacheFound:      true,
			cacheStale:      true,
			expectedError:   false,
			expectedWarning: "⚠️ API unavailable. Showing stale cached data.",
			expectRepoCall:  true,
		},
		{
			name: "Fail on invalid date format",
			mockResp: &openf1.Session{
				DateStart: "02-03-2024",
			},
			cacheFound:     false,
			expectedError:  true,
			expectRepoCall: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			mClient := &mockClient{
				fn: func(ctx context.Context, country, session, year string) (*openf1.Session, error) {
					return tc.mockResp, tc.mockErr
				},
			}

			mCache := &mockCache{
				storage: make(map[string]interface{}),
				found:   tc.cacheFound,
				isStale: tc.cacheStale,
			}

			if tc.cacheFound {
				key := "getsession:Bahrain:2024:Race"
				mCache.storage[key] = []domain.Session{{SessionName: "Cached Race"}}
			}

			s := New(mClient, mCache)
			res, err := s.GetSession(context.Background(), "Bahrain", "Race", "2024")

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
				t.Error("expected a session in response, got nil")
			}
		})
	}
}
