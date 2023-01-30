package tests

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/shakh9006/numbers-store/internal/apiserver/app"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func init() {
	gin.SetMode(gin.TestMode)
	zerolog.SetGlobalLevel(zerolog.Disabled)
}

func TestGetById(t *testing.T) {
	router := app.SetupApp()
	fmt.Println(router)
	testCases := []struct {
		name    string
		url     string
		status  string
		expCode int
	}{
		{
			name:    "valid",
			url:     "/v1/number/1",
			status:  "success",
			expCode: 200,
		},
		{
			name:    "not found",
			url:     "/v1/number/11",
			status:  "fail",
			expCode: 400,
		},
		{
			name:    "invalid params",
			url:     "/v1/number/s",
			status:  "fail",
			expCode: 400,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", tc.url, nil)
			req.Header.Set("Content-Type", "application/json; charset=utf-8")
			router.ServeHTTP(w, req)

			body, _ := io.ReadAll(w.Result().Body)
			var f interface{}
			json.Unmarshal(body, &f)
	
			myMap := f.(map[string]interface{})

			assert.Equal(t, tc.expCode, w.Code)
			assert.Equal(t, tc.status, myMap["status"])
			assert.NotNil(t, myMap["number"])
		})
	}
}
