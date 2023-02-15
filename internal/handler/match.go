package handler

import (
	"net/http"
	"strconv"
	"strings"
	
	"timeline/internal/domain"
	"timeline/internal/handler/response"
	"timeline/internal/handler/helpers"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (h *Handler) initMatchRoutes(api *gin.RouterGroup) {
	api.GET("/competition/:competition_id/season/:season_id/matches", h.getMatches)
	api.GET("/competition/:competition_id/season/:season_id/dates", h.getDates)
	api.GET("/match/:match_id", h.getMatchById)
}

//--

var matchStatusDefaultSlice = []string{
	"show_all",
	"canceled",
	"finished",
	"in_play",
	"live",
	"paused",
	"postponed",
	"scheduled",
	"suspended",
}

//--

func (h *Handler) getMatches(ctx *gin.Context) {
	competitionId, err := helpers.RequiredIntParam(ctx, "competition_id")
	if err != nil {
		log.Error().Err(err).Msg("")
		response.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())

		return
	}

	seasonId, err := helpers.RequiredIntParam(ctx, "season_id")
	if err != nil {
		log.Error().Err(err).Msg("")
		response.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())

		return
	}

	queryParams := domain.GetMatchesQueryParams{}
	queryParams.Datination = helpers.SetDatination(ctx.Query("min_date"), ctx.Query("max_date"))
	queryParams.ClubId, _ = strconv.Atoi(ctx.Query("club_id"))

	matchStatus := strings.ToLower(ctx.Query("match_status"))
	ok := helpers.InSlice(matchStatus, matchStatusDefaultSlice)
	if ok {
		queryParams.MatchStatus = ctx.Query("match_status")
	}

	res, err := h.s.Match.GetMatches(ctx, competitionId, seasonId, queryParams)
	if err != nil {
		log.Error().Err(err).Msg("")
		response.NewErrorResponse(ctx, http.StatusInternalServerError, domain.MsgInternalServerError)

		return
	}

	ctx.JSON(http.StatusOK, res)
}

//--

func (h *Handler) getDates(ctx *gin.Context) {
	competitionId, err := helpers.RequiredIntParam(ctx, "competition_id")
	if err != nil {
		log.Error().Err(err).Msg("")
		response.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())

		return
	}

	seasonId, err := helpers.RequiredIntParam(ctx, "season_id")
	if err != nil {
		log.Error().Err(err).Msg("")
		response.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())

		return
	}

	queryParams := domain.GetDatesQueryParams{}
	queryParams.ClubId, _ = strconv.Atoi(ctx.Query("club_id"))

	matchStatus := strings.ToLower(ctx.Query("match_status"))
	ok := helpers.InSlice(matchStatus, matchStatusDefaultSlice)
	if ok {
		queryParams.MatchStatus = ctx.Query("match_status")
	}

	res, resTotalCount, err := h.s.Match.GetMatchesDates(ctx, competitionId, seasonId, queryParams)
	if err != nil {
		log.Error().Err(err).Msg("")
		response.NewErrorResponse(ctx, http.StatusInternalServerError, domain.MsgInternalServerError)
		return
	}

	ctx.Header("Access-Control-Expose-Headers", "X-Total-Count")
	ctx.Header("X-Total-Count", strconv.Itoa(resTotalCount))
	ctx.JSON(http.StatusOK, res)
}

//--

func (h *Handler) getMatchById(ctx *gin.Context) {
	matchId, err := helpers.RequiredIntParam(ctx, "match_id")
	if err != nil {
		log.Error().Err(err).Msg("")
		response.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())

		return
	}

	res, err := h.s.Match.GetMatchById(ctx.Request.Context(), matchId)
	if err != nil {
		log.Error().Err(err).Msg("")
		response.NewErrorResponse(ctx, http.StatusInternalServerError, domain.MsgInternalServerError)

		return
	}

	ctx.JSON(http.StatusOK, res)
}
