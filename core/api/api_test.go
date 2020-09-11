package api

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
)

func TestGetAPIVersionFromFile(t *testing.T) {
	tests := map[string]struct {
		version  string
		expected interface{}
	}{
		"case1": {
			expected: &os.PathError{},
		},

		"case2": {
			version:  "6.6.6",
			expected: nil,
		},
	}

	var vFile string = "./test_api_version"

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if len(tc.version) > 1 {
				vByte := []byte(tc.version)
				err := ioutil.WriteFile(vFile, vByte, 0644)
				assert.NoError(t, err)
			}

			version, err := getAPIVersionFromFile(vFile)
			assert.IsType(t, tc.expected, err)
			assert.Equal(t, version, tc.version)

			_ = os.Remove(vFile)
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
