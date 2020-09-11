package handlers

import (
	"io/ioutil"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetAPIVersionHandler(t *testing.T) {
	tests := map[string]struct {
		init     bool
		expected int
	}{
		"case1": {
			init:     false,
			expected: 400,
		},

		"case2": {
			init:     true,
			expected: 200,
		},
	}

	h := New(ioutil.Discard, &redis.Client{})
	e := echo.New()
	os.Unsetenv("APIVERSION")

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.init {
				os.Setenv("APIVERSION", "6.6.6")
			}

			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()
			e.GET("/test", h.GetAPIVersion)
			e.ServeHTTP(w, req)

			assert.Equal(t, tc.expected, w.Code)
		})
	}
}
