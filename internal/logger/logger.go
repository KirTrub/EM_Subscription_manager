package logger

import (
	"io"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

func New(level string) zerolog.Logger {
	lvl, err := zerolog.ParseLevel(strings.ToLower(level))
	if err != nil {
		lvl = zerolog.InfoLevel
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	w := io.MultiWriter(os.Stdout)
	return zerolog.New(w).Level(lvl).With().Timestamp().Logger()
}
