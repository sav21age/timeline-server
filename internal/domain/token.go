package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RefreshTokenRecord struct {
	UserID       primitive.ObjectID `json:"userId" bson:"_id"`
	RefreshToken string             `json:"refreshToken" bson:"refreshToken"`
	RememberMe   bool               `json:"rememberMe" bson:"rememberMe"`
	CreatedAt    time.Time          `json:"createdAt" bson:"createdAt"`
}

type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}