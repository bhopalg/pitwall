package getsession

import (
	"context"
	"time"

	"github.com/bhopalg/pitwall/internal/openf1"
)

func GetSession(ctx context.Context) (*Session, error) {
	c := openf1.New()
	s, err := c.GetSession(ctx)

	date_start, err := parseDate(s.DateStart)
	if err != nil {
		return nil, err
	}

	date_end, err := parseDate(s.DateEnd)
	if err != nil {
		return nil, err
	}

	return &Session{
		SessionKey:  s.SessionKey,
		SessionName: s.SessionName,
		DateStart:   *date_start,
		DateEnd:     *date_end,
		Location:    s.Location,
		CountryName: s.CountryName,
		CircuitName: s.CircuitName,
		MeetingKey:  s.MeetingKey,
		Year:        s.Year,
	}, err
}

func parseDate(date string) (*time.Time, error) {
	parseDate, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return nil, err
	}

	return &parseDate, nil
}
