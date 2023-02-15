package handler

import (
	"errors"
	"net/http"
	"regexp"
	"strings"
	"timeline/internal/domain"
	"timeline/internal/handler/response"
	"timeline/pkg/email"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/rs/zerolog/log"
)

func (h *Handler) initUserRoutes(api *gin.RouterGroup) {
	auth := api.Group("/auth")
	{
		auth.POST("/sign-in", h.signIn)
		auth.POST("/sign-up", h.signUp)

		auth.GET("/refresh", h.refreshTokenPair)

		account := auth.Group("/account")
		{
			account.GET("/activate/:code", h.accountActivateByCode)
			account.POST("/activate/resend", h.accountActivateResendCode)
		}

		password := auth.Group("/password")
		{
			password.POST("/recovery", h.passwordRecovery)
			password.GET("/reset/:code", h.passwordRecoveryCodeVerify)
			password.POST("/reset", h.setNewPassword)
		}

		auth.GET("/sign-out", h.signOut)
		// isAuth := auth.Group("/", h.userIdentity)
		// {
		// isAuth.GET("/sign-out", h.signOut)
		// }
	}
}

// - Sign up -.

func (h *Handler) signUp(ctx *gin.Context) {
	var input domain.UserSignUpInput

	if err := ctx.BindJSON(&input); err != nil {
		log.Error().Err(err).Msg("")
		response.NewErrorResponse(ctx, http.StatusBadRequest, domain.MsgBadRequest)

		return
	}

	err := h.s.User.SignUp(ctx.Request.Context(), input)
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			log.Error().Err(err).Msg("")
			response.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())

			return
		}

		if !errors.Is(err, email.ErrMailServiceUnavailable) {
			log.Error().Err(err).Msg("")
			response.NewErrorResponse(ctx, http.StatusInternalServerError, domain.MsgInternalServerError)

			return
		}
	}

	ctx.JSON(http.StatusCreated, response.MsgWrapper{Message: "account created"})
}

// - Sign in -.

func (h *Handler) signInByUsername(ctx *gin.Context) (domain.Session, error) {
	var input domain.UserSignInUsernameInput
	err := ctx.ShouldBindBodyWith(&input, binding.JSON)
	if err != nil {
		log.Error().Err(err).Msg("")
		return domain.Session{}, errors.New(domain.MsgBadRequest)
	}

	return h.s.User.SignInByUsername(ctx.Request.Context(), input)
}

func (h *Handler) signInByEmail(ctx *gin.Context) (domain.Session, error) {
	var input domain.UserSignInEmailInput
	err := ctx.ShouldBindBodyWith(&input, binding.JSON)
	if err != nil {
		log.Error().Err(err).Msg("")
		return domain.Session{}, errors.New(domain.MsgBadRequest)
	}

	return h.s.User.SignInByEmail(ctx.Request.Context(), input)
}

func (h *Handler) signIn(ctx *gin.Context) {
	session, err := h.signInByEmail(ctx)

	if err != nil {
		if errors.Is(err, domain.ErrEmailOrPasswordInvalid) {
			log.Error().Err(err).Msg("")
			response.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())

			return
		}

		if errors.Is(err, domain.ErrAccountNotActivated) {
			log.Error().Err(err).Msg("")
			response.NewErrorResponse(ctx, http.StatusLocked, err.Error())

			return
		}

		session, err = h.signInByUsername(ctx)

		if err != nil {
			if errors.Is(err, domain.ErrUsernameOrPasswordInvalid) {
				log.Error().Err(err).Msg("")
				response.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())

				return
			}

			if errors.Is(err, domain.ErrAccountNotActivated) {
				log.Error().Err(err).Msg("")
				response.NewErrorResponse(ctx, http.StatusLocked, err.Error())

				return
			}

			log.Error().Err(err).Msg("")
			response.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())

			return
		}
	}

	response.SetCookieWithOk(ctx, h.cfg, session)
}

// - Sign out -.

