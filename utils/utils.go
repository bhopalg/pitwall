package utils

import (
	"fmt"
	"time"

	"github.com/bhopalg/pitwall/domain"
)

func ParseDate(date string) (*time.Time, error) {
	pd, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return nil, err
	}

	return &pd, nil
}

func PrintSessionStatus(s *domain.Session, now time.Time) {
	state := s.State(now)
	loc, _ := time.LoadLocation("Europe/London")

	fmt.Printf("Status: %s\n", state)

	switch state {
	case domain.StateFuture:
		startTimeUK := s.DateStart.In(loc)
		startsAtStr := startTimeUK.Format("Mon 02 Jan 2006, 15:04")
		diff := s.DateStart.Sub(now).Round(time.Minute)
		fmt.Printf("Starts at: %s (UK)\n", startsAtStr)
		fmt.Printf("Starts in: %s\n", formatDuration(diff))

	case domain.StateLive:
		if !s.DateEnd.IsZero() {
			endTime := s.DateEnd.In(loc).Format("15:04")
			diff := s.DateEnd.Sub(now).Round(time.Minute)
			fmt.Printf("Ends at: %s (UK)\n", endTime)
			fmt.Printf("Ends in: %s\n", formatDuration(diff))
		} else {
			fmt.Println("Ends at: TBD")
		}

	case domain.StateFinished:
		endTime := s.DateEnd.In(loc).Format("15:04")
		diff := now.Sub(s.DateEnd).Round(time.Minute)
		fmt.Printf("Ended at: %s (UK)\n", endTime)
		fmt.Printf("Ended: %s ago\n", formatDuration(diff))
	}
}

func formatDuration(d time.Duration) string {
	days := int(d.Hours()) / 24
	h := int(d.Hours()) % 24
	m := int(d.Minutes()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, h, m)
	}
	return fmt.Sprintf("%dh %dm", h, m)
}
