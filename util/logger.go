package util

import (
	"os"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/henilmalaviya/golw/env"
)

var logger *log.Logger = nil

func GetLogger() *log.Logger {
	return logger
}

// parseLogLevel converts string log level to log.Level
func parseLogLevel(level string) log.Level {
	switch strings.ToLower(level) {
	case "trace":
		return log.DebugLevel // charmbracelet/log doesn't have Trace, use Debug
	case "debug":
		return log.DebugLevel
	case "info":
		return log.InfoLevel
	case "warn", "warning":
		return log.WarnLevel
	case "error":
		return log.ErrorLevel
	case "fatal":
		return log.FatalLevel
	default:
		return log.InfoLevel // default fallback
	}
}

func init() {
	logger = log.NewWithOptions(os.Stdout, log.Options{
		ReportTimestamp: true,
		ReportCaller:    true,
	})

	// Get log level from environment configuration
	logLevel := parseLogLevel(env.Get().LogLevel)
	logger.SetLevel(logLevel)
}
