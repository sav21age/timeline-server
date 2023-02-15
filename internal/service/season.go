package service

import (
	"context"
	"timeline/internal/domain"
	"timeline/internal/repository"

	"github.com/gin-gonic/gin"
)

type SeasonService struct {
	r     repository.SeasonInterface
	cache CacheInterface
}

func NewSeasonService(r repository.SeasonInterface, cache CacheInterface) *SeasonService {
	return &SeasonService{
		r:     r,
		cache: cache,
	}
}

//go:generate mockgen -source=season.go -destination=mock/season.go

type SeasonInterface interface {
	GetSeasons(ctx *gin.Context, competitionId int, queryParams domain.GetSeasonsQueryParams) ([]domain.Season, int, error)
	GetSeasonById(ctx context.Context, seasonId int) (domain.Season, error)
	GetSeasonByMatchId(ctx *gin.Context, matchId int) (domain.Season, error)
}

//--

func (s *SeasonService) GetSeasons(ctx *gin.Context, competitionId int, queryParams domain.GetSeasonsQueryParams) ([]domain.Season, int, error) {
	var res = make([]domain.Season, 0)
	iRes, err := s.cache.Get(ctx.Request.URL.String())
	if err == nil {
		res = iRes.([]domain.Season)
	} else {
		res, err = s.r.GetSeasons(ctx, competitionId, queryParams)

		if err != nil {
			return []domain.Season{}, 0, err
		} else {
			s.cache.Set(ctx.Request.URL.String(), res)
		}
	}

	var resTotalCount int
	cacheName := ctx.Request.URL.String() +"&totalCount"
	iResTotalCount, err := s.cache.Get(cacheName)
	if err == nil {
		resTotalCount = iResTotalCount.(int)
	} else {
		resTotalCount, err = s.r.GetSeasonsTotalCount(ctx, competitionId, queryParams)

		if err != nil {
			return []domain.Season{}, 0, err
		} else {
			s.cache.Set(cacheName, resTotalCount)	
		}
	}

	return res, resTotalCount, err
}

//--

func (s *SeasonService) GetSeasonById(ctx context.Context, seasonId int) (domain.Season, error) {
	return s.r.GetSeasonById(ctx, seasonId)
}

//--

func (s *SeasonService) GetSeasonByMatchId(ctx *gin.Context, matchId int) (domain.Season, error) {
	if iRes, err := s.cache.Get(ctx.Request.URL.String()); err == nil {
		return iRes.(domain.Season), nil
	}

	res, err := s.r.GetSeasonByMatchId(ctx, matchId)

	if err == nil {
		s.cache.Set(ctx.Request.URL.String(), res)
	}
	
	return res, err
}
