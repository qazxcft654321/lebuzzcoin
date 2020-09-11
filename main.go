package main

import (
	"fmt"
	"net/http"
	"os"

	"lebuzzcoin/core"
	"lebuzzcoin/handlers"
	"lebuzzcoin/middlewares"

	"github.com/didip/tollbooth"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	err := godotenv.Load()
	if err != nil {
		e.Logger.Fatal("Error loading env file")
	}

	err = core.InitAPI()
	if err != nil {
		e.Logger.Fatalf("Error retrieving API version from file: %v \n", err)
	}
	fmt.Println("Lebuzzcoin-API v" + os.Getenv("APIVERSION"))

	// Setup some middlewares at router level
	e.Use(middleware.LoggerWithConfig(core.GetLoggerConfig()))
	e.Use(middlewares.LimitMiddleware(tollbooth.NewLimiter(1, nil))) // NOTE: set limit at 3/s
	e.Use(middleware.CORSWithConfig(core.GetCORSConfig()))
	e.Use(middleware.SecureWithConfig(core.GetSecureConfig()))
	e.Use(middleware.BodyLimit("3M"))

	// Setup groups
	v1 := e.Group("/v1")

	// Setup handler struct
	h := handlers.New(os.Stdout)

	// TODO: move to handler
	e.GET("/", h.GetAPIVersion)
	v1.POST("/fizzbuzz", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":  "success",
			"message": "...",
		})
	})

	e.Logger.Fatal(e.StartTLS(":"+os.Getenv("HTTP_PORT"), os.Getenv("CERT_PEM"), os.Getenv("KEY_PEM")))
}
