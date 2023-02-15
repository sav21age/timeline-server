package repository

import (
	"context"

	"timeline/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ClubRepository struct {
	// db *mongo.Collection
	db *mongo.Database
}

func NewClubRepository(db *mongo.Database) *ClubRepository {
	return &ClubRepository{
		// db: db.Collection("club"),
		db: db,
	}
}

type ClubInterface interface {
	GetClubs(ctx context.Context, queryParams domain.GetClubsQueryParams) ([]domain.ClubDTO, error)
	GetClubsTotalCount(ctx context.Context, queryParams domain.GetClubsQueryParams) (int, error)

	GetClubsBySeasonId(ctx context.Context, seasonId int) ([]domain.ClubDTO, error)
	GetClubById(ctx context.Context, clubId int) (domain.Club, error)
	GetClubsAreas(ctx context.Context) ([]domain.AreaClub, error)
}

//--

func ClubsQuery(queryParams domain.GetClubsQueryParams) (primitive.M, *options.FindOptions)  {
	o := options.Find()

	sort := bson.D{{Key: "area.name", Value: 1}, {Key: "name", Value: 1}}
	if queryParams.SortBy == "club" {
		sort = bson.D{{Key: "name", Value: 1}, {Key: "area.name", Value: 1}}
	}

	o.SetSort(sort)

	if (queryParams.Pagination != domain.Pagination{}) {
		o.SetSkip(int64(queryParams.Pagination.Skip))
		o.SetLimit(int64(queryParams.Pagination.Limit))
	}

	m := bson.M{}
	if queryParams.AreaId != 0 {
		m = bson.M{"area._id": queryParams.AreaId}
	}

	return m, o
}

func (r *ClubRepository) GetClubs(ctx context.Context, queryParams domain.GetClubsQueryParams) ([]domain.ClubDTO, error) {
	m, o := ClubsQuery(queryParams)

	cursor, err := r.db.Collection("club").Find(ctx, m, o)
	if err != nil {
		return []domain.ClubDTO{}, err
	}

	club := make([]domain.ClubDTO, 0)
	if err = cursor.All(ctx, &club); err != nil {
		return []domain.ClubDTO{}, err
	}

	return club, nil
}

func (r *ClubRepository) GetClubsTotalCount(ctx context.Context, queryParams domain.GetClubsQueryParams) (int, error) {
	m, _ := ClubsQuery(queryParams)
	totalCount, err := r.db.Collection("club").CountDocuments(ctx, m)
	
	if err != nil {
		return 0, err
	}

	return int(totalCount), nil
}

//--

func (r *ClubRepository) GetClubsBySeasonId(ctx context.Context, seasonId int) ([]domain.ClubDTO, error) {
	ids, err := r.db.Collection("match").Distinct(ctx, "homeTeam._id", bson.M{"season._id": seasonId}, options.Distinct())
	if err != nil {
		return []domain.ClubDTO{}, err
	}

	o := options.Find()
	o.SetSort(bson.D{{Key: "shortName", Value: 1}})

	cursor, err := r.db.Collection("club").Find(ctx, bson.M{"_id": bson.M{"$in": ids}}, o)
	if err != nil {
		return []domain.ClubDTO{}, err
	}
	
	club := make([]domain.ClubDTO, 0)
	if err = cursor.All(ctx, &club); err != nil {
		return []domain.ClubDTO{}, err
	}

	return club, nil
}

//--

func (r *ClubRepository) GetClubById(ctx context.Context, clubId int) (domain.Club, error) {
	var club domain.Club
	err := r.db.Collection("club").FindOne(ctx, bson.M{"_id": clubId}).Decode(&club)
	if err != nil {
		return domain.Club{}, err
	}

	return club, nil
}

//--

func (r *ClubRepository) GetClubsAreas(ctx context.Context) ([]domain.AreaClub, error) {

	ids, err := r.db.Collection("club").Distinct(ctx, "area._id", bson.M{}, options.Distinct())
	if err != nil {
		return []domain.AreaClub{}, err
	}

	cursor, err := r.db.Collection("area").Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return []domain.AreaClub{}, err
	}

	areas := make([]domain.AreaClub, 0)
	if err = cursor.All(ctx, &areas); err != nil {
		return []domain.AreaClub{}, err
	}

	return areas, nil
}
