package domain

import "time"

type Session struct {
	SessionKey  int
	SessionName string
	DateStart   time.Time
	DateEnd     time.Time
	Location    string
	CountryName string
	CircuitName string
	MeetingKey  int
	Year        int
}
