package handler

import (
	"net/http"
	"timeline/internal/domain"
	"timeline/internal/handler/response"
	"timeline/internal/handler/helpers"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (h *Handler) initCompetitionRoutes(api *gin.RouterGroup) {
	api.GET("/competitions", h.getCompetitions)
	api.GET("/competition/:competition_id", h.getCompetitionById)
}

//--

func (h *Handler) getCompetitions(ctx *gin.Context) {
	res, err := h.s.Competition.GetCompetitions(ctx)
	if err != nil {
		log.Error().Err(err).Msg("")
		response.NewErrorResponse(ctx, http.StatusInternalServerError, domain.MsgInternalServerError)

		return
	}

	ctx.JSON(http.StatusOK, res)
}

//--

func (h *Handler) getCompetitionById(ctx *gin.Context) {
	competitionId, err := helpers.RequiredIntParam(ctx, "competition_id")
	if err != nil {
		log.Error().Err(err).Msg("")
		response.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())

		return
	}

	res, err := h.s.Competition.GetCompetitionById(ctx.Request.Context(), competitionId)
	if err != nil {
		log.Error().Err(err).Msg("")
		response.NewErrorResponse(ctx, http.StatusInternalServerError, domain.MsgInternalServerError)

		return
	}

	ctx.JSON(http.StatusOK, res)
}
