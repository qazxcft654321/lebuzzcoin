package handlers

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

func (*Handler) GetAPIVersion(c echo.Context) error {
	APIVersion := os.Getenv("APIVERSION")
	if len(APIVersion) != 5 {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{
			"status":  "error",
			"message": "Server cannot process the request",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Lebuzzcoin-API v" + APIVersion,
	})
}
