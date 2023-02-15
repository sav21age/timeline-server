package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username    string             `json:"username" binding:"required"`
	Email       string             `json:"email" binding:"required,email"`
	Password    string             `json:"password" binding:"required, min=6"`
	LastVisitAt time.Time          `json:"lastVisitAt" bson:"lastVisitAt"`
	DateJoined  time.Time          `json:"dateJoined" bson:"dateJoined"`
	// LastVisitAt primitive.DateTime `json:"lastVisitAt"`
	// DateJoined  primitive.DateTime `json:"dateJoined"`
	Activation
	Recovery
	Area AreaUser
}

type Activation struct {
	Code        string `json:"activationCode" bson:"code"`
	IsActivated bool   `json:"isActivated" bson:"isActivated"`
}

type Recovery struct {
	Code string `json:"recoveryCode" bson:"code" binding:"required,alphanum,len=32"`
}

type AreaUser struct {
	ID   uint16 `json:"id" bson:"_id,omitempty"`
	Name string `json:"name" bson:"name"`
}

// - Input -.

type UserSignUpInput struct {
	Username string `json:"username" binding:"required,alphanum"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type UserSignInUsernameInput struct {
	Username   string `json:"usernameOrEmail" binding:"required,alphanum"`
	Password   string `json:"password" binding:"required,min=6"`
	RememberMe bool
}

type UserSignInEmailInput struct {
	Email      string `json:"usernameOrEmail" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=6"`
	RememberMe bool
}

type UsernameOrEmailInput struct {
	Username string `json:"usernameOrEmail" binding:"required,alphanum"`
}

type EmailOrUsernameInput struct {
	Email string `json:"usernameOrEmail" binding:"required,email"`
}

type EmailInput struct {
	Email string `json:"email" binding:"required,email"`
}

type NewPasswordInput struct {
	Password string `json:"password" binding:"required,min=6"`
	Code     string `json:"code" binding:"required,alphanum,len=32"`
}

// - DTO -.

type UserDTO struct {
	ID       primitive.ObjectID `json:"id"`
	Username string             `json:"username"`
	Email    string             `json:"email"`
}

type SignInDTO struct {
	User        UserDTO `json:"user"`
	AccessToken string  `json:"accessToken"`
}

// - Return -.

type Session struct {
	SignInDTO SignInDTO
	// IsActivated  bool   `json:"isActivated"`
	RememberMe   bool   `json:"-"`
	RefreshToken string `json:"-"`
}