func (h *Handler) signOut(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie("refreshToken")
	if err != nil {
		log.Error().Err(err).Msg("")
		response.NewErrorResponse(ctx, http.StatusBadRequest, domain.MsgBadRequest)

		return
	}

	err = h.s.User.SignOut(ctx.Request.Context(), refreshToken)
	if err != nil {
		log.Error().Err(err).Msg("")
		response.NewErrorResponse(ctx, http.StatusInternalServerError, domain.MsgInternalServerError)

		return
	}

	ctx.SetCookie("refreshToken", "", -1, "/", "", true, true)
	ctx.AbortWithStatus(204)
}

// - Refresh token -.

func (h *Handler) refreshTokenPair(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie("refreshToken")
	if err != nil {
		log.Error().Err(err).Msg("")
		response.NewErrorResponse(ctx, http.StatusBadRequest, domain.MsgBadRequest)

		return
	}

	signInDTO, err := h.s.User.RefreshSession(ctx.Request.Context(), refreshToken)
	if err != nil {
		if errors.Is(err, domain.ErrTokenInvalid) {
			// response.ErrBadReq(ctx, err.Error())
			log.Error().Err(err).Msg("")
			response.NewErrorResponse(ctx, http.StatusUnauthorized, err.Error())

			return
		}

		if errors.Is(err, domain.ErrUserIdNotExists) {
			log.Error().Err(err).Msg("")
			response.NewErrorResponse(ctx, http.StatusUnauthorized, err.Error())

			return
		}

		if errors.Is(err, domain.ErrAccountNotActivated) {
			log.Error().Err(err).Msg("")
			response.NewErrorResponse(ctx, http.StatusLocked, err.Error())

			return
		}

		log.Error().Err(err).Msg("")
		response.NewErrorResponse(ctx, http.StatusInternalServerError, domain.MsgInternalServerError)
		return
	}

	response.SetCookieWithOk(ctx, h.cfg, signInDTO)
}

// - Resend activation code -.

func (h *Handler) accountActivateResendCode(ctx *gin.Context) {
	var input domain.EmailInput

	if err := ctx.BindJSON(&input); err != nil {
		log.Error().Err(err).Msg("")
		response.NewErrorResponse(ctx, http.StatusBadRequest, domain.MsgBadRequest)

		return
	}

	err := h.s.User.AccountActivateResendCode(ctx.Request.Context(), input.Email)
	if err != nil {
		if errors.Is(err, domain.ErrEmailNotExists) {
			log.Error().Err(err).Msg("")
			response.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())

			return
		}
		if errors.Is(err, domain.ErrAccountAlreadyActivated) {
			log.Error().Err(err).Msg("")
			response.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())

			return
		}

		if !errors.Is(err, email.ErrMailServiceUnavailable) {
			log.Error().Err(err).Msg("")
			response.NewErrorResponse(ctx, http.StatusInternalServerError, domain.MsgInternalServerError)

			return
		}
	}

	ctx.JSON(http.StatusOK, response.MsgWrapper{Message: "email has been sent"})
}

// - Activate account by code -.

func (h *Handler) accountActivateByCode(ctx *gin.Context) {
	code := ctx.Param("code")

	var alphabet = regexp.MustCompile(`^[a-z0-9]+$`).MatchString
	if len(code) != 32 && !alphabet(code) {
		err := domain.ErrCodeInvalid
		log.Error().Err(err).Msg("")
		response.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())

		return
	}

	if err := h.s.User.AccountActivateByCode(ctx.Request.Context(), code); err != nil {
		if errors.Is(err, domain.ErrCodeInvalid) {
			response.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())

			return
		}

		response.NewErrorResponse(ctx, http.StatusInternalServerError, domain.MsgInternalServerError)

		return
	}

	ctx.JSON(http.StatusOK, true)
}

// - Recover password -.

func passwordRecoveryByUsername(h *Handler, ctx *gin.Context) error {
	var input domain.UsernameOrEmailInput
	err := ctx.ShouldBindBodyWith(&input, binding.JSON)
	if err != nil {
		log.Error().Err(err).Msg("")

		return err
	}

	return h.s.User.RecoveryPasswordByUsername(ctx.Request.Context(), input.Username)
}

