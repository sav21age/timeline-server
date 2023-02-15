package service

import (
	"errors"
	"context"
	"fmt"
	"time"

	"timeline/config"
	"timeline/internal/domain"
	"timeline/internal/repository"
	"timeline/pkg/cipher"

	"github.com/golang-jwt/jwt"
	// "github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TokenService struct {
	r   repository.TokenInterface
	cfg *config.Config
}

func NewTokenService(r repository.TokenInterface, config *config.Config) *TokenService {
	return &TokenService{
		r:   r,
		cfg: config,
	}
}

//go:generate mockgen -source=token.go -destination=mock/token.go

type TokenInterface interface {
	GenerateTokenPair(id primitive.ObjectID) (domain.TokenPair, error)
	ValidateAccessToken(token string) (string, error)

	SaveRefreshToken(ctx context.Context, refreshTokenRecord domain.RefreshTokenRecord) error
	DeleteRefreshToken(ctx context.Context, token string) error
	FindRefreshToken(ctx context.Context, token string) (domain.RefreshTokenRecord, error)
}

// - Generate tokens -.

func (s *TokenService) GenerateTokenPair(id primitive.ObjectID) (domain.TokenPair, error) {

	accessTTL, err := time.ParseDuration(s.cfg.AccessTTL.String())
	if err != nil {
		return domain.TokenPair{}, err
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(accessTTL).Unix(),
		Id:        id.Hex(),
	})

	accessToken, err := at.SignedString([]byte(s.cfg.AccessSecret))
	if err != nil {
		return domain.TokenPair{}, err
	}

	// ---

	refreshToken, err := cipher.GenerateCode()
	if err != nil {
		return domain.TokenPair{}, err
	}

	return domain.TokenPair{AccessToken: accessToken, RefreshToken: refreshToken}, err
}

// - Validate access token -.

func (s *TokenService) ValidateAccessToken(token string) (string, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(s.cfg.AccessSecret), nil
	})

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok || !t.Valid {
		return "", errors.New("invalid token")
	}

	return claims["jti"].(string), err
}

// - Save refresh token -.

func (s *TokenService) SaveRefreshToken(ctx context.Context, refreshToken domain.RefreshTokenRecord) error {
	return s.r.SaveRefreshToken(ctx, refreshToken)
}

// - Find refresh token -.

func (s *TokenService) FindRefreshToken(ctx context.Context, token string) (domain.RefreshTokenRecord, error) {
	return s.r.FindRefreshToken(ctx, token)
}

// - Delete refresh token -.

func (s *TokenService) DeleteRefreshToken(ctx context.Context, token string) error {
	return s.r.DeleteRefreshToken(ctx, token)
}
