package getsession

import (
	"context"
	"time"

	"github.com/bhopalg/pitwall/domain"
	"github.com/bhopalg/pitwall/internal/openf1"
)

type SessionProvider interface {
	GetSession(ctx context.Context, country_name, session_name, year string) (*openf1.Session, error)
}

type Service struct {
	openf1Client SessionProvider // Use interface here
}

func New(openf1Client SessionProvider) *Service {
	return &Service{openf1Client: openf1Client}
}

func (s *Service) GetSession(ctx context.Context, country_name, session_name, year string) (*domain.Session, error) {
	session, err := s.openf1Client.GetSession(ctx, country_name, session_name, year)

	date_start, err := parseDate(session.DateStart)
	if err != nil {
		return nil, err
	}

	date_end, err := parseDate(session.DateEnd)
	if err != nil {
		return nil, err
	}

	return &domain.Session{
		SessionKey:  session.SessionKey,
		SessionName: session.SessionName,
		DateStart:   *date_start,
		DateEnd:     *date_end,
		Location:    session.Location,
		CountryName: session.CountryName,
		CircuitName: session.CircuitName,
		MeetingKey:  session.MeetingKey,
		Year:        session.Year,
	}, err
}

func parseDate(date string) (*time.Time, error) {
	parseDate, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return nil, err
	}

	return &parseDate, nil
}
