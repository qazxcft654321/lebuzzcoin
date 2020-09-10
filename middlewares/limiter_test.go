package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/didip/tollbooth"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestLimitMiddleware(t *testing.T) {
	tests := map[string]struct {
		limitPerSecond float64
		loop           int
		expected       int
	}{
		"case1": {
			limitPerSecond: 5,
			loop:           3,
			expected:       200,
		},

		"case2": {
			limitPerSecond: 1,
			loop:           10,
			expected:       429,
		},
	}

	fakeHTTPHandler := func(c echo.Context) error {
		return c.String(http.StatusOK, "fake")
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			lmt := tollbooth.NewLimiter(tc.limitPerSecond, nil)
			e := echo.New()
			e.Use(LimitMiddleware(lmt))

			var code int

			for i := 0; i < tc.loop; i++ {
				req := httptest.NewRequest("GET", "/test", nil)
				w := httptest.NewRecorder()

				e.GET("/test", fakeHTTPHandler)
				e.ServeHTTP(w, req)
				code = w.Code
			}

			assert.Equal(t, tc.expected, code)
		})
	}
}
