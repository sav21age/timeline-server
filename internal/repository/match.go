package repository

import (
	"context"
	"timeline/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MatchRepository struct {
	// db *mongo.Collection
	db *mongo.Database
}

func NewMatchRepository(db *mongo.Database) *MatchRepository {
	return &MatchRepository{
		// db: db.Collection("match"),
		db: db,
	}
}

type MatchInterface interface {
	GetMatches(ctx context.Context, competitionId, seasonId int, queryParams domain.GetMatchesQueryParams) ([]domain.MatchDTO, error)
	GetMatchById(ctx context.Context, matchId int) (domain.Match, error)
	GetMatchesDates(ctx context.Context, competition_id, season_id int, queryParams domain.GetDatesQueryParams) ([]primitive.DateTime, int, error)
}

func (r *MatchRepository) GetMatches(ctx context.Context, competitionId, seasonId int, queryParams domain.GetMatchesQueryParams) ([]domain.MatchDTO, error) {
	// { "utcDate": {$gt: new Date('2023-06-04')} }

	// if (*matchAll.Datination != domain.Datination{}) {
	// 	q = bson.M{
	// 		"$and": []bson.M{
	// 			{"competition._id": competition_id},
	// 			{"season._id": season_id},
	// 			{"utcDate": bson.M{
	// 				"$gte": matchAll.Datination.MinDate,
	// 				"$lte": matchAll.Datination.MaxDate,
	// 			},
	// 			},
	// 		},
	// 	}
	// } else {
	// 	q = bson.M{
	// 		"$and": []bson.M{
	// 			{"competition._id": competition_id},
	// 			{"season._id": season_id},
	// 		},
	// 	}
	// }

	o := options.Find()
	o.SetSort(bson.D{{Key: "utcDate", Value: -1}})

	// var q primitive.M
	// https://stackoverflow.com/questions/67688353/how-to-append-to-bson-object
	// https://stackoverflow.com/questions/58365947/dynamic-bson-creation-in-golang-mongodb-driver

	f := []bson.M{
		{"competition._id": competitionId},
		{"season._id": seasonId},
	}

	if (queryParams.Datination != domain.Datination{}) {
		f = append(f, bson.M{"utcDate": bson.M{
			"$gte": queryParams.Datination.MinDate,
			"$lte": queryParams.Datination.MaxDate,
		}})
	}

	if queryParams.MatchStatus != "" {
		f = append(f, bson.M{"status": queryParams.MatchStatus})
	}

	if queryParams.ClubId != 0 {
		f = append(f, bson.M{"$or": []bson.M{
			{"homeTeam._id": queryParams.ClubId},
			{"awayTeam._id": queryParams.ClubId},
		}})
	}

	q := bson.M{
		"$and": f,
	}

	cursor, err := r.db.Collection("match").Find(ctx, q, o)
	if err != nil {
		return []domain.MatchDTO{}, err
	}

	matches := make([]domain.MatchDTO, 0)
	if err = cursor.All(ctx, &matches); err != nil {
		return []domain.MatchDTO{}, err
	}

	return matches, nil
}

//--

func (r *MatchRepository) GetMatchById(ctx context.Context, matchId int) (domain.Match, error) {
	var match domain.Match
	err := r.db.Collection("match").FindOne(ctx, bson.M{"_id": matchId}).Decode(&match)
	if err != nil {
		return domain.Match{}, err
	}

	return match, nil
}

//--

func (r *MatchRepository) GetMatchesDates(ctx context.Context, competitionId, seasonId int, queryParams domain.GetDatesQueryParams) ([]primitive.DateTime, int, error) {
	// [{'$match': {'competition._id': 2021,
	//             'season._id': 733,
	//             'status': 'FINISHED'
	//             }
	//  }, {'$match': {
	//             '$or': [{'homeTeam._id': 57},
	//                     {'awayTeam._id': 57}
	//                    ]
	//                }
	//  }, {'$sort': {'utcDate': -1}},
	//         {'$group': {'_id': 'utc',
	//                    'utcDate': {'$push': '$utcDate'}
	//                   }
	//         }
	// ]

	m := bson.D{
		{Key: "competition._id", Value: competitionId},
		{Key: "season._id", Value: seasonId},
	}

	if queryParams.MatchStatus != "" {
		m = append(m, bson.E{Key: "status", Value: queryParams.MatchStatus})
	}

	if queryParams.ClubId != 0 {
		m = append(m, bson.E{Key: "$or", Value: []bson.A{
			{"homeTeam._id", queryParams.ClubId},
			{"awayTeam._id", queryParams.ClubId}},
		})
	}

	q := bson.D{
		{Key: "$match", Value: m},
	}

	s := bson.D{{Key: "$sort", Value: bson.D{{Key: "utcDate", Value: -1}}}}

	g := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "utcDate"},
			{Key: "utcDate", Value: bson.D{{Key: "$push", Value: "$utcDate"}}},
		},
		},
	}

	cursor, err := r.db.Collection("match").Aggregate(ctx, mongo.Pipeline{q, s, g})
	if err != nil {
		return []primitive.DateTime{}, 0, err
	}

	var matchesDates []struct {
		Id    string               `bson:"_id"`
		Dates []primitive.DateTime `bson:"utcDate"`
	}

	err = cursor.All(context.TODO(), &matchesDates)
	if err != nil {
		return []primitive.DateTime{}, 0, err
	}

	totalCount := len(matchesDates[0].Dates)

	return matchesDates[0].Dates, totalCount, nil
}

// func (r *MatchRepository) GetDates(ctx context.Context, competition_id uint, season_id uint, queryParams domain.GetDatesQueryParams) ([]primitive.DateTime, int, error) {
// 	// o := options.Find()
// 	// o.SetSort(bson.D{{"utcDate", -1}})
// 	// o.SetLimit(10).SetSkip(0)

// 	// q := bson.M{
// 	// 	"$and": []bson.M{
// 	// 		{"competition._id": competition_id},
// 	// 		{"season._id": season_id},
// 	// 	},
// 	// }

// 	f := []bson.M{
// 		{"competition._id": competition_id},
// 		{"season._id": season_id},
// 	}

// 	if queryParams.MatchStatus != "" {
// 		f = append(f, bson.M{"status": queryParams.MatchStatus})
// 	}

// 	if queryParams.ClubId != nil {
// 		f = append(f, bson.M{"$or": []bson.M{
// 			{"homeTeam._id": queryParams.ClubId},
// 			{"awayTeam._id": queryParams.ClubId},
// 		}})
// 	}

// 	q := bson.M{
// 		"$and": f,
// 	}

// 	results, err := r.db.Collection("match").Distinct(ctx, "utcDate", q)
// 	if err != nil {
// 		log.Error().Err(err).Msg("")
// 		return []primitive.DateTime{}, 0, err
// 	}

// 	// for _, result := range results {
// 	// 	d := result.(primitive.DateTime)
// 	// 	layout  := "Monday, 2 Jan 2006"
// 	// 	fmt.Print(d.Time().Format(layout))
// 	// }

// 	totalCount := len(results)
// 	dates := make([]primitive.DateTime, len(results))

// 	if totalCount > 0 {
// 		for i, result := range results {
// 			dates[i] = result.(primitive.DateTime)
// 		}

// 		//reverse slice O_o
// 		for i, j := 0, len(dates)-1; i < j; i, j = i+1, j-1 {
// 			dates[i], dates[j] = dates[j], dates[i]
// 		}
// 	}

// 	return dates, totalCount, nil
// }
