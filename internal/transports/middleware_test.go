package handler_api

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/o-sokol-o/hub/pkg/jwt_processing"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestGetUserId(t *testing.T) {
	var getContext = func(id int) *gin.Context {
		ctx := &gin.Context{}
		ctx.Set(userCtx, id)
		return ctx
	}

	testTable := []struct {
		name       string
		ctx        *gin.Context
		id         int
		shouldFail bool
	}{
		{
			name: "Test case: Ok",
			ctx:  getContext(1),
			id:   1,
		},
		{
			ctx:        &gin.Context{},
			name:       "Test case: Empty",
			shouldFail: true,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			id, err := getUserIdFromContext(test.ctx)
			if test.shouldFail {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, id, test.id)
		})
	}
}

func getJWT(id int) string {
	jwt, _ := jwt_processing.GenerateToken(id)
	return jwt
}

func TestHandler_userIdentity_middleware(t *testing.T) {
	// Init Test Table

	testTable := []struct {
		name                 string
		headerName           string
		headerValue          string
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Test case: Ok",
			headerName:  "Authorization",
			headerValue: "Bearer " + getJWT(42),

			expectedStatusCode:   200,
			expectedResponseBody: "42",
		},
		{
			name:        "Test case: Invalid Header Name",
			headerName:  "",
			headerValue: "Bearer " + getJWT(1),

			expectedStatusCode:   401,
			expectedResponseBody: `{"status":"empty auth header"}`,
		},
		{
			name:        "Test case: Invalid Header Value",
			headerName:  "Authorization",
			headerValue: "Bearr " + getJWT(1),

			expectedStatusCode:   401,
			expectedResponseBody: `{"status":"invalid auth header"}`,
		},

		{
			name:        "Test case: Empty Token",
			headerName:  "Authorization",
			headerValue: "Bearer ",

			expectedStatusCode:   401,
			expectedResponseBody: `{"status":"token is empty"}`,
		},
		{
			name:        "Test case: Parse Error",
			headerName:  "Authorization",
			headerValue: "Bearer token",

			expectedStatusCode:   401,
			expectedResponseBody: `{"status":"invalid parse token"}`,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			// NewHandler(log *logrus.Logger, cache domain.Cache, a IServiceAuthentications, b IServiceChecklist, c IServiceChecklistItem, d IServiceAquahubList)
			h := NewHandler(logrus.New(), nil, nil, nil, nil, nil)

			// Init Endpoint
			r := gin.New()
			r.GET("/api", h.userIdentity_middleware, func(c *gin.Context) {
				id, _ := c.Get(userCtx)
				c.String(200, "%d", id)
			})

			// Init Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api", nil)
			req.Header.Set(test.headerName, test.headerValue)

			r.ServeHTTP(w, req)

			// Asserts
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}
