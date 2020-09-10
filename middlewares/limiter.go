package middlewares

import (
	"net/http"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/labstack/echo/v4"
)

func LimitMiddleware(lmt *limiter.Limiter) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return echo.HandlerFunc(func(c echo.Context) error {
			httpError := tollbooth.LimitByRequest(lmt, c.Response(), c.Request())
			if httpError != nil {
				return c.JSON(http.StatusTooManyRequests, map[string]string{
					"status":  "too many requests",
					"message": "You need to chill",
				})
			}
			return next(c)
		})
	}
}
