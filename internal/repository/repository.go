package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	Competition CompetitionInterface
	Club        ClubInterface
	Season      SeasonInterface
	Match       MatchInterface
	User        UserInterface
	Token       TokenInterface
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		Competition: NewCompetitionRepository(db),
		Club:        NewClubRepository(db),
		Season:      NewSeasonRepository(db),
		Match:       NewMatchRepository(db),
		User:        NewUserRepository(db),
		Token:       NewTokenRepository(db),
	}
}
