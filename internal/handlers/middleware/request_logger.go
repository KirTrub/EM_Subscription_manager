package middleware

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog"
)

func RequestLogger(log zerolog.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()
		err := c.Next()

		event := log.Info().
			Str("method", c.Method()).
			Str("path", c.Path()).
			Int("status", c.Response().StatusCode()).
			Str("ip", c.IP()).
			Dur("latency", time.Since(start))

		if err != nil {
			event = log.Error().Err(err).Str("method", c.Method()).Str("path", c.Path())
		}

		event.Msg("http request")
		return err
	}
}
