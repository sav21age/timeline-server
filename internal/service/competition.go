package service

import (
	"context"
	"timeline/internal/domain"
	"timeline/internal/repository"

	"github.com/gin-gonic/gin"
)

type CompetitionService struct {
	r repository.CompetitionInterface
	cache CacheInterface
}

func NewCompetitionService(r repository.CompetitionInterface, cache CacheInterface) *CompetitionService {
	return &CompetitionService{
		r: r,
		cache: cache,
	}
}

//go:generate mockgen -source=competition.go -destination=mock/competition.go

type CompetitionInterface interface {
	GetCompetitions(ctx *gin.Context) ([]domain.Competition, error)
	GetCompetitionById(ctx context.Context, id int) (domain.Competition, error)
}

//--

func (s *CompetitionService) GetCompetitions(ctx *gin.Context) ([]domain.Competition, error) {
	if i, err := s.cache.Get(ctx.Request.URL.String()); err == nil {
		return i.([]domain.Competition), nil
	}

	res, err := s.r.GetCompetitions(ctx)

	if err == nil {
		s.cache.Set(ctx.Request.URL.String(), res)
	}
	
	return res, err
}

//--

func (s *CompetitionService) GetCompetitionById(ctx context.Context, competitionId int) (domain.Competition, error) {
	return s.r.GetCompetitionById(ctx, competitionId)
}
