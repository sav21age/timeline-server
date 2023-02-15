package repository

import (
	"context"
	"timeline/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SeasonRepository struct {
	// db *mongo.Collection
	db *mongo.Database
}

func NewSeasonRepository(db *mongo.Database) *SeasonRepository {
	return &SeasonRepository{
		// db: db.Collection("season"),
		db: db,
	}
}

type SeasonInterface interface {
	GetSeasons(ctx context.Context, competitionId int, queryParams domain.GetSeasonsQueryParams) ([]domain.Season, error)
	GetSeasonsTotalCount(ctx context.Context, competitionId int, queryParams domain.GetSeasonsQueryParams) (int, error)
	GetSeasonById(ctx context.Context, seasonId int) (domain.Season, error)
	GetSeasonByMatchId(ctx context.Context, matchId int) (domain.Season, error)
}

//--

func (r *SeasonRepository) GetSeasons(ctx context.Context, competitionId int, queryParams domain.GetSeasonsQueryParams) ([]domain.Season, error) {
	o := options.Find()
	o.SetSort(bson.D{{Key: "startDate", Value: -1}})

	if (queryParams.Pagination != domain.Pagination{}) {
		o.SetSkip(int64(queryParams.Pagination.Skip))
		o.SetLimit(int64(queryParams.Pagination.Limit))
	}

	cursor, err := r.db.Collection("season").Find(ctx, bson.M{"competition._id": competitionId}, o)
	if err != nil {
		return []domain.Season{}, err
	}

	seasons := make([]domain.Season, 0)
	err = cursor.All(ctx, &seasons)
	if err != nil {
		return []domain.Season{}, err
	}

	return seasons, nil
}

func (r *SeasonRepository) GetSeasonsTotalCount(ctx context.Context, competitionId int, queryParams domain.GetSeasonsQueryParams) (int, error) {
	totalCount, err := r.db.Collection("season").CountDocuments(ctx, bson.M{"competition._id": competitionId})

	if err != nil {
		return 0, err
	}

	return int(totalCount), nil
}

//--

func (r *SeasonRepository) GetSeasonById(ctx context.Context, seasonId int) (domain.Season, error) {
	var season domain.Season
	err := r.db.Collection("season").FindOne(ctx, bson.M{"_id": seasonId}).Decode(&season)
	if err != nil {
		return domain.Season{}, err
	}

	return season, nil
}

//--

func (r *SeasonRepository) GetSeasonByMatchId(ctx context.Context, matchId int) (domain.Season, error) {
	var match domain.Match
	err := r.db.Collection("match").FindOne(ctx, bson.M{"_id": matchId}).Decode(&match)
	if err != nil {
		return domain.Season{}, err
	}

	var season domain.Season
	err = r.db.Collection("season").FindOne(ctx, bson.M{"_id": match.Season.ID}).Decode(&season)
	if err != nil {
		return domain.Season{}, err
	}

	return season, nil
}
