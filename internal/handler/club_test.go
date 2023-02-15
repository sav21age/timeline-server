package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"timeline/config"

	"timeline/internal/domain"
	"timeline/internal/service"
	"timeline/internal/service/mock"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"

	// "github.com/magiconair/properties/assert"
	"github.com/go-playground/assert/v2"
)

func TestHandler_getClubs(t *testing.T) {
	type mockBehavior func(r *mock_service.MockClubInterface, queryParams domain.GetClubsQueryParams)

	url := "/clubs"

	resStruct := []domain.ClubDTO{{
		ID:        uint16(58),
		ShortName: "Aston Villa",
		Name:      "Aston Villa FC",
		Crest:     "58.svg",
		Area: domain.AreaClub{
			ID:   uint16(2072),
			Name: "England",
			Code: "ENG",
			Flag: "2072.svg",
		},
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
		queryParams          domain.GetClubsQueryParams
	}{
		{
			name: "500 empty",
			url:  url,
			queryParams: domain.GetClubsQueryParams{
				Pagination: domain.Pagination{Skip: 0, Limit: 0},
				AreaId:     0,
				SortBy:     "",
			},
			mockBehavior: func(r *mock_service.MockClubInterface, queryParams domain.GetClubsQueryParams) {
				r.EXPECT().
					GetClubs(gomock.Any(), queryParams).
					Return([]domain.ClubDTO{}, 0, nil).
					AnyTimes()
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"message":"invalid query parameter: page"}`,
		},
		{
			name: "204 page=1000",
			url:  fmt.Sprintf("%s?page=1000", url),
			queryParams: domain.GetClubsQueryParams{
				Pagination: domain.Pagination{Skip: 11988, Limit: 12},
				AreaId:     0,
				SortBy:     "",
			},
			mockBehavior: func(r *mock_service.MockClubInterface, queryParams domain.GetClubsQueryParams) {
				r.EXPECT().
					GetClubs(gomock.Any(), queryParams).
					Return([]domain.ClubDTO{}, 0, nil).
					AnyTimes()
			},
			expectedStatusCode:   http.StatusNoContent,
			expectedResponseBody: "",
		},
		{
			name: "200 page=1",
			url:  fmt.Sprintf("%s?page=1", url),
			queryParams: domain.GetClubsQueryParams{
				Pagination: domain.Pagination{Skip: 0, Limit: 12},
				AreaId:     0,
				SortBy:     "",
			},
			mockBehavior: func(r *mock_service.MockClubInterface, queryParams domain.GetClubsQueryParams) {
				r.EXPECT().
					GetClubs(gomock.Any(), queryParams).
					Return(resStruct, 321, nil).
					AnyTimes()
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: string(resJson),
		},
		{
			name: "200 page=1&per_page=12",
			url:  fmt.Sprintf("%s?page=1&per_page=12", url),
			queryParams: domain.GetClubsQueryParams{
				Pagination: domain.Pagination{Skip: 0, Limit: 12},
				AreaId:     0,
				SortBy:     "",
			},
			mockBehavior: func(r *mock_service.MockClubInterface, queryParams domain.GetClubsQueryParams) {
				r.EXPECT().
					GetClubs(gomock.Any(), queryParams).
					Return(resStruct, 321, nil).
					AnyTimes()
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: string(resJson),
		},
		{
			name: "200 page=1&per_page=12&area_id=2002",
			url:  fmt.Sprintf("%s?page=1&per_page=12&area_id=2002", url),
			queryParams: domain.GetClubsQueryParams{
				Pagination: domain.Pagination{Skip: 0, Limit: 12},
				AreaId:     2002,
				SortBy:     "",
			},
			mockBehavior: func(r *mock_service.MockClubInterface, queryParams domain.GetClubsQueryParams) {
				r.EXPECT().
					GetClubs(gomock.Any(), queryParams).
					Return(resStruct, 321, nil).
					AnyTimes()
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: string(resJson),
		},
		{
			name: "200 page=1&per_page=12&area_id=2002&sort_by=club",
			url:  fmt.Sprintf("%s?page=1&per_page=12&area_id=2002&sort_by=club", url),
			queryParams: domain.GetClubsQueryParams{
				Pagination: domain.Pagination{Skip: 0, Limit: 12},
				AreaId:     2002,
				SortBy:     "club",
			},
			mockBehavior: func(r *mock_service.MockClubInterface, queryParams domain.GetClubsQueryParams) {
				r.EXPECT().
					GetClubs(gomock.Any(), queryParams).
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

			r := mock_service.NewMockClubInterface(c)
			test.mockBehavior(r, test.queryParams)

			cfg, _ := config.NewConfig()
			service := &service.Service{Club: r}
			handler := Handler{service, cfg}

			// Init Endpoint
			router := gin.New()
			router.GET("/clubs", handler.getClubs)

			// Create Request
			w := httptest.NewRecorder()

			req := httptest.NewRequest("GET", test.url, nil)

			// req.Header.Set("Content-type", "application/json")
			// req.URL.Query().Add("token", "someval")
			// q := req.URL.Query()
			// q.Add("page", "5")
			// q.Add("per_page", "6")
			// req.URL.RawQuery = q.Encode()

			// Make Request
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}
