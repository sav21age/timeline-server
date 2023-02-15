package service

import (
	"context"
	"timeline/internal/domain"
	"timeline/internal/repository"

	"github.com/gin-gonic/gin"
)

type ClubService struct {
	r     repository.ClubInterface
	cache CacheInterface
}

func NewClubService(r repository.ClubInterface, cache CacheInterface) *ClubService {
	return &ClubService{
		r:     r,
		cache: cache,
	}
}

//go:generate mockgen -source=club.go -destination=mock/club.go

type ClubInterface interface {
	GetClubs(ctx *gin.Context, queryParams domain.GetClubsQueryParams) ([]domain.ClubDTO, int, error)
	GetClubsBySeasonId(ctx *gin.Context, seasonId int) ([]domain.ClubDTO, error)
	GetClubById(ctx context.Context, clubId int) (domain.Club, error)
	GetClubsAreas(ctx *gin.Context) ([]domain.AreaClub, error)
}

func (s *ClubService) GetClubs(ctx *gin.Context, queryParams domain.GetClubsQueryParams) ([]domain.ClubDTO, int, error) {
	var res = make([]domain.ClubDTO, 0)
	iRes, err := s.cache.Get(ctx.Request.URL.String())
	if err == nil {
		res = iRes.([]domain.ClubDTO)
	} else {
		res, err = s.r.GetClubs(ctx, queryParams)

		if err != nil {
			return []domain.ClubDTO{}, 0, err
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
		resTotalCount, err = s.r.GetClubsTotalCount(ctx, queryParams)

		if err != nil {
			return []domain.ClubDTO{}, 0, err
		} else {
			s.cache.Set(cacheName, resTotalCount)	
		}
	}

	return res, resTotalCount, err
}

//--

func (s *ClubService) GetClubsBySeasonId(ctx *gin.Context, seasonId int) ([]domain.ClubDTO, error) {
	if iRes, err := s.cache.Get(ctx.Request.URL.String()); err == nil {
		return iRes.([]domain.ClubDTO), nil
	}

	res, err := s.r.GetClubsBySeasonId(ctx, seasonId)

	if err == nil {
		s.cache.Set(ctx.Request.URL.String(), res)
	}

	return res, err
}

//--

func (s *ClubService) GetClubById(ctx context.Context, clubId int) (domain.Club, error) {
	return s.r.GetClubById(ctx, clubId)
}

//--

func (s *ClubService) GetClubsAreas(ctx *gin.Context) ([]domain.AreaClub, error) {
	if iRes, err := s.cache.Get(ctx.Request.URL.String()); err == nil {
		return iRes.([]domain.AreaClub), nil
	}

	res, err := s.r.GetClubsAreas(ctx)

	if err == nil {
		s.cache.Set(ctx.Request.URL.String(), res)
	}

	return res, err
}
