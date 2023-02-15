package service

import (
	"timeline/config"
	"timeline/internal/repository"
)

type Service struct {
	Competition CompetitionInterface
	Club        ClubInterface
	Season      SeasonInterface
	Match       MatchInterface
	User        UserInterface
}

func NewService(r *repository.Repository, config *config.Config) *Service {
	EmailService := NewEmailService(config)
	TokenService := NewTokenService(r.Token, config)
	CacheService := NewCacheService(config)

	return &Service{
		Competition: NewCompetitionService(r.Competition, CacheService),
		Club:        NewClubService(r.Club, CacheService),
		Season:      NewSeasonService(r.Season, CacheService),
		Match:       NewMatchService(r.Match, CacheService),
		User: NewUserService(r.User, EmailService, TokenService, config),
	}
}
