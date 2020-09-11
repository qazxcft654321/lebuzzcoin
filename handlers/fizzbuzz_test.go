package handlers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"lebuzzcoin/core/cache"
	"lebuzzcoin/models"

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

	h := New(ioutil.Discard, &cache.Client{})
	e := echo.New()

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			os.Unsetenv("APIVERSION")
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

func TestComputeFizzbuzz(t *testing.T) {
	tests := map[string]struct {
		fizzbuzz *models.Fizzbuzz
		result   models.Result
		expected int
		repeat   bool
	}{
		"case1": {
			fizzbuzz: &models.Fizzbuzz{},
			result:   models.Result{},
			expected: http.StatusForbidden,
			repeat:   false,
		},

		"case2": {
			fizzbuzz: &models.Fizzbuzz{
				ModA:     4,
				ReplaceA: "less than 20Bytes",
			},
			result:   models.Result{},
			expected: http.StatusForbidden,
			repeat:   false,
		},

		"case3": {
			fizzbuzz: &models.Fizzbuzz{
				ModA:     32768,
				ReplaceA: "less than 20Bytes",
			},
			result:   models.Result{},
			expected: http.StatusForbidden,
			repeat:   false,
		},

		"case4": {
			fizzbuzz: &models.Fizzbuzz{
				ModA:     10,
				ReplaceA: "lets write something incorrect longer than 20 bytes",
			},
			result:   models.Result{},
			expected: http.StatusForbidden,
			repeat:   false,
		},

		"case5": {
			fizzbuzz: &models.Fizzbuzz{
				ModA:     1,
				ModB:     1,
				Limit:    32768,
				ReplaceA: "less than 20Bytes",
				ReplaceB: "less than 20Bytes",
			},
			result:   models.Result{},
			expected: http.StatusForbidden,
			repeat:   false,
		},

		"case6": {
			fizzbuzz: &models.Fizzbuzz{
				ModA:     1,
				ModB:     1,
				Limit:    100,
				ReplaceA: "less than 20Bytes",
				ReplaceB: "lets write something incorrect longer than 20 bytes",
			},
			result:   models.Result{},
			expected: http.StatusForbidden,
			repeat:   false,
		},

		"case7": {
			fizzbuzz: &models.Fizzbuzz{
				ModA:     32768,
				ModB:     32768,
				Limit:    32768,
				ReplaceA: "lets write something incorrect longer than 20 bytes",
				ReplaceB: "lets write something incorrect longer than 20 bytes",
			},
			result:   models.Result{},
			expected: http.StatusForbidden,
			repeat:   false,
		},

		"case8": {
			fizzbuzz: &models.Fizzbuzz{
				ModA:     2,
				ModB:     3,
				Limit:    1000,
				ReplaceA: "fizz",
				ReplaceB: "buzz",
			},
			result:   models.Result{},
			expected: http.StatusOK,
			repeat:   false,
		},
	}

	cache, err := cache.NewTestCache()
	assert.NoError(t, err)

	h := New(ioutil.Discard, cache)
	e := echo.New()

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			payload, err := json.Marshal(tc.fizzbuzz)
			assert.NoError(t, err)

			req := httptest.NewRequest("POST", "/test", bytes.NewBuffer(payload))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			e.POST("/test", h.ComputeFizzbuzz)
			e.ServeHTTP(w, req)

			result := models.Result{}
			err = json.Unmarshal(w.Body.Bytes(), &result)
			assert.NoError(t, err)

			assert.EqualValues(t, tc.result, result)
			assert.Equal(t, tc.expected, w.Code)
		})
	}
}
