package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"timeline/config"

	"timeline/internal/domain"
	"timeline/internal/service"
	"timeline/internal/service/mock"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	// "go.mongodb.org/mongo-driver/bson/primitive"
)

func TestHandler_signUp(t *testing.T) {
	type mockBehavior func(r *mock_service.MockUserInterface, input domain.UserSignUpInput)
	ctx := context.Background()
	tests := []struct {
		name                 string
		inputBody            string
		inputUser            domain.UserSignUpInput
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"username": "noname", "email": "noname@noname.org", "password": "qwerty"}`,
			inputUser: domain.UserSignUpInput{
				Username: "noname",
				Email:    "noname@noname.org",
				Password: "qwerty",
			},

			mockBehavior: func(r *mock_service.MockUserInterface, input domain.UserSignUpInput) {
				r.EXPECT().SignUp(ctx, input).Return(nil)
			},
			expectedStatusCode:   http.StatusCreated,
			expectedResponseBody: `{"message":"account created"}`,
		},
		{
			name:      "Wrong input",
			inputBody: `{"username": "noname", "password": "qwerty"}`,
			inputUser: domain.UserSignUpInput{
				Username: "noname",
				Email:    "noname@noname.org",
				Password: "qwerty",
			},
			mockBehavior:         func(r *mock_service.MockUserInterface, input domain.UserSignUpInput) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "User already exists",
			inputBody: `{"username": "noname", "email": "noname@noname.org", "password": "qwerty"}`,
			inputUser: domain.UserSignUpInput{
				Username: "noname",
				Email:    "noname@noname.org",
				Password: "qwerty",
			},
			mockBehavior: func(r *mock_service.MockUserInterface, input domain.UserSignUpInput) {
				r.EXPECT().SignUp(ctx, input).Return(domain.ErrUserAlreadyExists)
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"user already exists"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			r := mock_service.NewMockUserInterface(c)
			test.mockBehavior(r, test.inputUser)

			cfg, _ := config.NewConfig()
			service := &service.Service{User: r}
			handler := Handler{service, cfg}

			// Init Endpoint
			router := gin.New()
			router.POST("/sign-up", handler.signUp)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-up",
				bytes.NewBufferString(test.inputBody))

			// Make Request
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

// func TestHandler_signInUsername(t *testing.T) {
// 	type mockBehavior func(r *mock_service.MockUserInterface, input domain.UserSignInUsernameInput)

// 	userId := primitive.NewObjectID()
// 	ctx := context.Background()

// 	session := domain.Session{
// 		SignInDTO: domain.SignInDTO{
// 			User: domain.UserDTO{
// 				ID:          userId,
// 				Username:    "noname",
// 				Email:       "noname@noname.org",
// 			},
// 			AccessToken: "987654321",
// 		},
// 		RememberMe:   true,
// 		RefreshToken: "123456789",
// 	}

// 	tests := []struct {
// 		name                 string
// 		inputBody            string
// 		inputUser            domain.UserSignInUsernameInput
// 		mockBehavior         mockBehavior
// 		expectedStatusCode   int
// 		expectedResponseBody string
// 		expectedCookie       string
// 	}{
// 		{
// 			name:      "Ok",
// 			inputBody: `{"usernameOrEmail": "noname", "password": "qwerty", "rememberMe": true}`,
// 			inputUser: domain.UserSignInUsernameInput{
// 				Username:   "noname",
// 				Password:   "qwerty",
// 				RememberMe: true,
// 			},
// 			mockBehavior: func(r *mock_service.MockUserInterface, input domain.UserSignInUsernameInput) {
// 				r.EXPECT().SignInByUsername(ctx, input).Return(session, nil)
// 			},
// 			expectedStatusCode: http.StatusOK,
// 			expectedResponseBody: fmt.Sprintf(
// 				`{"user":{"id":"%s","username":"noname","email":"noname@noname.org"},"accessToken":"987654321"}`,
// 				userId.Hex(),
// 			),
// 			expectedCookie: fmt.Sprintf("%v", []*http.Cookie{
// 				{
// 					Name:     "refreshToken",
// 					Value:    "123456789",
// 					Path:     "/",
// 					MaxAge:   3888000,
// 					HttpOnly: true,
// 				},
// 			}),
// 		},
// 		{
// 			name:      "Wrong input",
// 			inputBody: `{"usernameOrEmail": "noname", "rememberMe": true}`,
// 			inputUser: domain.UserSignInUsernameInput{
// 				Username:   "noname",
// 				Password:   "qwerty",
// 				RememberMe: true,
// 			},
// 			mockBehavior:         func(r *mock_service.MockUserInterface, input domain.UserSignInUsernameInput) {},
// 			expectedStatusCode:   http.StatusBadRequest,
// 			expectedResponseBody: `{"message":"invalid input body"}`,
// 			expectedCookie:       fmt.Sprintf("%v", []*http.Cookie{}),
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Init Dependencies
// 			c := gomock.NewController(t)
// 			defer c.Finish()

// 			i := mock_service.NewMockUserInterface(c)
// 			tt.mockBehavior(i, tt.inputUser)

// 			cfg, _ := config.NewConfig()
// 			service := &service.Service{User: i}
// 			handler := Handler{service, cfg}

// 			// Init Endpoint
// 			router := gin.New()
// 			router.POST("/sign-in", handler.signIn)

// 			// Create Request
// 			w := httptest.NewRecorder()
// 			req := httptest.NewRequest("POST", "/sign-in",
// 				bytes.NewBufferString(tt.inputBody))

// 			// Make Request
// 			router.ServeHTTP(w, req)

// 			// Assert
// 			assert.Equal(t, w.Code, tt.expectedStatusCode)
// 			assert.Equal(t, w.Body.String(), tt.expectedResponseBody)
// 			assert.Equal(t, fmt.Sprintf("%v", w.Result().Cookies()), tt.expectedCookie)
// 		})
// 	}
// }

// func TestHandler_signInEmail(t *testing.T) {
// 	type mockBehavior func(r *mock_service.MockUserInterface, input domain.UserSignInEmailInput)

// 	userId := primitive.NewObjectID()
// 	ctx := context.Background()

// 	session := domain.Session{
// 		SignInDTO: domain.SignInDTO{
// 			User: domain.UserDTO{
// 				ID:          userId,
// 				Username:    "noname",
// 				Email:       "noname@noname.org",
// 			},
// 			AccessToken: "987654321",
// 		},
// 		// IsActivated: true,
// 		RememberMe:   true,
// 		RefreshToken: "123456789",
// 	}

// 	tests := []struct {
// 		name                 string
// 		inputBody            string
// 		inputUser            domain.UserSignInEmailInput
// 		mockBehavior         mockBehavior
// 		expectedStatusCode   int
// 		expectedResponseBody string
// 		// expectedCookie       []*http.Cookie
// 		expectedCookie string
// 	}{
// 		{
// 			name:      "Ok",
// 			inputBody: `{"usernameOrEmail": "noname@noname.org", "password": "qwerty", "rememberMe": true}`,
// 			inputUser: domain.UserSignInEmailInput{
// 				Email:      "noname@noname.org",
// 				Password:   "qwerty",
// 				RememberMe: true,
// 			},
// 			mockBehavior: func(r *mock_service.MockUserInterface, input domain.UserSignInEmailInput) {
// 				r.EXPECT().SignInByEmail(ctx, input).Return(session, nil)
// 			},
// 			expectedStatusCode: http.StatusOK,
// 			expectedResponseBody: fmt.Sprintf(
// 				`{"user":{"id":"%s","username":"noname","email":"noname@noname.org"},"accessToken":"987654321"}`,
// 				userId.Hex(),
// 			),
// 			expectedCookie: fmt.Sprintf("%v", []*http.Cookie{
// 				{
// 					Name:     "refreshToken",
// 					Value:    "123456789",
// 					Path:     "/",
// 					MaxAge:   3888000,
// 					HttpOnly: true,
// 				},
// 			},
// 			),
// 			// expectedCookie: "[refreshToken=123456789; Path=/; Max-Age=3888000; HttpOnly]",
// 		},
// 		{
// 			name:      "Wrong input",
// 			inputBody: `{"usernameOrEmail": "noname@noname.org", "rememberMe": true}`,
// 			inputUser: domain.UserSignInEmailInput{
// 				Email:      "noname@noname.org",
// 				Password:   "qwerty",
// 				RememberMe: true,
// 			},
// 			mockBehavior:         func(r *mock_service.MockUserInterface, input domain.UserSignInEmailInput) {},
// 			expectedStatusCode:   http.StatusBadRequest,
// 			expectedResponseBody: `{"message":"invalid input body"}`,
// 			expectedCookie:       fmt.Sprintf("%v", []*http.Cookie{}),
// 			// expectedCookie:       []*http.Cookie{},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Init Dependencies
// 			c := gomock.NewController(t)
// 			defer c.Finish()

// 			i := mock_service.NewMockUserInterface(c)
// 			tt.mockBehavior(i, tt.inputUser)

// 			cfg, _ := config.NewConfig()
// 			service := &service.Service{User: i}
// 			handler := Handler{service, cfg}

// 			// Init Endpoint
// 			router := gin.New()
// 			router.POST("/sign-in", handler.signIn)

// 			// Create Request
// 			w := httptest.NewRecorder()
// 			req := httptest.NewRequest("POST", "/sign-in",
// 				bytes.NewBufferString(tt.inputBody))

// 			// Make Request
// 			router.ServeHTTP(w, req)

// 			// Assert
// 			assert.Equal(t, w.Code, tt.expectedStatusCode)
// 			assert.Equal(t, w.Body.String(), tt.expectedResponseBody)

// 			// fmt.Printf("%#v", w.Result().Cookies())
// 			// fmt.Println("")
// 			// fmt.Printf("%#v", tt.expectedCookie)
// 			// fmt.Println("")
// 			// fmt.Printf("%v", w.Result().Cookies())
// 			// fmt.Println("")
// 			// fmt.Printf("%v", tt.expectedCookie)
// 			// fmt.Println("")

// 			// assert.Equal(t, w.Result().Cookies(), tt.expectedCookie)
			
			
// 			assert.Equal(t, fmt.Sprintf("%v", w.Result().Cookies()), tt.expectedCookie)
// 		})
// 	}
// }