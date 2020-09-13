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
		response string
		expected int
		repeat   bool
	}{
		"case1": {
			fizzbuzz: &models.Fizzbuzz{},
			expected: http.StatusForbidden,
			repeat:   false,
		},

		"case2": {
			fizzbuzz: &models.Fizzbuzz{
				ModA:     4,
				ReplaceA: "less than 20Bytes",
			},
			expected: http.StatusForbidden,
			repeat:   false,
		},

		"case3": {
			fizzbuzz: &models.Fizzbuzz{
				ModA:     32768,
				ReplaceA: "less than 20Bytes",
			},
			expected: http.StatusForbidden,
			repeat:   false,
		},

		"case4": {
			fizzbuzz: &models.Fizzbuzz{
				ModA:     10,
				ReplaceA: "lets write something incorrect longer than 20 bytes",
			},
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
			expected: http.StatusForbidden,
			repeat:   false,
		},

		"case8": {
			fizzbuzz: &models.Fizzbuzz{
				ModA:     2,
				ModB:     3,
				Limit:    100,
				ReplaceA: "testA",
				ReplaceB: "testB",
			},
			response: "{\"data\":{\"hash\":\"dea6044b07f822744ab3de793dc4d9e5d0fd7c572e7c2c3e68fafe182551d110\",\"fizzbuzz\":{\"mod_a\":2,\"mod_b\":3,\"limit\":100,\"replace_a\":\"testA\",\"replace_b\":\"testB\"},\"result\":[\"1\",\"testA\",\"testB\",\"testA\",\"testAtestB\",\"testA\",\"7\",\"testA\",\"testB\",\"testAtestB\",\"11\",\"testA\",\"13\",\"testA\",\"testAtestB\",\"testA\",\"17\",\"testA\",\"19\",\"testAtestB\",\"testB\",\"testA\",\"23\",\"testA\",\"testAtestB\",\"testA\",\"testB\",\"testA\",\"29\",\"testAtestB\",\"31\",\"testA\",\"testB\",\"testA\",\"testAtestB\",\"testA\",\"37\",\"testA\",\"testB\",\"testAtestB\",\"41\",\"testA\",\"43\",\"testA\",\"testAtestB\",\"testA\",\"47\",\"testA\",\"49\",\"testAtestB\",\"testB\",\"testA\",\"53\",\"testA\",\"testAtestB\",\"testA\",\"testB\",\"testA\",\"59\",\"testAtestB\",\"61\",\"testA\",\"testB\",\"testA\",\"testAtestB\",\"testA\",\"67\",\"testA\",\"testB\",\"testAtestB\",\"71\",\"testA\",\"73\",\"testA\",\"testAtestB\",\"testA\",\"77\",\"testA\",\"79\",\"testAtestB\",\"testB\",\"testA\",\"83\",\"testA\",\"testAtestB\",\"testA\",\"testB\",\"testA\",\"89\",\"testAtestB\",\"91\",\"testA\",\"testB\",\"testA\",\"testAtestB\",\"testA\",\"97\",\"testA\",\"testB\"],\"state\":\"built\"},\"status\":\"success\"}\n",
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

			if len(tc.response) > 0 {
				assert.EqualValues(t, tc.response, w.Body.String())
			}
			assert.Equal(t, tc.expected, w.Code)
		})
	}
}

func TestGetFizzbuzzFromHash(t *testing.T) {
	tests := map[string]struct {
		hash     string
		data     *models.Result
		response string
		cache    bool
		expected int
	}{
		"case1": {
			hash:     "qwert+=123",
			data:     nil,
			cache:    false,
			expected: http.StatusForbidden,
		},

		"case2": {
			hash:     "8c226d1fc7d66d45192f02eeefd85c0d99b729926c17c1632a32f53dbbd5657a",
			data:     nil,
			cache:    false,
			expected: http.StatusBadRequest,
		},

		"case3": {
			hash: "d695743ec1d546f0bfec475d27edf126212b02b457b724d4ff95dd3e5e3d2476",
			data: &models.Result{
				Hash: "d695743ec1d546f0bfec475d27edf126212b02b457b724d4ff95dd3e5e3d2476",
				Fizzbuzz: &models.Fizzbuzz{
					ModA:     2,
					ModB:     3,
					Limit:    1000,
					ReplaceA: "testA",
					ReplaceB: "testB",
				},
			},
			response: "{\"data\":{\"hash\":\"d695743ec1d546f0bfec475d27edf126212b02b457b724d4ff95dd3e5e3d2476\",\"fizzbuzz\":{\"mod_a\":2,\"mod_b\":3,\"limit\":1000,\"replace_a\":\"testA\",\"replace_b\":\"testB\"},\"result\":null,\"state\":\"cached\"},\"status\":\"success\"}\n",
			cache:    true,
			expected: http.StatusOK,
		},
	}

	cache, err := cache.NewTestCache()
	assert.NoError(t, err)

	h := New(ioutil.Discard, cache)
	e := echo.New()

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.cache {
				json, err := json.Marshal(tc.data)
				assert.NoError(t, err)

				err = cache.Set(tc.hash, json, 0)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest("GET", "/test/"+tc.hash, nil)
			w := httptest.NewRecorder()
			e.GET("/test/:hash", h.GetFizzbuzzFromHash)
			e.ServeHTTP(w, req)

			if len(tc.response) > 0 {
				assert.EqualValues(t, tc.response, w.Body.String())
			}
			assert.Equal(t, tc.expected, w.Code)
		})
	}
}

