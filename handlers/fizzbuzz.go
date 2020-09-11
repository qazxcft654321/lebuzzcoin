package handlers

import (
	"net/http"
	"os"

	"lebuzzcoin/models"

	"github.com/labstack/echo/v4"
)

func (h *Handler) GetAPIVersion(c echo.Context) error {
	APIVersion := os.Getenv("APIVERSION")
	if len(APIVersion) != 5 {
		h.LogErrorMessage("handlers.fizzbuzz", nil, "Error retrieving API version from env variable")
		return h.RespondJSONBadRequest()
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  RESPONSE_STATUS_SUCCESS,
		"message": "Lebuzzcoin-API v" + APIVersion,
	})
}

func (h *Handler) ComputeFizzbuzz(c echo.Context) error {
	fizzbuzz := &models.Fizzbuzz{}

	// Decoding request body
	if err := c.Bind(fizzbuzz); err != nil {
		h.LogErrorMessage("handlers.fizzbuzz", err, "Error decoding request body")
		return h.RespondJSONBadRequest()
	}

	// Struct validation
	if err := h.validator.Struct(fizzbuzz); err != nil {
		h.LogErrorMessage("handlers.fizzbuzz", err, "Data validation failed")
		return h.RespondJSONForbidden()
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": RESPONSE_STATUS_SUCCESS,
		"data":   fizzbuzz,
	})
}
