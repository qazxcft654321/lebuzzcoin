package main

import (
	"bufio"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"os"
)

var HTTPPort string = "8080"

// TODO: move to something like helper later + better name convention?
func GetApiVersion(versionFile string) (string, error) {
	var version string
	if _, err := os.Stat(versionFile); err != nil {
		return version, err
	}

	file, err := os.Open(versionFile)
	if err != nil {
		return version, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		version = scanner.Text()
	}

	err = scanner.Err()
	return version, err
}

func main() {
	e := echo.New()

	APIVersion, err := GetApiVersion("VERSION")
	if err != nil {
		log.Fatalf("error while retrieving API version from file: %v", err)
	}

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":  "success",
			"message": "Lebuzzcoin-API v" + APIVersion,
		})
	})

	log.Fatal(e.Start(":" + HTTPPort))
}
