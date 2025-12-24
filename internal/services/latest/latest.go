package latest

import (
	"context"
	"time"

	"github.com/bhopalg/pitwall/domain"
	"github.com/bhopalg/pitwall/internal/openf1"
	"github.com/bhopalg/pitwall/utils"
)

type NextSessionProivder interface {
	Next(ctx context.Context) (*openf1.Session, error)
}

type NextSessionService struct {
	openf1Client NextSessionProivder
}

func New(openf1Client NextSessionProivder) *NextSessionService {
	return &NextSessionService{openf1Client: openf1Client}
}

func (n *NextSessionService) Next(ctx context.Context) (*domain.Session, error) {
	session, err := n.openf1Client.Next(ctx)
	if err != nil {
		return nil, err
	}

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
