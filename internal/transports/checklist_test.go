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

func getPointerString(s string) *string {
	return &s
}

func TestHandler_createList(t *testing.T) {
	type mockBehavior func(s *mock_service.MockIServiceChecklist, Checklist domain.CreateChecklist)

	fmt.Println("----------------  Test signUp ---------------")

	testTable := []struct {
		name                 string
		inputCtxUserId       int
		inputBody            string
		inputChecklist       domain.CreateChecklist
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:           "Test case: OK",
			inputCtxUserId: 1,
			inputBody:      `{"title":"Checklist Title","description":"Checklist Description"}`,
			inputChecklist: domain.CreateChecklist{
				Title:       getPointerString("Checklist Title"),
				Description: getPointerString("Checklist Description"),
			},
			mockBehavior: func(s *mock_service.MockIServiceChecklist, Checklist domain.CreateChecklist) {
				s.EXPECT().Create(1, Checklist).Return(1, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1}`,
		},
		{
			name:                 "Test case: Invalid input body",
			inputCtxUserId:       1,
			inputBody:            `{"tle":"Checklist Title"}`,
			mockBehavior:         func(s *mock_service.MockIServiceChecklist, Checklist domain.CreateChecklist) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"status":"User send invalid input body"}`,
		},
		{
			name:                 "Test case: User is unauthorized",
			inputCtxUserId:       0,
			mockBehavior:         func(s *mock_service.MockIServiceChecklist, Checklist domain.CreateChecklist) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"status":"user is unauthorized"}`,
		},
		{
			name:           "Test case: Service Failure",
			inputCtxUserId: 1,
			inputBody:      `{"title":"Checklist Title","description":"Checklist Description"}`,
			inputChecklist: domain.CreateChecklist{
				Title:       getPointerString("Checklist Title"),
				Description: getPointerString("Checklist Description"),
			},
			mockBehavior: func(s *mock_service.MockIServiceChecklist, Checklist domain.CreateChecklist) {
				s.EXPECT().Create(1, Checklist).Return(0, errors.New("service failure: something went wrong"))
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

			chk := mock_service.NewMockIServiceChecklist(c)
			testCase.mockBehavior(chk, testCase.inputChecklist)

			// NewHandler(log *logrus.Logger, cache domain.Cache, a IServiceAuthentications, b IServiceChecklist, c IServiceChecklistItem, d IServiceAquahubList)
			h := NewHandler(logrus.New(), nil, nil, chk, nil, nil)

			// Init Endpoint
			router := gin.Default()
			router.POST(
				"/api/lists/",

				func(c *gin.Context) {
					if testCase.inputCtxUserId == 1 {
						c.Set(userCtx, 1)
					}
				},

				h.createList)

			// Create request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/lists/", bytes.NewBufferString(testCase.inputBody))

			// Send request
			router.ServeHTTP(w, req)

			//Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			body := w.Body.String()
			assert.Equal(t, testCase.expectedResponseBody, body)

			fmt.Printf("---------------- End %s ---------------\n", testCase.name)
		})
	}
}
