package response

import (
	"net/http"
	"timeline/config"
	"timeline/internal/domain"

	"github.com/gin-gonic/gin"
)

type MsgWrapper struct {
	Message string `json:"message"`
}

func NewErrorResponse(ctx *gin.Context, statusCode int, message string) {
	ctx.AbortWithStatusJSON(statusCode, MsgWrapper{message})
}

func SetCookieWithOk(ctx *gin.Context, cfg *config.Config, session domain.Session) {

	// responseDTO := domain.ResponseDTO{
	// 	User: domain.UserDTO{
	// 		ID:          user.ID,
	// 		Username:    user.Username,
	// 		Email:       user.Email,
	// 		IsActivated: user.IsActivated,
	// 	},
	// 	AccessToken: tokens.AccessToken,
	// }

	var maxAge int
	if session.RememberMe {
		maxAge = int(cfg.RefreshTTL.Seconds())
	} else {
		maxAge = 0
	}

	ctx.SetCookie(
		"refreshToken",
		session.RefreshToken,
		maxAge,
		"/",
		"", // os.Getenv("AUTH_DOMAIN")
		cfg.Secure,
		true,
	)

	ctx.JSON(http.StatusOK, session.SignInDTO)
}
