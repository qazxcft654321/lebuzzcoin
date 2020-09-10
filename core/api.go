package core

import (
	"bufio"
	"os"

	"github.com/labstack/echo/v4/middleware"
)

func GetAPIVersionFromFile(file string) (string, error) {
	var version string
	if _, err := os.Stat(file); err != nil {
		return version, err
	}

	f, err := os.Open(file)
	if err != nil {
		return version, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		version = scanner.Text()
	}

	err = scanner.Err()
	return version, err
}

func GetLoggerConfig() middleware.LoggerConfig {
	return middleware.LoggerConfig{
		Output: os.Stdout,
		Format: "time=${time_rfc3339}, status=${status}, method=${method}, uri=${uri}, remote_ip=${remote_ip}, userAgent=${user_agent}, latency=${latency_human}\n",
	}
}

func GetCORSConfig() middleware.CORSConfig {
	return middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "OPTION"},
		AllowHeaders: []string{"Accept", "Content-Type", "X-Requested-With"},
		MaxAge:       3600,
	}
}

func GetSecureConfig() middleware.SecureConfig {
	return middleware.SecureConfig{
		XFrameOptions:         "DENY",
		HSTSMaxAge:            3600,
		ContentSecurityPolicy: "frame-ancestors 'none'",
	}
}
