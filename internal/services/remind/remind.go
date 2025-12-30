package remind

import (
	"time"
)

// ShouldRemind checks if the duration until start is within the window (0, threshold]
func ShouldRemind(now, start time.Time, threshold int) (bool, time.Duration) {
	diff := start.Sub(now)
	if diff > 0 && diff <= time.Duration(threshold)*time.Minute {
		return true, diff
	}

	return false, diff
}
