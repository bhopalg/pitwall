package getsession

import (
	"context"
	"time"

	"github.com/bhopalg/pitwall/domain"
	"github.com/bhopalg/pitwall/internal/openf1"
	"github.com/bhopalg/pitwall/utils"
)

type GetSessionSessionProvider interface {
	GetSession(ctx context.Context, country_name, session_name, year string) (*openf1.Session, error)
}

type GetSessionService struct {
	openf1Client GetSessionSessionProvider // Use interface here
}

func New(openf1Client GetSessionSessionProvider) *GetSessionService {
	return &GetSessionService{openf1Client: openf1Client}
}

func (s *GetSessionService) GetSession(ctx context.Context, country_name, session_name, year string) (*domain.Session, error) {
	session, err := s.openf1Client.GetSession(ctx, country_name, session_name, year)

	date_start, err := utils.ParseDate(session.DateStart)
	if err != nil {
		return nil, err
	}

	date_end, err := utils.ParseDate(session.DateEnd)
	if err != nil {
		return nil, err
	}

	mappedSession := &domain.Session{
		SessionKey:  session.SessionKey,
		SessionName: session.SessionName,
		DateStart:   *date_start,
		DateEnd:     *date_end,
		Location:    session.Location,
		CountryName: session.CountryName,
		CircuitName: session.CircuitName,
		MeetingKey:  session.MeetingKey,
		Year:        session.Year,
	}

	now := time.Now().UTC()
	mappedSession.SessionState = mappedSession.State(now)

	return mappedSession, nil
}
