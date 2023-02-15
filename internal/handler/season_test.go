package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"timeline/config"

	"timeline/internal/domain"
	"timeline/internal/service"
	"timeline/internal/service/mock"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"

	// "github.com/magiconair/properties/assert"
	"github.com/go-playground/assert/v2"
)

func TestHandler_getSeasons(t *testing.T) {
	type mockBehavior func(r *mock_service.MockSeasonInterface, competitionId int, queryParams domain.GetSeasonsQueryParams)

	competitionId := 2021

	url := fmt.Sprintf("/competition/%d/seasons", competitionId)

	sd, err := time.Parse(time.RFC3339, "2022-08-05T00:00:00Z")
	if err != nil {
		fmt.Println(err)
		return
	}

	ed, err := time.Parse(time.RFC3339, "2023-05-28T00:00:00Z")
	if err != nil {
		fmt.Println(err)
		return
	}

	resStruct := []domain.Season{{
		ID:              uint16(1490),
		StartDate:       sd,
		EndDate:         ed,
		CurrentMatchday: uint8(7),
		Winner:          nil,
		Competition: domain.CompetitionSeason{
			ID:     uint16(2021),
			Name:   "Premier League",
			Code:   "PL",
			Emblem: "PL.png",
			Area: domain.AreaSeason{
				ID:          uint16(2072),
				Name:        "England",
				CountryCode: "ENG",
				Ensign:      "770.svg",
			},
		},
		Available: true,
	}}

	resJson, err := json.Marshal(resStruct)
	if err != nil {
		fmt.Println(err)
		return
	}

	tests := []struct {
		name                 string
		url                  string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
		queryParams          domain.GetSeasonsQueryParams
		competitionId        int
	}{
		{
			name: "500 empty",
			url:  url,
			queryParams: domain.GetSeasonsQueryParams{
				Pagination: domain.Pagination{Skip: 0, Limit: 0},
			},
			competitionId: competitionId,
			mockBehavior: func(r *mock_service.MockSeasonInterface, competitionId int, queryParams domain.GetSeasonsQueryParams) {
				r.EXPECT().
					GetSeasons(gomock.Any(), competitionId, queryParams).
					Return([]domain.Season{}, 0, nil).
					AnyTimes()
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"message":"invalid query parameter: page"}`,
		},
		{
			name: "400 competitionId:EPL!page=1",
			url:  fmt.Sprintf("/competition/%s/seasons?page=1", "EPL"),
			queryParams: domain.GetSeasonsQueryParams{
				Pagination: domain.Pagination{Skip: 0, Limit: 9},
			},
			competitionId: competitionId,
			mockBehavior: func(r *mock_service.MockSeasonInterface, competitionId int, queryParams domain.GetSeasonsQueryParams) {
				r.EXPECT().
					GetSeasons(gomock.Any(), "EPL", queryParams).
					Return([]domain.Season{}, 0, nil).
					AnyTimes()
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"invalid query parameter: competition_id"}`,
		},
		{
			name: "204 competitionId:10000!page=1",
			url:  fmt.Sprintf("/competition/%d/seasons?page=1", 10000),
			queryParams: domain.GetSeasonsQueryParams{
				Pagination: domain.Pagination{Skip: 0, Limit: 9},
			},
			competitionId: competitionId,
			mockBehavior: func(r *mock_service.MockSeasonInterface, competitionId int, queryParams domain.GetSeasonsQueryParams) {
				r.EXPECT().
					GetSeasons(gomock.Any(), 10000, queryParams).
					Return([]domain.Season{}, 0, nil).
					AnyTimes()
			},
			expectedStatusCode:   http.StatusNoContent,
			expectedResponseBody: "",
		},
		{
			name: "204 page=1000",
			url:  fmt.Sprintf("%s?page=1000", url),
			queryParams: domain.GetSeasonsQueryParams{
				Pagination: domain.Pagination{Skip: 8991, Limit: 9},
			},
			competitionId: competitionId,
			mockBehavior: func(r *mock_service.MockSeasonInterface, competitionId int, queryParams domain.GetSeasonsQueryParams) {
				r.EXPECT().
					GetSeasons(gomock.Any(), competitionId, queryParams).
					Return([]domain.Season{}, 0, nil).
					AnyTimes()
			},
			expectedStatusCode:   http.StatusNoContent,
			expectedResponseBody: "",
		},
		{
			name: "200 page=1",
			url:  fmt.Sprintf("%s?page=1", url),
			queryParams: domain.GetSeasonsQueryParams{
				Pagination: domain.Pagination{Skip: 0, Limit: 9},
			},
			competitionId: competitionId,
			mockBehavior: func(r *mock_service.MockSeasonInterface, competitionId int, queryParams domain.GetSeasonsQueryParams) {
				r.EXPECT().
					GetSeasons(gomock.Any(), competitionId, queryParams).
					Return(resStruct, 321, nil).
					AnyTimes()
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: string(resJson),
		},
		{
			name: "200 page=1&per_page=9",
			url:  fmt.Sprintf("%s?page=1&per_page=9", url),
			queryParams: domain.GetSeasonsQueryParams{
				Pagination: domain.Pagination{Skip: 0, Limit: 9},
			},
			competitionId: competitionId,
			mockBehavior: func(r *mock_service.MockSeasonInterface, competitionId int, queryParams domain.GetSeasonsQueryParams) {
				r.EXPECT().
					GetSeasons(gomock.Any(), competitionId, queryParams).
					Return(resStruct, 321, nil).
					AnyTimes()
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: string(resJson),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			r := mock_service.NewMockSeasonInterface(c)
			test.mockBehavior(r, competitionId, test.queryParams)

			cfg, _ := config.NewConfig()
			service := &service.Service{Season: r}
			handler := Handler{service, cfg}

			// Init Endpoint
			router := gin.New()
			router.GET("/competition/:competition_id/seasons", handler.getSeasons)

			// Create Request
			w := httptest.NewRecorder()

			req := httptest.NewRequest("GET", test.url, nil)

			// Make Request
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}
