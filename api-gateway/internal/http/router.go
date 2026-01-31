package http

import (
	"github.com/DobryySoul/DeliveryFlow/api-gateway/internal/http/handler"
	natsrpc "github.com/DobryySoul/DeliveryFlow/api-gateway/internal/nats"
	"github.com/DobryySoul/DeliveryFlow/api-gateway/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type Router struct {
	app        *fiber.App
	orderUC    *usecase.OrderUsecase
	natsClient *natsrpc.Client
	logger     *zerolog.Logger
}

func NewRouter(app *fiber.App, orderUC *usecase.OrderUsecase, natsClient *natsrpc.Client, logger *zerolog.Logger) *Router {
	return &Router{
		app:        app,
		orderUC:    orderUC,
		natsClient: natsClient,
		logger:     logger,
	}
}

func (r *Router) RegisterRoutes() {
	orderHandler := handler.NewHandlerOrder(r.orderUC, r.logger)

	r.app.Get("/ready", func(c *fiber.Ctx) error {
		if !r.natsClient.IsConnected() {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error": "NATS is not connected",
			})
		}

		return c.SendString("OK")
	})

	api := r.app.Group("/api/v1")
	api.Post("/orders", orderHandler.CreateOrder)

	api.Get("/orders/:id", orderHandler.GetOrderStatus)

	api.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})
}
