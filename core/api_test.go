package core

import (
	"os"
	"testing"

	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
)

func TestGetAPIVersionFromFile(t *testing.T) {
	tests := map[string]struct {
		file         string
		stringLenght int
		expected     interface{}
	}{
		"case1": {
			file:         "BADFILE",
			stringLenght: 0,
			expected:     &os.PathError{},
		},

		"case2": {
			file:         "../VERSION",
			stringLenght: 5,
			expected:     nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			version, err := GetAPIVersionFromFile(tc.file)
			assert.IsType(t, tc.expected, err)
			assert.Equal(t, len(version), tc.stringLenght)
		})
	}
}

func TestGetLoggerConfig(t *testing.T) {
	loggerConfig := GetLoggerConfig()

	assert.IsType(t, loggerConfig, middleware.LoggerConfig{})
	assert.NotEmpty(t, loggerConfig.Output)
}

func TestGetCORSConfig(t *testing.T) {
	CORSConfig := GetCORSConfig()

	assert.IsType(t, CORSConfig, middleware.CORSConfig{})
	assert.GreaterOrEqual(t, len(CORSConfig.AllowOrigins), 1)
	assert.GreaterOrEqual(t, len(CORSConfig.AllowMethods), 3)
	assert.GreaterOrEqual(t, len(CORSConfig.AllowHeaders), 3)
	assert.NotEmpty(t, CORSConfig.MaxAge)
}

func TestGetSecureConfig(t *testing.T) {
	secureConfig := GetSecureConfig()

	assert.IsType(t, secureConfig, middleware.SecureConfig{})
	assert.NotEmpty(t, secureConfig.XFrameOptions)
	assert.NotEmpty(t, secureConfig.ContentSecurityPolicy)
	assert.NotEmpty(t, secureConfig.HSTSMaxAge)
}
