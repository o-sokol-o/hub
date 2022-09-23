package handler_api

import (
	"bytes"
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mock_service "github.com/o-sokol-o/hub/internal/transports/mocks"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/o-sokol-o/hub/internal/domain"
)

func TestHandler_signUp(t *testing.T) {
	type mockBehavior func(s *mock_service.MockIServiceAuthentications, user domain.User)

	fmt.Println("----------------  Test signUp ---------------")

	testTable := []struct {
		name                 string
		inputBody            string
		inputUser            domain.User
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Test case: OK",
			inputBody: `{"first_name":"Test","last_name":"Test","email":"test@gmail.com","password":"qwerty"}`,
			inputUser: domain.User{
				FirstName:    "Test",
				LastName:     "Test",
				Email:        "test@gmail.com",
				PasswordHash: "qwerty",
			},
			mockBehavior: func(s *mock_service.MockIServiceAuthentications, user domain.User) {
				s.EXPECT().CreateUser(user).Return(1, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"token":"Bearer `,
		},
		{
			name:                 "Test case: Invalid input body",
			inputBody:            `{"email":"test@gmail.com","password":"qwerty"}`,
			mockBehavior:         func(s *mock_service.MockIServiceAuthentications, user domain.User) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"status":"User send invalid input body"}`,
		},
		{
			name:                 "Test case: Email not valid",
			inputBody:            `{"first_name":"Test","email":"testgmailcom","password":"qwerty"}`,
			mockBehavior:         func(s *mock_service.MockIServiceAuthentications, user domain.User) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"status":"User send invalid input body"}`,
		},

		{
			name:      "Test case: Service Failure",
			inputBody: `{"first_name":"Test","last_name":"Test","email":"test@gmail.com","password":"qwerty"}`,
			inputUser: domain.User{
				FirstName:    "Test",
				LastName:     "Test",
				Email:        "test@gmail.com",
				PasswordHash: "qwerty",
			},
			mockBehavior: func(s *mock_service.MockIServiceAuthentications, user domain.User) {
				s.EXPECT().CreateUser(user).Return(0, errors.New("service failure: something went wrong"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"status":"service failure: something went wrong"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			fmt.Printf("---------------- Start %s ---------------\n", testCase.name)

			//Init deps
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockIServiceAuthentications(c)
			testCase.mockBehavior(auth, testCase.inputUser)

			// NewHandler(log *logrus.Logger, cache domain.Cache, a IServiceAuthentications, b IServiceChecklist, c IServiceChecklistItem, d IServiceAquahubList)
			h := NewHandler(logrus.New(), nil, auth, nil, nil, nil)

			// Test Server
			router := gin.Default()
			router.POST("/auth/sign-up", h.signUp)

			// Create request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/auth/sign-up", bytes.NewBufferString(testCase.inputBody))

			// Send request
			router.ServeHTTP(w, req)

			//Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)

			body := w.Body.String()
			if testCase.name == "Test case: OK" {
				body = body[:17]
			}
			assert.Equal(t, testCase.expectedResponseBody, body)

			fmt.Printf("---------------- End %s ---------------\n", testCase.name)
		})
	}

}

func TestHandler_signIn(t *testing.T) {
	type mockBehavior func(s *mock_service.MockIServiceAuthentications, user signInInput)

	fmt.Println("----------------  Test signIn ---------------")

	testTable := []struct {
		name                 string
		inputBody            string
		inputUser            signInInput
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Test case: OK",
			inputBody: `{"email":"test@gmail.com","password":"qwerty"}`,
			inputUser: signInInput{
				Email:    "test@gmail.com",
				Password: "qwerty",
			},
			mockBehavior: func(s *mock_service.MockIServiceAuthentications, user signInInput) {
				s.EXPECT().Authenticate(user.Email, user.Password).Return(&domain.User{ID: 1}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"token":"Bearer `,
		},
		{
			name:                 "Test case: Invalid input body",
			inputBody:            `{"email":"test@gmail.com","pas__sword":"qwerty"}`,
			mockBehavior:         func(s *mock_service.MockIServiceAuthentications, user signInInput) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"status":"User send invalid input body"}`,
		},
		{
			name:                 "Test case: Email not valid",
			inputBody:            `{"email":"","password":"qwerty"}`,
			mockBehavior:         func(s *mock_service.MockIServiceAuthentications, user signInInput) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"status":"User send invalid input body"}`,
		},

		{
			name:      "Test case: Service Failure",
			inputBody: `{"email":"test@gmail.com","password":"qwerty"}`,
			inputUser: signInInput{
				Email:    "test@gmail.com",
				Password: "qwerty",
			},
			mockBehavior: func(s *mock_service.MockIServiceAuthentications, user signInInput) {
				s.EXPECT().Authenticate(user.Email, user.Password).Return(nil, errors.New("service failure: something went wrong"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"status":"service failure: something went wrong"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			fmt.Printf("---------------- Start %s ---------------\n", testCase.name)

			//Init deps
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockIServiceAuthentications(c)
			testCase.mockBehavior(auth, testCase.inputUser)

			// NewHandler(log *logrus.Logger, cache domain.Cache, a IServiceAuthentications, b IServiceChecklist, c IServiceChecklistItem, d IServiceAquahubList)
			h := NewHandler(logrus.New(), nil, auth, nil, nil, nil)

			// Test Server
			router := gin.Default()
			router.POST("/auth/sign-in", h.signIn)

			// Create request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/auth/sign-in", bytes.NewBufferString(testCase.inputBody))

			// Send request
			router.ServeHTTP(w, req)

			//Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)

			body := w.Body.String()
			if testCase.name == "Test case: OK" {
				body = body[:17]
			}
			assert.Equal(t, testCase.expectedResponseBody, body)

			fmt.Printf("---------------- End %s ---------------\n", testCase.name)
		})
	}

	fmt.Println()
	fmt.Println()

}
