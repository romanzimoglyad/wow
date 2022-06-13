package logger

import (
	"github.com/romanzimoglyad/wow/internal/config"
	"github.com/rs/zerolog"
)

// InitZeroLog initializes the logger for the application to work
func InitZeroLog(cfg *config.Config) {
	var (
		err      error
		logLevel zerolog.Level
	)

	if logLevel, err = zerolog.ParseLevel(cfg.LogLevel); err != nil {
		logLevel = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(logLevel)
}
