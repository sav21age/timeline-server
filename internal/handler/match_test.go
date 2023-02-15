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
	"go.mongodb.org/mongo-driver/bson/primitive"

	// "github.com/magiconair/properties/assert"
	"github.com/go-playground/assert/v2"
)

func TestHandler_getMatches(t *testing.T) {
	type mockBehavior func(r *mock_service.MockMatchInterface, competitionId, seasonId int, queryParams domain.GetMatchesQueryParams)

	competitionId := 2021
	seasonId := 1490

	url := fmt.Sprintf("/competition/%d/season/%d/matches", competitionId, seasonId)

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

	resStruct := []domain.MatchDTO{
		{
			ID: uint(416005),
			Competition: domain.CompetitionMatchShort{
				ID:   uint16(2021),
				Name: "Premier League",
				Area: domain.AreaMatch{
					Name:        "England",
					CountryCode: "",
					Ensign:      "770.svg",
				},
			},
			Season: domain.SeasonMatchShort{
				ID:        uint16(1490),
				StartDate: sd,
				EndDate:   ed,
			},
			UtcDate:  ed,
			Status:   "SCHEDULED",
			Matchday: 38,
			Score: domain.ScoreMatchShort{
				FullTime: domain.MatchScore{},
			},
			HomeTeam: domain.TeamShortMatch{
				ID:   uint16(340),
				Name: "Southampton FC",
			},
			AwayTeam: domain.TeamShortMatch{
				ID:   64,
				Name: "Liverpool FC",
			},
		},
	}

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
		queryParams          domain.GetMatchesQueryParams
		competitionId        int
		seasonId             int
	}{
		{
			name: "200 empty",
			url:  url,
			queryParams: domain.GetMatchesQueryParams{
				Datination:  domain.Datination{MinDate: time.Time{}, MaxDate: time.Time{}},
				ClubId:      0,
				MatchStatus: "",
			},
			competitionId: competitionId,
			seasonId:      seasonId,
			mockBehavior: func(r *mock_service.MockMatchInterface, competitionId, seasonId int, queryParams domain.GetMatchesQueryParams) {
				r.EXPECT().
					GetMatches(gomock.Any(), competitionId, seasonId, queryParams).
					// Return([]domain.MatchDTO{}, nil).
					Return(resStruct, nil).
					AnyTimes()
			},
			expectedStatusCode: http.StatusOK,
			// expectedResponseBody: `{"message":"invalid query parameter: page or per_page"}`,
			expectedResponseBody: string(resJson),
		},
		{
			name: "400 competitionId:EPL",
			url:  fmt.Sprintf("/competition/EPL/season/%d/matches", seasonId),
			queryParams: domain.GetMatchesQueryParams{
				Datination:  domain.Datination{MinDate: sd, MaxDate: ed},
				ClubId:      0,
				MatchStatus: "",
			},
			competitionId: competitionId,
			seasonId:      seasonId,
			mockBehavior: func(r *mock_service.MockMatchInterface, competitionId, seasonId int, queryParams domain.GetMatchesQueryParams) {
				r.EXPECT().
					GetMatches(gomock.Any(), "EPL", seasonId, queryParams).
					Return(resStruct, nil).
					AnyTimes()
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"invalid query parameter: competition_id"}`,
		},
		{
			name: "400 seasonId:NEW",
			url:  fmt.Sprintf("/competition/%d/season/NEW/matches", competitionId),
			queryParams: domain.GetMatchesQueryParams{
				Datination:  domain.Datination{MinDate: sd, MaxDate: ed},
				ClubId:      0,
				MatchStatus: "",
			},
			competitionId: competitionId,
			seasonId:      seasonId,
			mockBehavior: func(r *mock_service.MockMatchInterface, competitionId, seasonId int, queryParams domain.GetMatchesQueryParams) {
				r.EXPECT().
					GetMatches(gomock.Any(), competitionId, "NEW", queryParams).
					Return(resStruct, nil).
					AnyTimes()
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"invalid query parameter: season_id"}`,
		},
		{
			name: "200 club_id=57",
			url:  fmt.Sprintf("%s?club_id=57", url),
			queryParams: domain.GetMatchesQueryParams{
				Datination:  domain.Datination{MinDate: time.Time{}, MaxDate: time.Time{}},
				ClubId:      57,
				MatchStatus: "",
			},
			competitionId: competitionId,
			seasonId:      seasonId,
			mockBehavior: func(r *mock_service.MockMatchInterface, competitionId, seasonId int, queryParams domain.GetMatchesQueryParams) {
				r.EXPECT().
					GetMatches(gomock.Any(), competitionId, seasonId, queryParams).
					Return(resStruct, nil).
					AnyTimes()
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: string(resJson),
		},
		{
			name: "200 match_status=FINISHED",
			url:  fmt.Sprintf("%s?match_status=FINISHED", url),
			queryParams: domain.GetMatchesQueryParams{
				Datination:  domain.Datination{MinDate: time.Time{}, MaxDate: time.Time{}},
				ClubId:      0,
				MatchStatus: "FINISHED",
			},
			competitionId: competitionId,
			seasonId:      seasonId,
			mockBehavior: func(r *mock_service.MockMatchInterface, competitionId, seasonId int, queryParams domain.GetMatchesQueryParams) {
				r.EXPECT().
					GetMatches(gomock.Any(), competitionId, seasonId, queryParams).
					Return(resStruct, nil).
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

			r := mock_service.NewMockMatchInterface(c)
			test.mockBehavior(r, competitionId, seasonId, test.queryParams)

			cfg, _ := config.NewConfig()
			service := &service.Service{Match: r}
			handler := Handler{service, cfg}

			// Init Endpoint
			router := gin.New()
			router.GET("/competition/:competition_id/season/:season_id/matches", handler.getMatches)

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

func TestHandler_getDates(t *testing.T) {
	type mockBehavior func(r *mock_service.MockMatchInterface, competitionId, seasonId int, queryParams domain.GetDatesQueryParams)

	competitionId := 2021
	seasonId := 1490
	url := fmt.Sprintf("/competition/%d/season/%d/dates", competitionId, seasonId)

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

	resArray := []primitive.DateTime{
		primitive.NewDateTimeFromTime(sd),
		primitive.NewDateTimeFromTime(ed),
	}

	resJson, err := json.Marshal(resArray)
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
		queryParams          domain.GetDatesQueryParams
		competitionId        int
		seasonId             int
	}{
		{
			name: "200 empty",
			url:  url,
			queryParams: domain.GetDatesQueryParams{
				ClubId: 0,
				MatchStatus: "",
			},
			competitionId: competitionId,
			seasonId: seasonId,
			mockBehavior: func(r *mock_service.MockMatchInterface, competitionId, seasonId int, queryParams domain.GetDatesQueryParams) {
				r.EXPECT().
					GetMatchesDates(gomock.Any(), competitionId, seasonId, queryParams).
					Return(resArray, int(2), nil).
					AnyTimes()
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: string(resJson),
		},
		{
			name: "400 competitionId:EPL",
			url:  fmt.Sprintf("/competition/EPL/season/%d/dates", seasonId),
			queryParams: domain.GetDatesQueryParams{
				ClubId: 0,
				MatchStatus: "",
			},
			competitionId: competitionId,
			seasonId: seasonId,
			mockBehavior: func(r *mock_service.MockMatchInterface, competitionId, seasonId int, queryParams domain.GetDatesQueryParams) {
				r.EXPECT().
					GetMatchesDates(gomock.Any(), competitionId, seasonId, queryParams).
					Return(resArray, int(2), nil).
					AnyTimes()
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"invalid query parameter: competition_id"}`,
		},		
		{
			name: "400 seasonId:NEW",
			url:  fmt.Sprintf("/competition/%d/season/NEW/dates", competitionId),
			queryParams: domain.GetDatesQueryParams{
				ClubId: 0,
				MatchStatus: "",
			},
			competitionId: competitionId,
			seasonId: seasonId,
			mockBehavior: func(r *mock_service.MockMatchInterface, competitionId, seasonId int, queryParams domain.GetDatesQueryParams) {
				r.EXPECT().
					GetMatchesDates(gomock.Any(), competitionId, seasonId, queryParams).
					Return(resArray, int(2), nil).
					AnyTimes()
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"invalid query parameter: season_id"}`,
		},
		{
			name: "200 club_id=57",
			url:  fmt.Sprintf("%s?club_id=57", url),
			queryParams: domain.GetDatesQueryParams{
				ClubId: 57,
				MatchStatus: "",
			},
			competitionId: competitionId,
			seasonId: seasonId,
			mockBehavior: func(r *mock_service.MockMatchInterface, competitionId, seasonId int, queryParams domain.GetDatesQueryParams) {
				r.EXPECT().
					GetMatchesDates(gomock.Any(), competitionId, seasonId, queryParams).
					Return(resArray, int(2), nil).
					AnyTimes()
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: string(resJson),
		},
		{
			name: "200 match_status=FINISHED",
			url:  fmt.Sprintf("%s?match_status=FINISHED", url),
			queryParams: domain.GetDatesQueryParams{
				ClubId: 0,
				MatchStatus: "FINISHED",
			},
			competitionId: competitionId,
			seasonId: seasonId,
			mockBehavior: func(r *mock_service.MockMatchInterface, competitionId, seasonId int, queryParams domain.GetDatesQueryParams) {
				r.EXPECT().
					GetMatchesDates(gomock.Any(), competitionId, seasonId, queryParams).
					Return(resArray, int(2), nil).
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

			r := mock_service.NewMockMatchInterface(c)
			test.mockBehavior(r, competitionId, seasonId, test.queryParams)

			cfg, _ := config.NewConfig()
			service := &service.Service{Match: r}
			handler := Handler{service, cfg}

			// Init Endpoint
			router := gin.New()
			router.GET("/competition/:competition_id/season/:season_id/dates", handler.getDates)

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
