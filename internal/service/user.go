package service

import (
	"context"
	"time"
	"timeline/config"
	"timeline/internal/domain"
	"timeline/internal/repository"
	"timeline/pkg/cipher"
)

type UserService struct {
	r   repository.UserInterface
	e   EmailInterface
	tkn TokenInterface
	cfg *config.Config
}

func NewUserService(
	r repository.UserInterface,
	email EmailInterface,
	token TokenInterface,
	config *config.Config) *UserService {
	return &UserService{
		r:   r,
		e:   email,
		tkn: token,
		cfg: config,
	}
}

//go:generate mockgen -source=user.go -destination=mock/user.go

type UserInterface interface {
	SignUp(ctx context.Context, input domain.UserSignUpInput) error

	CreateSession(ctx context.Context, user domain.User, rememberMe bool) (domain.Session, error)
	SignInByUsername(ctx context.Context, input domain.UserSignInUsernameInput) (domain.Session, error)
	SignInByEmail(ctx context.Context, input domain.UserSignInEmailInput) (domain.Session, error)

	SignOut(ctx context.Context, refreshToken string) error

	RefreshSession(ctx context.Context, refreshToken string) (domain.Session, error)

	ValidateAccessToken(accessToken string) (string, error)

	VerifyRecoveryCode(ctx context.Context, code string) error
	SetRecoveryCode(ctx context.Context, user domain.User) error
	SetNewPassword(ctx context.Context, password, code string) error

	RecoveryPasswordByUsername(ctx context.Context, username string) error
	RecoveryPasswordByEmail(ctx context.Context, email string) error

	AccountActivateByCode(ctx context.Context, code string) error
	AccountActivateResendCode(ctx context.Context, email string) error
}

// - Sign up -.

func (s *UserService) SignUp(ctx context.Context, input domain.UserSignUpInput) error {
	passwordHash := cipher.GeneratePassword(input.Password, s.cfg.Cipher.Salt)

	code, err := cipher.GenerateCode()
	if err != nil {
		return err
	}

	user := domain.User{
		Username:    input.Username,
		Password:    passwordHash,
		Email:       input.Email,
		DateJoined:  time.Now(),
		LastVisitAt: time.Now(),
		Activation: domain.Activation{
			Code:        code,
			IsActivated: false,
		},
	}

	if err = s.r.CreateUser(ctx, user); err != nil {
		return err
	}

	if err = s.e.AccountActivateEmail(user); err != nil {
		return err
	}

	return err
}

// - Sign in -.

func (s *UserService) CreateSession(ctx context.Context, user domain.User, rememberMe bool) (domain.Session, error) {
	tokens, err := s.tkn.GenerateTokenPair(user.ID)
	if err != nil {
		return domain.Session{}, err
	}

	refreshTokenRecord := domain.RefreshTokenRecord{
		UserID:       user.ID,
		RefreshToken: tokens.RefreshToken,
		RememberMe:   rememberMe,
		CreatedAt:    time.Now(),
	}

	if err = s.tkn.SaveRefreshToken(ctx, refreshTokenRecord); err != nil {
		return domain.Session{}, err
	}

	session := domain.Session{
		SignInDTO: domain.SignInDTO{
			User: domain.UserDTO{
				ID:       user.ID,
				Username: user.Username,
				Email:    user.Email,
			},
			AccessToken: tokens.AccessToken,
		},
		// IsActivated:  user.IsActivated,
		RememberMe:   rememberMe,
		RefreshToken: tokens.RefreshToken,
	}

	return session, err
}

func (s *UserService) SignInByUsername(ctx context.Context, input domain.UserSignInUsernameInput) (domain.Session, error) {
	passwordHash := cipher.GeneratePassword(input.Password, s.cfg.Cipher.Salt)
	user, err := s.r.AuthByUsername(ctx, input.Username, passwordHash)
	if err != nil {
		return domain.Session{}, err
	}

	session, err := s.CreateSession(ctx, user, input.RememberMe)
	if err != nil {
		return domain.Session{}, err
	}

	return session, nil
}

func (s *UserService) SignInByEmail(ctx context.Context, input domain.UserSignInEmailInput) (domain.Session, error) {
	passwordHash := cipher.GeneratePassword(input.Password, s.cfg.Cipher.Salt)
	user, err := s.r.AuthByEmail(ctx, input.Email, passwordHash)
	if err != nil {
		return domain.Session{}, err
	}

	session, err := s.CreateSession(ctx, user, input.RememberMe)
	if err != nil {
		return domain.Session{}, err
	}

	return session, nil
}

