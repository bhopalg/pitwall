package weekend

import (
	"context"
	"fmt"
	"time"

	"github.com/bhopalg/pitwall/domain"
	"github.com/bhopalg/pitwall/internal/cache"
	"github.com/bhopalg/pitwall/internal/openf1"
	"github.com/bhopalg/pitwall/utils"
)

type WeekendProvider interface {
	GetSessions(ctx context.Context, country_name, year string) (*[]openf1.Session, error)
}

type WeekendResponse struct {
	Sessions *[]domain.Session
	Warning  string
}

type WeekendService struct {
	openf1Client WeekendProvider
	cache        cache.Cache
}

func New(openf1Client WeekendProvider, cache cache.Cache) *WeekendService {
	return &WeekendService{
		openf1Client: openf1Client,
		cache:        cache,
	}
}

func (w *WeekendService) Weekend(ctx context.Context, country_name, year string) (WeekendResponse, error) {
	cacheKey := "weekend:" + country_name + ":" + year
	var cachedSessions []domain.Session

	found, isStale, _ := w.cache.Get(cacheKey, &cachedSessions)

	if found && !isStale {
		return WeekendResponse{
			Sessions: &cachedSessions,
		}, nil
	}

	apiSessions, err := w.openf1Client.GetSessions(ctx, country_name, year)
	if err != nil && found {
		return WeekendResponse{
			Sessions: &cachedSessions,
			Warning:  "⚠️ API unavailable. Showing stale cached data.",
		}, nil
	}

	if err != nil {
		return WeekendResponse{}, err
	}

	sessions := mapToDomain(apiSessions)

	if len(*sessions) == 0 {
		return WeekendResponse{}, nil
	}

	_ = w.cache.Set(cacheKey, sessions, 24*time.Hour)
	return WeekendResponse{Sessions: sessions}, nil
}

func mapToDomain(apiSession *[]openf1.Session) *[]domain.Session {
	now := time.Now().UTC()

	var sessions []domain.Session
	for _, session := range *apiSession {
		date_start, err := utils.ParseDate(session.DateStart)
		if err != nil {
			fmt.Errorf("error parsing date: %v", err)
			continue
		}

		date_end, err := utils.ParseDate(session.DateEnd)
		if err != nil {
			fmt.Errorf("error parsing date: %v", err)
			continue
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

	return &sessions
}
