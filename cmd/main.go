// @title EM Subscription Manager API
// @version 1.0
// @description Subscription management service
// @host localhost:3000
// @BasePath /
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	appcfg "em_subscription_manager/internal/config"
	"em_subscription_manager/internal/db"
	"em_subscription_manager/internal/handlers"
	"em_subscription_manager/internal/handlers/middleware"
	"em_subscription_manager/internal/logger"
	"em_subscription_manager/internal/repo"
	"em_subscription_manager/internal/services"

	swaggerui "github.com/gofiber/contrib/v3/swaggerui"
	"github.com/gofiber/fiber/v3"
)

func main() {
	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		cfgPath = "configs/config.yaml"
	}

	cfg, err := appcfg.Load(cfgPath)
	if err != nil {
		panic(err)
	}

	log := logger.New(cfg.Logger.Level)
	log.Info().Str("app", cfg.App.Name).Msg("config loaded")

	database, err := db.Connect(cfg.Postgres.DSN)
	if err != nil {
		log.Fatal().Err(err).Msg("database connection failed")
	}
	defer database.Close()

	err = db.RunMigrations(database)
	if err != nil {
		log.Fatal().Err(err).Msg("database migrations failed")
	}
	log.Info().Msg("database migrations completed")

	subscriptionRepo := repo.New(database)
	subscriptionService := services.NewSubscriptionService(subscriptionRepo)

	app := fiber.New(fiber.Config{
		AppName: cfg.App.Name,
	})

	app.Use(middleware.RequestLogger(log))

	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	handlers.RegisterSubscriptionRoutes(app, subscriptionService)

	if cfg.Swagger.Enabled {
		app.Use(swaggerui.New(swaggerui.Config{
			BasePath: "/",
			FilePath: "./docs/swagger.json",
			Path:     "swagger",
			Title:    "EM Subscription Manager API",
		}))
	}

	addr := fmt.Sprintf(":%d", cfg.App.Port)

	go func() {
		log.Info().Str("addr", addr).Msg("server started")
		if err := app.Listen(addr); err != nil {
			log.Fatal().Err(err).Msg("server stopped unexpectedly")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Error().Err(err).Msg("graceful shutdown failed")
	}
	log.Info().Msg("shutdown complete")
}
