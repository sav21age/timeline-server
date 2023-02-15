package repository

import (
	"context"

	"timeline/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CompetitionRepository struct {
	// db *mongo.Collection
	db *mongo.Database
}

func NewCompetitionRepository(db *mongo.Database) *CompetitionRepository {
	return &CompetitionRepository{
		db: db,
	}
}

type CompetitionInterface interface {
	GetCompetitions(ctx context.Context) ([]domain.Competition, error)
	GetCompetitionById(ctx context.Context, competitionId int) (domain.Competition, error)
}

func (r *CompetitionRepository) GetCompetitions(ctx context.Context) ([]domain.Competition, error) {
	cursor, err := r.db.Collection("competition").Find(ctx, bson.M{})
	if err != nil {
		return []domain.Competition{}, err
	}

	competitions := make([]domain.Competition, 0)
	err = cursor.All(ctx, &competitions)
	if err != nil {
		return []domain.Competition{}, err
	}

	return competitions, nil
}

//--

func (r *CompetitionRepository) GetCompetitionById(ctx context.Context, competitionId int) (domain.Competition, error) {
	var competition domain.Competition
	err := r.db.Collection("competition").FindOne(ctx, bson.M{"_id": competitionId}).Decode(&competition)
	if err != nil {
		return domain.Competition{}, err
	}

	return competition, nil
}
