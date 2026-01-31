package http

import (
	"context"

	"github.com/DobryySoul/DeliveryFlow/api-gateway/internal/config"
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	App *fiber.App
}

func NewServer(cfg *config.HTTPConfig) *Server {
	app := fiber.New(fiber.Config{
		Prefork:      cfg.Prefork,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
		Concurrency:  cfg.Concurrency,
	})

	return &Server{
		App: app,
	}
}

func (s *Server) Start(ctx context.Context, cfg *config.Config) error {
	return s.App.Listen(cfg.HTTPConfig.HTTPAddr)
}

func (s *Server) Stop(ctx context.Context) error {
	return s.App.Shutdown()
}