// - Sign out -.

func (s *UserService) SignOut(ctx context.Context, refreshToken string) error {
	if err := s.tkn.DeleteRefreshToken(ctx, refreshToken); err != nil {
		return err
	}

	return nil
}

// - Refresh session -.

func (s *UserService) RefreshSession(ctx context.Context, refreshToken string) (domain.Session, error) {
	refreshTokenRecord, err := s.tkn.FindRefreshToken(ctx, refreshToken)
	if err != nil {
		return domain.Session{}, err
	}

	user, err := s.r.FindUserById(ctx, refreshTokenRecord.UserID)
	if err != nil {
		return domain.Session{}, err
	}

	tokens, err := s.tkn.GenerateTokenPair(user.ID)
	if err != nil {
		return domain.Session{}, err
	}

	refreshTokenRecord.RefreshToken = tokens.RefreshToken

	if err = s.tkn.SaveRefreshToken(ctx, refreshTokenRecord); err != nil {
		return domain.Session{}, err
	}

	session := domain.Session{
		SignInDTO: domain.SignInDTO{
			User: domain.UserDTO{
				ID:       user.ID,
				Username: user.Username,
				Email:    user.Email,
			},
			AccessToken: tokens.AccessToken,
		},
		// IsActivated:  user.IsActivated,
		RememberMe:   refreshTokenRecord.RememberMe,
		RefreshToken: tokens.RefreshToken,
	}

	return session, err
}

// - Validate access token -.

func (s *UserService) ValidateAccessToken(accessToken string) (string, error) {
	userId, err := s.tkn.ValidateAccessToken(accessToken)
	if err != nil {
		return "", err
	}

	return userId, err
}

// - Resend activation code -.

func (s *UserService) AccountActivateResendCode(ctx context.Context, email string) error {
	user, err := s.r.FindUserByEmail(ctx, email)
	if err != nil {
		return err
	}

	if user.Activation.IsActivated {
		return domain.ErrAccountAlreadyActivated
	}

	// if !user.Activation.IsActivated {
	// 	log.Error().Err(domain.ErrAccountNotActivated)
	// 	return domain.ErrAccountNotActivated
	// }

	if err = s.e.AccountActivateEmail(user); err != nil {
		return err
	}

	return nil
}

// - Activate account by code -.

func (s *UserService) AccountActivateByCode(ctx context.Context, code string) error {
	if err := s.r.AccountActivateByCode(ctx, code);	err != nil {
		return err
	}

	return nil
}

// - Recover password -.

func (s *UserService) SetRecoveryCode(ctx context.Context, user domain.User) error {
	code, err := cipher.GenerateCode()
	if err != nil {
		return err
	}

	if err = s.r.SetRecoveryCode(ctx, user, code); err != nil {
		return err
	}

	user.Recovery.Code = code
	if err = s.e.PasswordRecoveryEmail(user); err != nil {
		return err
	}

	return nil
}

func (s *UserService) RecoveryPasswordByUsername(ctx context.Context, username string) error {
	user, err := s.r.FindUserByUsername(ctx, username)
	if err != nil {
		return err
	}

	if err = s.SetRecoveryCode(ctx, user); err != nil {
		return err
	}

	return nil
}

func (s *UserService) RecoveryPasswordByEmail(ctx context.Context, email string) error {
	user, err := s.r.FindUserByEmail(ctx, email)
	if err != nil {
		return err
	}

	if err = s.SetRecoveryCode(ctx, user); err != nil {
		return err
	}

	return nil
}

// - Verify recovery code -.

func (s *UserService) VerifyRecoveryCode(ctx context.Context, code string) error {
	if err := s.r.VerifyRecoveryCode(ctx, code); err != nil {
		return err
	}

	return nil
}

// - Set new password -.

func (s *UserService) SetNewPassword(ctx context.Context, password, code string) error {
	passwordHash := cipher.GeneratePassword(password, s.cfg.Cipher.Salt)
	if err := s.r.SetNewPassword(ctx, passwordHash, code); err != nil {
		return err
	}

	return nil
}
