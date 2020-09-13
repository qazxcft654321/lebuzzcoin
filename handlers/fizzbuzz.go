package handlers

import (
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"

	"lebuzzcoin/core/cache"
	"lebuzzcoin/models"

	"github.com/labstack/echo/v4"
)

const SortedSetKey string = "computeScore"

func (h *Handler) GetAPIVersion(c echo.Context) error {
	APIVersion := os.Getenv("APIVERSION")
	if len(APIVersion) != 5 {
		h.LogErrorMessage("handlers.fizzbuzz", nil, "Error retrieving API version from env variable")
		return h.RespondJSONBadRequest()
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status":  RESPONSE_STATUS_SUCCESS,
		"message": "Lebuzzcoin-API v" + APIVersion,
	})
}

// TODO: Recheck int parameters limit range
func (h *Handler) ComputeFizzbuzz(c echo.Context) error {
	fizzbuzz := &models.Fizzbuzz{}
	if err := c.Bind(fizzbuzz); err != nil {
		h.LogErrorMessage("handlers.fizzbuzz", err, "Error decoding request body")
		return h.RespondJSONBadRequest()
	}

	if err := h.validator.Struct(fizzbuzz); err != nil {
		h.LogErrorMessage("handlers.fizzbuzz", err, "Data validation failed")
		return h.RespondJSONForbidden()
	}

	hash := fizzbuzz.HashData()

	// Retrieve from cache
	data, err := h.cache.Get(hash)
	if err != nil && err.Error() != cache.Empty {
		h.LogErrorMessage("handlers.fizzbuzz", err, "Error retrieving data from cache")
		return h.RespondJSONBadRequest()
	}

	// Building result from cache
	result := &models.Result{}
	if len(data) > 1 {
		err := json.Unmarshal([]byte(data), result)
		if err != nil {
			h.LogErrorMessage("handlers.fizzbuzz", err, "Error decoding struct from json")
			return h.RespondJSONBadRequest()
		}
		result.Flag = "cached"

		// Increment sorted set member
		go func() {
			_, err := h.cache.ZIncr(SortedSetKey, &cache.ZMember{Score: 1, Member: hash})
			if err != nil {
				h.LogErrorMessage("handlers.fizzbuzz", err, "Error caching sorted set")
			}
		}()
	}

	// Building new result
	if result.Hash != hash {
		result.Hash = hash
		result.Fizzbuzz = fizzbuzz
		// TODO: process calc

		json, err := json.Marshal(result)
		if err != nil {
			h.LogErrorMessage("handlers.fizzbuzz", err, "Error encoding struct to json")
			return h.RespondJSONBadRequest()
		}

		err = h.cache.Set(hash, json, 0)
		if err != nil {
			h.LogErrorMessage("handlers.fizzbuzz", err, "Error caching data")
			return h.RespondJSONBadRequest()
		}
		result.Flag = "built"

		// Add member to sorted set
		go func() {
			_, err := h.cache.ZAdd(SortedSetKey, &cache.ZMember{Score: 1, Member: hash})
			if err != nil {
				h.LogErrorMessage("handlers.fizzbuzz", err, "Error caching sorted set")
			}
		}()
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": RESPONSE_STATUS_SUCCESS,
		"data":   result,
	})
}

func (h *Handler) GetFizzbuzzFromHash(c echo.Context) error {
	hash := c.Param("hash")
	buffer, err := hex.DecodeString(hash)

	// Check if hex & 256bits / 32Bytes long
	if err != nil || len(buffer) != 32 {
		h.LogErrorMessage("handlers.fizzbuzz", nil, "Submited hash incorrect")
		return h.RespondJSONForbidden()
	}

	// Retrieve from cache
	cache, err := h.cache.Get(hash)
	if err != nil {
		h.LogErrorMessage("handlers.fizzbuzz", err, "Error retrieving data from cache")
		return h.RespondJSONBadRequest()
	}

	// Building result from cache
	result := &models.Result{}
	if len(cache) > 1 {
		err := json.Unmarshal([]byte(cache), result)
		if err != nil {
			h.LogErrorMessage("handlers.fizzbuzz", err, "Error decoding struct from json")
			return h.RespondJSONBadRequest()
		}
		result.Flag = "cached"
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": RESPONSE_STATUS_SUCCESS,
		"data":   result,
	})
}

// TODO: make test
func (h *Handler) GetComputeByDescHitScore(c echo.Context) error {
	members, err := h.cache.ZRevRangeWithScores(SortedSetKey, 0, 3)
	if err != nil {
		h.LogErrorMessage("handlers.fizzbuzz", err, "Error retrieving data from cache")
		return h.RespondJSONBadRequest()
	}

	type Stats struct {
		Hits   float64
		Result *models.Result
	}

	// Retrieve from cache
	stats := make([]Stats, 0, len(members))
	if len(members) > 0 {
		for _, v := range members {
			cache, err := h.cache.Get(v.Member.(string))
			if err != nil {
				h.LogErrorMessage("handlers.fizzbuzz", err, "Error retrieving data from cache")
				return h.RespondJSONBadRequest()
			}

			result := &models.Result{}
			err = json.Unmarshal([]byte(cache), result)
			if err != nil {
				h.LogErrorMessage("handlers.fizzbuzz", err, "Error decoding struct from json")
				return h.RespondJSONBadRequest()
			}
			result.Flag = "cached"

			stats = append(stats, Stats{Hits: v.Score, Result: result})
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": RESPONSE_STATUS_SUCCESS,
		"data":   stats,
	})
}
