package domain

import "time"

type SessionState string

const (
	StateFuture   SessionState = "Future"
	StateLive     SessionState = "Live"
	StateFinished SessionState = "Finished"
)

type Session struct {
	SessionKey   int
	SessionName  string
	DateStart    time.Time
	DateEnd      time.Time
	Location     string
	CountryName  string
	CircuitName  string
	MeetingKey   int
	Year         int
	SessionState SessionState
}

func (s *Session) State(now time.Time) SessionState {
	if now.Before(s.DateStart) {
		return StateFuture
	}

	if s.DateEnd.IsZero() || now.Before(s.DateEnd) {
		return StateLive
	}

	return StateFinished
}
