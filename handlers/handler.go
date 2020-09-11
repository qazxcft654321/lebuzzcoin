package handlers

import (
	"io"
	"net/http"
	"runtime"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

const (
	RESPONSE_STATUS_ERROR              string = "error"
	RESPONSE_STATUS_SUCCESS            string = "success"
	RESPONSE_MESSAGE_BAD_REQUEST       string = "Server cannot process the request"
	RESPONSE_MESSAGE_RESSOURCE_CREATED string = "Ressource created"
)

type Handler struct {
	logger echo.Logger
}

func New(logOutput io.Writer) *Handler {
	logger := echo.New().Logger
	logger.SetOutput(logOutput)
	logger.SetLevel(log.INFO)
	logger.SetHeader("time=${time_rfc3339}, level=${level}, message=${message}")

	return &Handler{logger: logger}
}

func (h *Handler) LogErrorMessage(message, ressource string, err error) {
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
