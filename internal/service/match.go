package service

import (
	"context"
	"timeline/internal/domain"
	"timeline/internal/repository"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MatchService struct {
	r     repository.MatchInterface
	cache CacheInterface
}

func NewMatchService(r repository.MatchInterface, cache CacheInterface) *MatchService {
	return &MatchService{
		r:     r,
		cache: cache,
	}
}

//go:generate mockgen -source=match.go -destination=mock/match.go

type MatchInterface interface {
	GetMatches(ctx *gin.Context, competitionId, seasonId int, queryParams domain.GetMatchesQueryParams) ([]domain.MatchDTO, error)
	GetMatchById(ctx context.Context, matchId int) (domain.Match, error)
	GetMatchesDates(ctx *gin.Context, competitionId, seasonId int, queryParams domain.GetDatesQueryParams) ([]primitive.DateTime, int, error)
}

//--

func (s *MatchService) GetMatches(ctx *gin.Context, competitionId, seasonId int, queryParams domain.GetMatchesQueryParams) ([]domain.MatchDTO, error) {
	if iRes, err := s.cache.Get(ctx.Request.URL.String()); err == nil {
		return iRes.([]domain.MatchDTO), nil
	}

	res, err := s.r.GetMatches(ctx, competitionId, seasonId, queryParams)

	if err == nil {
		s.cache.Set(ctx.Request.URL.String(), res)
	}

	return res, err
}

//--

func (s *MatchService) GetMatchById(ctx context.Context, matchId int) (domain.Match, error) {
	return s.r.GetMatchById(ctx, matchId)
}

//--

func (s *MatchService) GetMatchesDates(ctx *gin.Context, competitionId, seasonId int, queryParams domain.GetDatesQueryParams) ([]primitive.DateTime, int, error) {
	if iRes, err := s.cache.Get(ctx.Request.URL.String()); err == nil {
		res := iRes.([]primitive.DateTime)
		return res, len(res), nil
	}

	res, resTotalCount, err := s.r.GetMatchesDates(ctx, competitionId, seasonId, queryParams)

	if err == nil {
		s.cache.Set(ctx.Request.URL.String(), res)
	}

	return res, resTotalCount, err
}
