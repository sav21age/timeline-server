package repository

import (
	"context"
	"errors"
	"timeline/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TokenRepository struct {
	db *mongo.Database
}

func NewTokenRepository(db *mongo.Database) *TokenRepository {
	return &TokenRepository{
		db: db,
	}
}

type TokenInterface interface {
	SaveRefreshToken(ctx context.Context, refreshTokenRecord domain.RefreshTokenRecord) error
	DeleteRefreshToken(ctx context.Context, token string) error
	FindRefreshToken(ctx context.Context, token string) (domain.RefreshTokenRecord, error)
}

// - Save refresh token -.

func (r *TokenRepository) SaveRefreshToken(ctx context.Context, refreshToken domain.RefreshTokenRecord) error {
	// o := options.Update().SetUpsert(true)
	// f := bson.M{"_id": userId}
	// q := bson.M{"$set": bson.M{"refreshToken": refreshToken, "createdAt": time.Now()}}

	o := options.Update().SetUpsert(true)
	f := bson.M{"_id": refreshToken.UserID}
	q := bson.M{"$set": refreshToken}

	// cursor, err := r.db.Collection("token").UpdateOne(ctx, f, q, o)
	_, err := r.db.Collection("token").UpdateOne(ctx, f, q, o)
	if err != nil {
		return err
	}

	// if cursor.ModifiedCount == 0 {
	// 	return domain.ErrTokenInvalid
	// }

	return nil
}

// - Delete refresh token -.

func (r *TokenRepository) DeleteRefreshToken(ctx context.Context, refreshToken string) error {
	cursor, err := r.db.Collection("token").DeleteOne(ctx, bson.M{"refreshToken": refreshToken})
	if err != nil {
		return err
	}

	if cursor.DeletedCount == 0 {
		return domain.ErrTokenInvalid
	}

	return nil
}

// - Find refresh token -.

func (r *TokenRepository) FindRefreshToken(ctx context.Context, refreshToken string) (domain.RefreshTokenRecord, error) {
	var token domain.RefreshTokenRecord
	err := r.db.Collection("token").FindOne(ctx, bson.M{"refreshToken": refreshToken}).Decode(&token)
	
	if errors.Is(err, mongo.ErrNoDocuments) {
		return domain.RefreshTokenRecord{}, domain.ErrTokenInvalid
	}

	if err != nil {
		return domain.RefreshTokenRecord{}, err
	}

	return token, nil
}
