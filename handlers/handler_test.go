package handlers

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
)

func TestNewHandler(t *testing.T) {
	h := New(os.Stdout)

	assert.IsType(t, h, &Handler{})
	assert.NotEmpty(t, h)
}

func TestLogErrorMessage(t *testing.T) {
	var buffer bytes.Buffer
	logger := echo.New().Logger
	logger.SetOutput(&buffer)
	logger.SetLevel(log.INFO)
	logger.SetHeader("time=${time_rfc3339}, level=${level}, message=${message}")

	h := &Handler{logger: logger}
	h.LogErrorMessage("test ressource", errors.New("test error"), "test message")

	ressource := bytes.ContainsAny(buffer.Bytes(), "test ressource")
	error := bytes.ContainsAny(buffer.Bytes(), "test error")
	message := bytes.ContainsAny(buffer.Bytes(), "test message")

	assert.True(t, ressource)
	assert.True(t, error)
	assert.True(t, message)
}

func TestRespondJSONBadRequest(t *testing.T) {
	e := echo.New()
	fakeHTTPHandler := func(c echo.Context) error {
		h := &Handler{}
		return h.RespondJSONBadRequest()
	}

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	e.GET("/test", fakeHTTPHandler)
	e.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRespondJSONForbidden(t *testing.T) {
	e := echo.New()
	fakeHTTPHandler := func(c echo.Context) error {
		h := &Handler{}
		return h.RespondJSONForbidden()
	}

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	e.GET("/test", fakeHTTPHandler)
	e.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}
