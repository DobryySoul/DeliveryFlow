package main

import (
	"context"

	"github.com/DobryySoul/DeliveryFlow/api-gateway/internal/app"
	"github.com/DobryySoul/DeliveryFlow/api-gateway/internal/config"
	"github.com/DobryySoul/DeliveryFlow/api-gateway/observability"
)

func main() {
	ctx := context.Background()
	logger := observability.NewLogger()

	cfg, err := config.Load()
	if err != nil {
		logger.Fatal().Err(err).Msg("config load failed")
	}

	logger.Info().Msg("config loaded successfully")

	if err := app.Run(ctx, logger, cfg); err != nil {
		logger.Fatal().Err(err).Msg("gateway run failed")
	}
}
