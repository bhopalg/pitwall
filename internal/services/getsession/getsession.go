package getsession

import (
	"context"
	"time"

	"github.com/bhopalg/pitwall/domain"
	"github.com/bhopalg/pitwall/internal/cache"
	"github.com/bhopalg/pitwall/internal/openf1"
	"github.com/bhopalg/pitwall/utils"
)

type GetSessionResponse struct {
	Session *domain.Session
	Warning string
}

type GetSessionSessionProvider interface {
	GetSession(ctx context.Context, country_name, session_name, year string) (*openf1.Session, error)
}

type GetSessionService struct {
	openf1Client GetSessionSessionProvider
	cache        cache.Cache
}

func New(openf1Client GetSessionSessionProvider, cache cache.Cache) *GetSessionService {
	return &GetSessionService{openf1Client: openf1Client, cache: cache}
}

func (s *GetSessionService) GetSession(ctx context.Context, country_name, session_name, year string) (GetSessionResponse, error) {
	cacheKey := "getsession:" + country_name + ":" + year + ":" + session_name
	var cachedSessions []domain.Session

	found, isStale, _ := s.cache.Get(cacheKey, &cachedSessions)

	if found && !isStale {
		return GetSessionResponse{
			Session: &cachedSessions[0],
		}, nil
	}

	session, err := s.openf1Client.GetSession(ctx, country_name, session_name, year)
	if err != nil && found {
		return GetSessionResponse{
			Session: &cachedSessions[0],
			Warning: "⚠️ API unavailable. Showing stale cached data.",
		}, nil
	}

	if err != nil {
		return GetSessionResponse{}, err
	}

	mappedSession, err := utils.MapToDomain(session)
	if err != nil {
		return GetSessionResponse{}, err
	}

	_ = s.cache.Set(cacheKey, session, 24*time.Hour)
	return GetSessionResponse{
		Session: mappedSession,
	}, nil
}