func TestGetComputeByDescHitScore(t *testing.T) {
	tests := map[string]struct {
		data      []*models.Result
		response  string
		sortedSet bool
		expected  int
	}{
		"case1": {
			response: "{\"data\":[],\"status\":\"success\"}\n",
			expected: http.StatusOK,
		},

		"case2": {
			data: []*models.Result{
				&models.Result{
					Hash: "d695743ec1d546f0bfec475d27edf126212b02b457b724d4ff95dd3e5e3d2476",
					Fizzbuzz: &models.Fizzbuzz{
						ModA:     2,
						ModB:     3,
						Limit:    1000,
						ReplaceA: "testA",
						ReplaceB: "testB",
					},
				},
				&models.Result{
					Hash: "8e7a9c5cf786f6b80e6f1708b44b8d4fc298761711112ba8009484938fb05c2a",
					Fizzbuzz: &models.Fizzbuzz{
						ModA:     2,
						ModB:     4,
						Limit:    1000,
						ReplaceA: "testA",
						ReplaceB: "testB",
					},
				},
				&models.Result{
					Hash: "8c226d1fc7d66d45192f02eeefd85c0d99b729926c17c1632a32f53dbbd5657a",
					Fizzbuzz: &models.Fizzbuzz{
						ModA:     2,
						ModB:     5,
						Limit:    1000,
						ReplaceA: "testA",
						ReplaceB: "testB",
					},
				},
			},
			response:  "{\"data\":[{\"Hits\":3,\"Result\":{\"hash\":\"8c226d1fc7d66d45192f02eeefd85c0d99b729926c17c1632a32f53dbbd5657a\",\"fizzbuzz\":{\"mod_a\":2,\"mod_b\":5,\"limit\":1000,\"replace_a\":\"testA\",\"replace_b\":\"testB\"},\"result\":null,\"state\":\"cached\"}},{\"Hits\":2,\"Result\":{\"hash\":\"8e7a9c5cf786f6b80e6f1708b44b8d4fc298761711112ba8009484938fb05c2a\",\"fizzbuzz\":{\"mod_a\":2,\"mod_b\":4,\"limit\":1000,\"replace_a\":\"testA\",\"replace_b\":\"testB\"},\"result\":null,\"state\":\"cached\"}},{\"Hits\":1,\"Result\":{\"hash\":\"d695743ec1d546f0bfec475d27edf126212b02b457b724d4ff95dd3e5e3d2476\",\"fizzbuzz\":{\"mod_a\":2,\"mod_b\":3,\"limit\":1000,\"replace_a\":\"testA\",\"replace_b\":\"testB\"},\"result\":null,\"state\":\"cached\"}}],\"status\":\"success\"}\n",
			sortedSet: true,
			expected:  http.StatusOK,
		},
	}

	cacheClient, err := cache.NewTestCache()
	assert.NoError(t, err)

	h := New(ioutil.Discard, cacheClient)
	e := echo.New()

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := cacheClient.FlushAll()
			assert.NoError(t, err)
			// cache data
			if len(tc.data) > 0 {
				for _, v := range tc.data {
					json, err := json.Marshal(v)
					assert.NoError(t, err)

					err = h.cache.Set(v.Hash, json, 0)
					assert.NoError(t, err)
				}
			}

			// sorted set
			if tc.sortedSet {
				for i := 0; i < len(tc.data); i++ {
					_, err := h.cache.ZAdd(SortedSetKey, &cache.ZMember{
						Score:  float64(i + 1),
						Member: tc.data[i].Hash,
					})
					assert.NoError(t, err)
				}
			}

			// http req
			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()
			e.GET("/test", h.GetComputeByDescHitScore)
			e.ServeHTTP(w, req)

			assert.EqualValues(t, tc.response, w.Body.String())
			assert.Equal(t, tc.expected, w.Code)
		})
	}
}
