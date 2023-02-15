package handler

import (
	"net/http"
	"strconv"
	
	"timeline/internal/domain"
	"timeline/internal/handler/response"
	"timeline/internal/handler/helpers"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (h *Handler) initSeasonRoutes(api *gin.RouterGroup) {
	api.GET("/competition/:competition_id/seasons", h.getSeasons)
	api.GET("/season/:season_id", h.getSeasonById)
	api.GET("/season/match/:match_id", h.getSeasonByMatchId)
}

//--

func (h *Handler) getSeasons(ctx *gin.Context) {
	competitionId, err := helpers.RequiredIntParam(ctx, "competition_id")
	if err != nil {
		log.Error().Err(err).Msg("")
		response.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())

		return
	}

	queryParams := domain.GetSeasonsQueryParams{}
	perPageDefaultSlice := []int{9}
	perPageDefault := perPageDefaultSlice[0]

	queryParams.Pagination, err = helpers.SetPagination(ctx.Query("page"), ctx.Query("per_page"), perPageDefaultSlice, perPageDefault)
	if err != nil {
		log.Error().Err(err).Msg("")
		response.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	res, resTotalCount, err := h.s.Season.GetSeasons(ctx, competitionId, queryParams)
	if err != nil {
		log.Error().Err(err).Msg("")
		response.NewErrorResponse(ctx, http.StatusInternalServerError, domain.MsgInternalServerError)
		return
	}

	if resTotalCount != 0 {
		ctx.Header("Access-Control-Expose-Headers", "X-Total-Count")
		ctx.Header("X-Total-Count", strconv.Itoa(resTotalCount))
		ctx.JSON(http.StatusOK, res)
	} else {
		ctx.JSON(http.StatusNoContent, nil)
	}
}

//--

func (h *Handler) getSeasonById(ctx *gin.Context) {
	seasonId, err := helpers.RequiredIntParam(ctx, "season_id")
	if err != nil {
		log.Error().Err(err).Msg("")
		response.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())

		return
	}

	res, err := h.s.Season.GetSeasonById(ctx.Request.Context(), seasonId)
	if err != nil {
		log.Error().Err(err).Msg("")
		response.NewErrorResponse(ctx, http.StatusInternalServerError, domain.MsgInternalServerError)

		return
	}

	ctx.JSON(http.StatusOK, res)
}

//--

func (h *Handler) getSeasonByMatchId(ctx *gin.Context) {
	matchId, err := helpers.RequiredIntParam(ctx, "match_id")
	if err != nil {
		log.Error().Err(err).Msg("")
		response.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())

		return
	}

	res, err := h.s.Season.GetSeasonByMatchId(ctx, matchId)
	if err != nil {
		log.Error().Err(err).Msg("")
		response.NewErrorResponse(ctx, http.StatusInternalServerError, domain.MsgInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, res)
}
