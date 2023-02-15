package domain

// import "github.com/pkg/errors"
import "errors"

var (
	ErrUserAlreadyExists         = errors.New("user already exists")
	ErrUserIdNotExists           = errors.New("user with such id not exists")
	ErrEmailNotExists            = errors.New("user with such email not exists")
	ErrUsernameOrPasswordInvalid = errors.New("username or password is invalid")
	ErrEmailOrPasswordInvalid    = errors.New("email or password is invalid")
	ErrUsernameOrEmailInvalid    = errors.New("username or email is invalid")
	ErrAccountNotActivated       = errors.New("account not activated")
	ErrAccountAlreadyActivated   = errors.New("account already activated")
	ErrCodeInvalid               = errors.New("code is invalid")
	ErrTokenInvalid              = errors.New("token is invalid")

	MsgQueryParamInvalid   = "invalid query parameter"
	MsgInternalServerError = "internal server error"
	MsgBadRequest          = "invalid input body"
)
