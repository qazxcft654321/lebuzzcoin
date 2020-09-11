package handlers

import (
	"io"
	"net/http"
	"runtime"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"gopkg.in/go-playground/validator.v9"
)

const (
	RESPONSE_STATUS_ERROR        string = "error"
	RESPONSE_STATUS_FAIL         string = "fail"
	RESPONSE_STATUS_SUCCESS      string = "success"
	RESPONSE_MESSAGE_BAD_REQUEST string = "Server cannot process the request"
	RESPONSE_MESSAGE_FORBIDDEN   string = "Invalid request"
)

type Handler struct {
	logger    echo.Logger
	validator *validator.Validate
	rdb       *redis.Client
}

func New(logOutput io.Writer, redisClient *redis.Client) *Handler {
	logger := echo.New().Logger
	logger.SetOutput(logOutput)
	logger.SetLevel(log.INFO)
	logger.SetHeader("time=${time_rfc3339}, level=${level}, message=${message}")

	return &Handler{
		logger:    logger,
		validator: validator.New(),
		rdb:       redisClient,
	}
}

func (h *Handler) LogErrorMessage(ressource string, err error, message string) {
	pc, _, line, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()

	if err != nil {
		if ressource != "" {
			h.logger.Errorf("[%s:%d] %s (%s) -> %v", funcName, line, message, ressource, err)
		} else {
			h.logger.Errorf("[%s:%d] %s -> %v", funcName, line, message, err)
		}
	} else {
		if ressource != "" {
			h.logger.Errorf("[%s:%d] %s (%s)", funcName, line, message, ressource)
		} else {
			h.logger.Error(message)
		}
	}
}

func (h *Handler) RespondJSONBadRequest() error {
	return echo.NewHTTPError(http.StatusBadRequest, map[string]string{
		"status":  RESPONSE_STATUS_ERROR,
		"message": RESPONSE_MESSAGE_BAD_REQUEST,
	})
}

func (h *Handler) RespondJSONForbidden() error {
	return echo.NewHTTPError(http.StatusForbidden, map[string]string{
		"status":  RESPONSE_STATUS_FAIL,
		"message": RESPONSE_MESSAGE_FORBIDDEN,
	})
}
