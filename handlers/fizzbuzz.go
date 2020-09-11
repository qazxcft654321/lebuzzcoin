package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"lebuzzcoin/models"

	"github.com/go-redis/redis"
	"github.com/labstack/echo/v4"
)

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

// TODO: make test
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

	// Retrive from cache
	cache, err := h.rdb.Get(hash).Result()
	if err != nil && err != redis.Nil {
		h.LogErrorMessage("handlers.fizzbuzz", err, "Error retrieving data from cache")
		return h.RespondJSONBadRequest()
	}

	// TODO: abstract redis to avoid dependencies
	// Building result from cache
	result := &models.Result{}
	if len(cache) > 1 {
		err := json.Unmarshal([]byte(cache), result)
		if err != nil {
			h.LogErrorMessage("handlers.fizzbuzz", err, "Error decoding struct from json")
			return h.RespondJSONBadRequest()
		}
		result.State = "cached"
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

		err = h.rdb.Set(hash, json, 0).Err()
		if err != nil {
			h.LogErrorMessage("handlers.fizzbuzz", err, "Error caching data")
			return h.RespondJSONBadRequest()
		}
		result.State = "built"
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": RESPONSE_STATUS_SUCCESS,
		"data":   result,
	})
}
