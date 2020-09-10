package main

import (
	"log"
	"net/http"
	"os"

	"lebuzzcoin/core"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading env file")
	}

	APIVersion, err := core.GetAPIVersionFromFile("VERSION")
	if err != nil {
		log.Fatalf("Error while retrieving API version from file: %v \n", err)
	}

	e := echo.New()
	e.Use(middleware.LoggerWithConfig(core.GetLoggerConfig()))
	e.Use(middleware.CORSWithConfig(core.GetCORSConfig()))
	e.Use(middleware.SecureWithConfig(core.GetSecureConfig()))
	e.Use(middleware.BodyLimit("3M"))

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":  "success",
			"message": "Lebuzzcoin-API v" + APIVersion,
		})
	})

	log.Fatal(e.StartTLS(":"+os.Getenv("HTTP_PORT"), os.Getenv("CERT_PEM"), os.Getenv("KEY_PEM")))
}
