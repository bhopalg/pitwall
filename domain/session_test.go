package domain

import (
	"testing"
	"time"
)

func TestSession_State(t *testing.T) {
	start := time.Date(2025, 3, 15, 10, 0, 0, 0, time.UTC) // 10:00 AM
	end := time.Date(2025, 3, 15, 11, 0, 0, 0, time.UTC)   // 11:00 AM

	tests := []struct {
		name      string
		dateStart time.Time
		dateEnd   time.Time
		now       time.Time
		want      SessionState
	}{
		{
			name:      "Future: now is well before start",
			dateStart: start,
			dateEnd:   end,
			now:       start.Add(-1 * time.Hour),
			want:      StateFuture,
		},
		{
			name:      "Live: now is exactly at start time",
			dateStart: start,
			dateEnd:   end,
			now:       start,
			want:      StateLive,
		},
		{
			name:      "Live: now is during the session",
			dateStart: start,
			dateEnd:   end,
			now:       start.Add(30 * time.Minute),
			want:      StateLive,
		},
		{
			name:      "Live: now is just before end time",
			dateStart: start,
			dateEnd:   end,
			now:       end.Add(-1 * time.Second),
			want:      StateLive,
		},
		{
			name:      "Live: end date is missing (zero value)",
			dateStart: start,
			dateEnd:   time.Time{}, // Zero value
			now:       start.Add(2 * time.Hour),
			want:      StateLive,
		},
		{
			name:      "Finished: now is exactly at end time",
			dateStart: start,
			dateEnd:   end,
			now:       end,
			want:      StateFinished,
		},
		{
			name:      "Finished: now is after end time",
			dateStart: start,
			dateEnd:   end,
			now:       end.Add(1 * time.Minute),
			want:      StateFinished,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{
				DateStart: tt.dateStart,
				DateEnd:   tt.dateEnd,
			}
			if got := s.State(tt.now); got != tt.want {
				t.Errorf("Session.State() = %v, want %v (Now: %v, Start: %v, End: %v)",
					got, tt.want, tt.now, tt.dateStart, tt.dateEnd)
			}
		})
	}
}
