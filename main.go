package main

import (
	"log"
	"net/http"
	"os"

	"lebuzzcoin/core"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading env file")
	}

	APIVersion, err := core.GetAPIVersionFromFile("VERSION")
	if err != nil {
		log.Fatalf("error while retrieving API version from file: %v", err)
	}

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":  "success",
			"message": "Lebuzzcoin-API v" + APIVersion,
		})
	})

	log.Fatal(e.Start(":" + os.Getenv("HTTP_PORT")))
}
