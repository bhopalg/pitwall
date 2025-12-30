package remind

import (
	"testing"
	"time"
)

func TestShouldRemind(t *testing.T) {
	now := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	threshold := 30

	testcases := []struct {
		name     string
		start    time.Time
		expected bool
	}{
		{
			name:     "29 minutes before -> Trigger",
			start:    now.Add(29 * time.Minute),
			expected: true,
		},
		{
			name:     "31 minutes before -> No trigger",
			start:    now.Add(31 * time.Minute),
			expected: false,
		},
		{
			name:     "Exactly 30 minutes -> Trigger",
			start:    now.Add(30 * time.Minute),
			expected: true,
		},
		{
			name:     "5 minutes after start -> No trigger",
			start:    now.Add(-5 * time.Minute),
			expected: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got, _ := ShouldRemind(now, tc.start, threshold)
			if got != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, got)
			}
		})
	}
}
