package main

import (
	"os"

	"lebuzzcoin/core/api"
	"lebuzzcoin/core/cache"
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

	// API
	err = api.Init()
	if err != nil {
		e.Logger.Fatalf("Error retrieving API version from file: %v \n", err)
	}

	// Redis
	cache, err := cache.NewCache(cache.Options{Address: os.Getenv("REDIS_ADDR")})
	if err != nil {
		e.Logger.Fatalf("Cannot ping: %v \n", err)
	}

	// Setup some middlewares at router level
	e.Use(middleware.LoggerWithConfig(api.GetLoggerConfig()))
	e.Use(middlewares.LimitMiddleware(tollbooth.NewLimiter(2, nil))) // NOTE: set limit at 2/s (hardcore mode)
	e.Use(middleware.CORSWithConfig(api.GetCORSConfig()))
	e.Use(middleware.SecureWithConfig(api.GetSecureConfig()))
	e.Use(middleware.BodyLimit("3M"))

	// Setup groups
	v1 := e.Group("/v1")

	// Setup handler struct
	h := handlers.New(os.Stdout, cache)

	// Routes
	e.GET("/", h.GetAPIVersion)
	v1.POST("/fizzbuzz/compute", h.ComputeFizzbuzz) // NOTE: due to many parameters I am going for POST (dont't like overstuffing URLs)
	// TODO: v1/fizzbuzz/stats
	// TODO: v1/fizzbuzz/:hash (retrieve specific result)

	// Server
	e.Logger.Fatal(e.StartTLS(":"+os.Getenv("HTTP_PORT"), os.Getenv("CERT_PEM"), os.Getenv("KEY_PEM")))
}
