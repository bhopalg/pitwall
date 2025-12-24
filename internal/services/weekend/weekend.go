package weekend

import (
	"context"
	"time"

	"github.com/bhopalg/pitwall/domain"
	"github.com/bhopalg/pitwall/internal/openf1"
	"github.com/bhopalg/pitwall/utils"
)

type WeekendProvider interface {
	GetSessions(ctx context.Context, country_name, year string) (*[]openf1.Session, error)
}

type WeekendService struct {
	openf1Client WeekendProvider
}

func New(openf1Client WeekendProvider) *WeekendService {
	return &WeekendService{
		openf1Client: openf1Client,
	}
}

func (w *WeekendService) Weekend(ctx context.Context, country_name, year string) (*[]domain.Session, error) {
	now := time.Now().UTC()

	apiSessions, err := w.openf1Client.GetSessions(ctx, country_name, year)
	if err != nil {
		return nil, err
	}

	var sessions []domain.Session
	for _, session := range *apiSessions {
		date_start, err := utils.ParseDate(session.DateStart)
		if err != nil {
			return nil, err
		}

		date_end, err := utils.ParseDate(session.DateEnd)
		if err != nil {
			return nil, err
		}

		mappedSession := domain.Session{
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

		mappedSession.SessionState = mappedSession.State(now)

		sessions = append(sessions, mappedSession)
	}

	if len(sessions) == 0 {
		return nil, nil
	}

	return &sessions, nil
}
