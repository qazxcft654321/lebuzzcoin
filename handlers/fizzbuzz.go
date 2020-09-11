package handlers

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

func (h *Handler) GetAPIVersion(c echo.Context) error {
	APIVersion := os.Getenv("APIVERSION")
	if len(APIVersion) != 5 {
		h.LogErrorMessage("Error retrieving API version from env variable", "handlers.fizzbuzz", nil)
		return h.RespondJSONBadRequest()
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  RESPONSE_STATUS_SUCCESS,
		"message": "Lebuzzcoin-API v" + APIVersion,
	})
}
