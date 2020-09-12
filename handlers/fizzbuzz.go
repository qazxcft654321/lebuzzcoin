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
	cache, err := h.cache.Get(hash)
	if err != nil && err != redis.Nil {
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
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": RESPONSE_STATUS_SUCCESS,
		"data":   result,
	})
}

func (h *Handler) GetFizzbuzzFromHash(c echo.Context) error {
	hash := c.Param("hash")
	// TODO: switch to validator with hexa only?
	// TODO: hex.DecodeString() check if []byte is 256/32 long
	if len(hash) != 64 {
		h.LogErrorMessage("handlers.fizzbuzz", nil, "Submited hash incorrect")
		return h.RespondJSONForbidden()
	}

	// Retrive from cache
	cache, err := h.cache.Get(hash)
	if err != nil || err == redis.Nil {
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
