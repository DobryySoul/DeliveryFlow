package app

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/DobryySoul/DeliveryFlow/api-gateway/internal/config"
	"github.com/DobryySoul/DeliveryFlow/api-gateway/internal/http"
	natsrpc "github.com/DobryySoul/DeliveryFlow/api-gateway/internal/nats"
	"github.com/DobryySoul/DeliveryFlow/api-gateway/internal/usecase"

	"github.com/rs/zerolog"
)

func Run(ctx context.Context, logger *zerolog.Logger, cfg *config.Config) error {
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger.Info().Msg("starting gateway service")

	natsClient := natsrpc.NewClient()
	if err := natsClient.Connect(ctx, cfg.NATSConfig); err != nil {
		logger.Error().Err(err).Msgf("failed to connect to NATS: %s", cfg.NATSConfig)
		return fmt.Errorf("failed to connect to NATS: %w", err)
	}

	orderUC := usecase.NewOrderUsecase(natsClient, logger)

	httpServer := http.NewServer(cfg.HTTPConfig)
	router := http.NewRouter(httpServer.App, orderUC, natsClient, logger)
	router.RegisterRoutes()

	go func() {
		if err := httpServer.Start(ctx, cfg); err != nil {
			logger.Error().Err(err).Msg("failed to start API Gateway HTTP server")
		}
	}()

	logger.Info().Str("addr", cfg.HTTPConfig.HTTPAddr).Msg("API Gateway service started successfully")

	<-ctx.Done()

	logger.Info().Msg("shutting down API Gateway service")

	if err := httpServer.Stop(ctx); err != nil {
		logger.Error().Err(err).Msg("failed to stop API Gateway HTTP server")
	}

	if err := natsClient.Close(); err != nil {
		logger.Error().Err(err).Msg("failed to close NATS client")
	}

	return nil
}
