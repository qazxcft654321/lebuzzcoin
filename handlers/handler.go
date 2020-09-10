package handlers

import (
	"io"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type Handler struct {
	logger echo.Logger
}

func New(logOutput io.Writer) *Handler {
	logger := echo.New().Logger
	logger.SetOutput(logOutput)
	logger.SetLevel(log.INFO)
	logger.SetHeader("time=${time_rfc3339}, level=${level}, file=${long_file}, line=${line}, message=${message}")

	return &Handler{logger: logger}
}
