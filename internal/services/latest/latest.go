package latest

import (
	"context"
	"time"

	"github.com/bhopalg/pitwall/domain"
	"github.com/bhopalg/pitwall/internal/cache"
	"github.com/bhopalg/pitwall/internal/openf1"
	"github.com/bhopalg/pitwall/utils"
)

type LatestResponse struct {
	Session *domain.Session
	Warning string
}

type NextSessionProivder interface {
	Next(ctx context.Context) (*openf1.Session, error)
}

type NextSessionService struct {
	openf1Client NextSessionProivder
	cache        cache.Cache
}

func New(openf1Client NextSessionProivder, cache cache.Cache) *NextSessionService {
	return &NextSessionService{openf1Client: openf1Client, cache: cache}
}

func (n *NextSessionService) Next(ctx context.Context) (LatestResponse, error) {
	cacheKey := "latest"
	var cachedSessions []domain.Session

	found, isStale, _ := n.cache.Get(cacheKey, &cachedSessions)

	if found && !isStale {
		return LatestResponse{
			Session: &cachedSessions[0],
		}, nil
	}

	session, err := n.openf1Client.Next(ctx)
	if err != nil && found {
		return LatestResponse{
			Session: &cachedSessions[0],
			Warning: "⚠️ API unavailable. Showing stale cached data.",
		}, nil
	}

	if err != nil {
		return LatestResponse{}, err
	}

	mappedSession, err := utils.MapToDomain(session)
	if err != nil {
		return LatestResponse{}, err
	}

	_ = n.cache.Set(cacheKey, session, 24*time.Hour)
	return LatestResponse{
		Session: mappedSession,
	}, nil
}
