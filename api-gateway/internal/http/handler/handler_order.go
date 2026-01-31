package handler

import (
	"errors"

	"github.com/DobryySoul/DeliveryFlow/api-gateway/internal/dto"
	"github.com/DobryySoul/DeliveryFlow/api-gateway/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type HandlerOrder struct {
	orderUC *usecase.OrderUsecase
	logger  *zerolog.Logger
}

func NewHandlerOrder(orderUC *usecase.OrderUsecase, logger *zerolog.Logger) *HandlerOrder {
	return &HandlerOrder{orderUC: orderUC, logger: logger}
}

func (h *HandlerOrder) CreateOrder(c *fiber.Ctx) error {
	h.logger.Info().Msg("Creating order")

	var order dto.OrderRequest
	if err := c.BodyParser(&order); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	orderID, status, err := h.orderUC.CreateOrder(c.Context(), &order)
	if err != nil {
		// TODO: switch errors and return appropriate status code
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   err.Error(),
			"message": "failed to create order",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"order_id": orderID,
		"status":   status,
	})
}

func (h *HandlerOrder) GetOrderStatus(c *fiber.Ctx) error {
	h.logger.Info().Msg("Getting order status")

	orderID := c.Params("id")
	if orderID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "order ID is required",
		})
	}

	order, err := h.orderUC.GetOrderStatus(c.Context(), orderID)
	if err != nil {
		switch {
		case order != nil && order.IsNotFound():
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   order.Err,
				"message": "failed to get order status",
			})
		case errors.Is(err, dto.ErrFailedToRequestNATS):
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error":   err.Error(),
				"message": "failed to get request to order service",
			})
		case errors.Is(err, dto.ErrTimeout):
			return c.Status(fiber.StatusGatewayTimeout).JSON(fiber.Map{
				"error":   err.Error(),
				"message": "failed to get request to order service",
			})
		case errors.Is(err, dto.ErrFailedToUnmarshalOrderResponse):
			return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
				"error":   err.Error(),
				"message": "invalid response from order service",
			})
		default:
			return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
				"error":   err.Error(),
				"message": "failed to get order",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"order_id":   order.OrderID,
		"status":     order.Status,
		"updated_at": order.UpdatedAt,
	})
}
