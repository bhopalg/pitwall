package weekend

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bhopalg/pitwall/domain"
	"github.com/bhopalg/pitwall/internal/cache"
	"github.com/bhopalg/pitwall/internal/openf1"
)

// MockCache implements the cache.Cache interface for testing
type MockCache struct {
	storage map[string]interface{}
	isStale bool
	found   bool
}

func (m *MockCache) Get(key string, target interface{}) (bool, bool, error) {
	if !m.found {
		return false, false, nil
	}
	if data, ok := m.storage[key]; ok {
		if sessions, ok := data.(*[]domain.Session); ok {
			*(target.(*[]domain.Session)) = *sessions
		}
	}
	return m.found, m.isStale, nil
}

func (m *MockCache) Set(key string, value interface{}, ttl time.Duration) error {
	m.storage[key] = value
	return nil
}

func (m *MockCache) Clear() (int, error) {
	m.storage = make(map[string]interface{})
	return 0, nil
}

func (m *MockCache) Info() ([]cache.InfoEntry, string, error) {
	return []cache.InfoEntry{}, "", nil
}

type MockOpenF1 struct {
	sessions *[]openf1.Session
	err      error
	called   bool
}

func (m *MockOpenF1) GetSessions(ctx context.Context, country, year string) (*[]openf1.Session, error) {
	m.called = true
	return m.sessions, m.err
}

func TestWeekendService_Weekend(t *testing.T) {
	testcases := []struct {
		name            string
		mockResp        []openf1.Session
		mockErr         error
		cacheFound      bool
		cacheStale      bool
		expectedError   bool
		expectedWarning string
		expectedLen     int
		expectRepoCall  bool
	}{
		{
			name: "Cache Hit - Fresh (Repo not called)",
			mockResp: []openf1.Session{
				{SessionName: "Practice 1", DateStart: "2023-07-28T11:30:00Z", DateEnd: "2023-07-28T12:30:00Z"},
			},
			cacheFound:     true,
			cacheStale:     false,
			expectedError:  false,
			expectedLen:    1,
			expectRepoCall: false,
		},
		{
			name: "Cache Miss - Call Repo Success",
			mockResp: []openf1.Session{
				{SessionName: "Practice 1", DateStart: "2023-07-28T11:30:00Z", DateEnd: "2023-07-28T12:30:00Z"},
			},
			cacheFound:     false,
			expectedError:  false,
			expectedLen:    1,
			expectRepoCall: true,
		},
		{
			name:            "Stale Cache + Repo Failure (Fallback)",
			mockErr:         errors.New("api down"),
			cacheFound:      true,
			cacheStale:      true,
			expectedError:   false,
			expectedWarning: "⚠️ API unavailable. Showing stale cached data.",
			expectedLen:     0, // Length depends on what's in mock storage
			expectRepoCall:  true,
		},
		{
			name:           "API Error - No Cache",
			mockErr:        errors.New("network failure"),
			cacheFound:     false,
			expectedError:  true,
			expectRepoCall: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup Mocks
			mockClient := &MockOpenF1{
				sessions: &tc.mockResp,
				err:      tc.mockErr,
			}
			mockCache := &MockCache{
				storage: make(map[string]interface{}),
				found:   tc.cacheFound,
				isStale: tc.cacheStale,
			}

			// If cache is "found", pre-populate it with dummy data for the test
			if tc.cacheFound {
				dummy := []domain.Session{{SessionName: "Cached Session"}}
				mockCache.storage["weekend:Belgium:2023"] = &dummy
			}

			service := New(mockClient, mockCache)

			// Execute
			resp, err := service.Weekend(context.Background(), "Belgium", "2023")

			// Assertions
			if (err != nil) != tc.expectedError {
				t.Fatalf("expected error: %v, got: %v", tc.expectedError, err)
			}

			if mockClient.called != tc.expectRepoCall {
				t.Errorf("expected repo call: %v, but was: %v", tc.expectRepoCall, mockClient.called)
			}

			if resp.Warning != tc.expectedWarning {
				t.Errorf("expected warning %q, got %q", tc.expectedWarning, resp.Warning)
			}

			if !tc.expectedError && tc.expectedLen > 0 {
				if resp.Sessions == nil || len(*resp.Sessions) != tc.expectedLen {
					t.Errorf("expected %d sessions, got %v", tc.expectedLen, resp.Sessions)
				}
			}
		})
	}
}
