package repository

import (
	"context"
	"errors"

	"timeline/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	db *mongo.Database
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{db: db}
}

type UserInterface interface {
	CreateUser(ctx context.Context, user domain.User) error

	FindUserById(ctx context.Context, id primitive.ObjectID) (domain.User, error)
	FindUserByUsername(ctx context.Context, username string) (domain.User, error)
	FindUserByEmail(ctx context.Context, email string) (domain.User, error)

	AuthByUsername(ctx context.Context, username, password string) (domain.User, error)
	AuthByEmail(ctx context.Context, email, password string) (domain.User, error)

	VerifyRecoveryCode(ctx context.Context, code string) error
	SetRecoveryCode(ctx context.Context, user domain.User, code string) error

	SetNewPassword(ctx context.Context, passwordHash, code string) error

	AccountActivateByCode(ctx context.Context, code string) error
}

// - Create user -.

func (r *UserRepository) CreateUser(ctx context.Context, user domain.User) error {
	_, err := r.db.Collection("user").InsertOne(ctx, user)

	if mongo.IsDuplicateKeyError(err) {
		return domain.ErrUserAlreadyExists
	}

	if err != nil {
		return err
	}

	return nil
}

// - Find user by id -.

func (r *UserRepository) FindUserById(ctx context.Context, id primitive.ObjectID) (domain.User, error) {
	var user domain.User
	err := r.db.Collection("user").FindOne(ctx, bson.M{"_id": id}).Decode(&user)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return domain.User{}, domain.ErrUserIdNotExists
	}
	
	if !user.Activation.IsActivated {
		return domain.User{}, domain.ErrAccountNotActivated
	}

	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

// - Auth by username -.

func (r *UserRepository) AuthByUsername(ctx context.Context, username, password string) (domain.User, error) {
	var user domain.User

	// q := bson.M{
	// 	"$and": []bson.M{
	// 		{"username": username},
	// 		{"password": password},
	// 	},
	// }
	err := r.db.Collection("user").FindOne(ctx, bson.M{"username": username, "password": password}).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return domain.User{}, domain.ErrUsernameOrPasswordInvalid
	}

	if !user.Activation.IsActivated {
		return domain.User{}, domain.ErrAccountNotActivated
	}

	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

// - Auth by email -.

func (r *UserRepository) AuthByEmail(ctx context.Context, email, password string) (domain.User, error) {
	var user domain.User

	q := bson.M{
		"$and": []bson.M{
			{"email": email},
			{"password": password},
		},
	}

	err := r.db.Collection("user").FindOne(ctx, q).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return domain.User{}, domain.ErrEmailOrPasswordInvalid
	}

	if !user.Activation.IsActivated {
		return domain.User{}, domain.ErrAccountNotActivated
	}

	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

// - Find user by username -.

func (r *UserRepository) FindUserByUsername(ctx context.Context, username string) (domain.User, error) {
	var user domain.User

	err := r.db.Collection("user").FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return domain.User{}, domain.ErrUsernameOrPasswordInvalid
	}

	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

// - Find user by email -.

func (r *UserRepository) FindUserByEmail(ctx context.Context, email string) (domain.User, error) {
	var user domain.User

	q := bson.M{
		"$and": []bson.M{
			{"email": email},
		},
	}

	err := r.db.Collection("user").FindOne(ctx, q).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return domain.User{}, domain.ErrEmailNotExists
	}

	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

// - Verify recovery code -.

func (r *UserRepository) VerifyRecoveryCode(ctx context.Context, code string) error {
	var user domain.User

	err := r.db.Collection("user").FindOne(ctx, bson.M{"recovery.code": code}).Decode(&user)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return domain.ErrCodeInvalid
	}

	if err != nil {
		return err
	}

	return nil
}

// - Set new password -.

func (r *UserRepository) SetNewPassword(ctx context.Context, passwordHash, code string) error {
	cursor, err := r.db.Collection("user").UpdateOne(ctx,
		bson.M{"recovery.code": code},
		bson.M{"$set": bson.M{"password": passwordHash, "recovery.code": ""}})
	
	if err != nil {
		return err
	}

	if cursor.ModifiedCount == 0 {
		return domain.ErrCodeInvalid
	}

	return nil
}

// - Set recovery code -.

func (r *UserRepository) SetRecoveryCode(ctx context.Context, user domain.User, code string) error {
	cursor, err := r.db.Collection("user").UpdateOne(ctx,
		bson.M{"_id": user.ID},
		bson.M{"$set": bson.M{"recovery.code": code}})
	
	if err != nil {
		return err
	}

	if cursor.ModifiedCount == 0 {
		return domain.ErrCodeInvalid
	}

	return nil
}

// - Activate account by code -.

func (r *UserRepository) AccountActivateByCode(ctx context.Context, code string) error {
	cursor, err := r.db.Collection("user").UpdateOne(ctx,
		bson.M{"activation.code": code},
		bson.M{"$set": bson.M{"activation.isActivated": true, "activation.code": ""}})
	
	if err != nil {
		return err
	}

	if cursor.ModifiedCount == 0 {
		return domain.ErrCodeInvalid
	}

	return nil
}

// func (r *UserRepository) ResendEmail(ctx context.Context, email string) (domain.User, error) {
// 	var user domain.User

// 	err := r.db.Collection("user").FindOne(ctx, bson.M{"email": email}).Decode(&user)
// 	if errors.Is(err, mongo.ErrNoDocuments){
// 		return domain.User{}, domain.ErrEmailNotExists
// 	}

// 	if err != nil {
// 		return domain.User{}, err
// 	}

// 	if user.Activation.IsActivated {
// 		return domain.User{}, domain.ErrAccountAlreadyActivated
// 	}

// 	return user, nil
// }
