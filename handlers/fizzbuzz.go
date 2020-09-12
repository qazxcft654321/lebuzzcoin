package handlers

import (
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"

	dbCache "lebuzzcoin/core/cache"
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
	if err != nil && err.Error() != dbCache.Empty {
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
	buffer, err := hex.DecodeString(hash)

	// Check if hex & 256bits / 32Bytes long
	if err != nil || len(buffer) != 32 {
		h.LogErrorMessage("handlers.fizzbuzz", nil, "Submited hash incorrect")
		return h.RespondJSONForbidden()
	}

	// Retrive from cache
	cache, err := h.cache.Get(hash)
	if err != nil || err.Error() == dbCache.Empty {
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
