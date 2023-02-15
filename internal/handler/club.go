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

func (h *Handler) initClubRoutes(api *gin.RouterGroup) {
	api.GET("/clubs", h.getClubs)
	api.GET("/season/:season_id/clubs", h.getClubsBySeasonId)
	api.GET("/clubs/area", h.GetClubsAreas)

	// api.GET("/club/:club_id", h.getClubById)
	isAuth := api.Group("/", h.userIdentity)
	{
		isAuth.GET("/club/:club_id", h.getClubById)
	}
}

//--

func (h *Handler) getClubs(ctx *gin.Context) {
	sortByDefaultSlice := []string{"country", "club"}
	perPageDefaultSlice := []int{6, 12, 18, 24}
	perPageDefault := perPageDefaultSlice[1]

	queryParams := domain.GetClubsQueryParams{}
	var err error
	queryParams.Pagination, err = helpers.SetPagination(ctx.Query("page"), ctx.Query("per_page"), perPageDefaultSlice, perPageDefault)
	if err != nil {
		log.Error().Err(err).Msg("")
		response.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		
		return
	}
	
	queryParams.AreaId, _ = strconv.Atoi(ctx.Query("area_id"))

	sortBy := strings.ToLower(ctx.Query("sort_by"))
	ok := helpers.InSlice(sortBy, sortByDefaultSlice)
	if ok {
		queryParams.SortBy = sortBy
	}

	res, resTotalCount, err := h.s.Club.GetClubs(ctx, queryParams)
	if err != nil {
		log.Error().Err(err).Msg("")
		response.NewErrorResponse(ctx, http.StatusInternalServerError, domain.MsgInternalServerError)

		return
	}

	// !reflect.DeepEqual(seasons, []domain.Season{}
	if resTotalCount != 0 {
		ctx.Header("Access-Control-Expose-Headers", "X-Total-Count")
		ctx.Header("X-Total-Count", strconv.Itoa(resTotalCount))
		ctx.JSON(http.StatusOK, res)
	} else {
		ctx.JSON(http.StatusNoContent, nil)
	}
}

//--

func (h *Handler) getClubsBySeasonId(ctx *gin.Context) {
	seasonId, err := helpers.RequiredIntParam(ctx, "season_id")
	if err != nil {
		log.Error().Err(err).Msg("")
		response.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())

		return
	}

	res, err := h.s.Club.GetClubsBySeasonId(ctx, seasonId)
	if err != nil {
		log.Error().Err(err).Msg("")
		response.NewErrorResponse(ctx, http.StatusInternalServerError, domain.MsgInternalServerError)

		return
	}

	ctx.JSON(http.StatusOK, res)
}

//--

func (h *Handler) getClubById(ctx *gin.Context) {
	clubId, err := helpers.RequiredIntParam(ctx, "club_id")
	if err != nil {
		log.Error().Err(err).Msg("")
		response.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())

		return
	}

	res, err := h.s.Club.GetClubById(ctx.Request.Context(), clubId)
	if err != nil {
		log.Error().Err(err).Msg("")
		response.NewErrorResponse(ctx, http.StatusInternalServerError, domain.MsgInternalServerError)

		return
	}

	ctx.JSON(http.StatusOK, res)
}

//--

func (h *Handler) GetClubsAreas(ctx *gin.Context) {
	res, err := h.s.Club.GetClubsAreas(ctx)
	if err != nil {
		log.Error().Err(err).Msg("")
		response.NewErrorResponse(ctx, http.StatusInternalServerError, domain.MsgInternalServerError)

		return
	}

	ctx.JSON(http.StatusOK, res)
}