func passwordRecoveryByEmail(h *Handler, ctx *gin.Context) error {
	var input domain.EmailOrUsernameInput
	err := ctx.ShouldBindBodyWith(&input, binding.JSON)
	if err != nil {
		log.Error().Err(err).Msg("")

		return err
	}

	return h.s.User.RecoveryPasswordByEmail(ctx.Request.Context(), input.Email)
}

func (h *Handler) passwordRecovery(ctx *gin.Context) {
	err := passwordRecoveryByEmail(h, ctx)
	if err != nil {
		if !errors.Is(err, email.ErrMailServiceUnavailable) {

			if errors.Is(err, domain.ErrUsernameOrEmailInvalid) {
				log.Error().Err(err).Msg("")
				response.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())

				return
			}

			err = passwordRecoveryByUsername(h, ctx)

			if err != nil {
				if errors.Is(err, domain.ErrUsernameOrEmailInvalid) {
					log.Error().Err(err).Msg("")
					response.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())

					return
				}

				if !errors.Is(err, email.ErrMailServiceUnavailable) {
					log.Error().Err(err).Msg("")
					response.NewErrorResponse(ctx, http.StatusInternalServerError, domain.MsgInternalServerError)

					return
				}
			}
		}
	}

	ctx.JSON(http.StatusOK, true)
}

// - Reset password -.

func (h *Handler) passwordRecoveryCodeVerify(ctx *gin.Context) {
	code := ctx.Param("code")

	// var alphabet = regexp.MustCompile(`^[a-z0-9]+$`).MatchString
	// if len(code) != 32 && !alphabet(code) {
	// 	msg := "code is invalid"
	// 	response.ErrBadReq(ctx, msg)
	// 	log.Error().Msg(msg)
	// 	return
	// }

	var alphabet = regexp.MustCompile(`^[a-z0-9]+$`).MatchString
	if len(code) != 32 && !alphabet(code) {
		response.NewErrorResponse(ctx, http.StatusBadRequest, domain.ErrCodeInvalid.Error())
		log.Error().Msg(domain.ErrCodeInvalid.Error())

		return
	}

	if err := h.s.User.VerifyRecoveryCode(ctx.Request.Context(), code); err != nil {
		if errors.Is(err, domain.ErrCodeInvalid) {
			response.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())

			return
		}

		response.NewErrorResponse(ctx, http.StatusInternalServerError, domain.MsgInternalServerError)

		return
	}

	ctx.JSON(http.StatusOK, true)
}

func (h *Handler) setNewPassword(ctx *gin.Context) {
	var input domain.NewPasswordInput

	if err := ctx.BindJSON(&input); err != nil {
		response.NewErrorResponse(ctx, http.StatusBadRequest, domain.MsgBadRequest)

		return
	}

	if err := h.s.User.SetNewPassword(ctx.Request.Context(), input.Password, input.Code); err != nil {
		if errors.Is(err, domain.ErrCodeInvalid) {
			response.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())

			return
		}

		response.NewErrorResponse(ctx, http.StatusInternalServerError, domain.MsgInternalServerError)

		return
	}

	ctx.JSON(http.StatusOK, true)
}

// - User identity -.

func (h *Handler) userIdentity(ctx *gin.Context) {
	id, err := h.parseAuthHeader(ctx)
	if err != nil {
		// newResponse(c, http.StatusUnauthorized, err.Error())
		log.Error().Err(err).Msg("")
		response.NewErrorResponse(ctx, http.StatusUnauthorized, err.Error())
	}

	// c.Set(userCtx, id)
	ctx.Set("userId", id)
}

func (h *Handler) parseAuthHeader(ctx *gin.Context) (string, error) {
	header := ctx.GetHeader("Authorization")
	if header == "" {
		return "", errors.New("empty auth header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("invalid auth header")
	}

	if len(headerParts[1]) == 0 {
		return "", errors.New("token is empty")
	}

	return h.s.User.ValidateAccessToken(headerParts[1])
}
