package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/DobryySoul/DeliveryFlow/api-gateway/internal/dto"
	"github.com/DobryySoul/DeliveryFlow/api-gateway/internal/nats"
	"github.com/rs/zerolog"
)

type OrderUsecase struct {
	natsClient *nats.Client
	logger     *zerolog.Logger
}

func NewOrderUsecase(natsClient *nats.Client, logger *zerolog.Logger) *OrderUsecase {
	return &OrderUsecase{
		natsClient: natsClient,
		logger:     logger,
	}
}

func (uc *OrderUsecase) CreateOrder(ctx context.Context, order *dto.OrderRequest) (string, dto.OrderStatus, error) {
	uc.logger.Info().Msg("Creating order")

	orderData, err := json.Marshal(order)
	if err != nil {
		uc.logger.Error().Err(err).Msg("failed to marshal order")
		return "", "", fmt.Errorf("failed to marshal order: %w", err)
	}

	response, err := uc.natsClient.Request(ctx, nats.SubjectCreateOrder, orderData)
	if err != nil {
		uc.logger.Error().Err(err).Msg("failed to request NATS")
		return "", "", fmt.Errorf("failed to request NATS: %w", err)
	}

	var orderResponse dto.OrderResponse
	if err := json.Unmarshal(response, &orderResponse); err != nil {
		uc.logger.Error().Err(err).Msg("failed to unmarshal order response")
		return "", "", fmt.Errorf("failed to unmarshal order response: %w", err)
	}

	if orderResponse.Status != dto.StatusCreated {
		uc.logger.Error().Msgf("failed to create order: %s", orderResponse.Status)
		return "", "", fmt.Errorf("failed to create order: %s", orderResponse.Status)
	}

	return orderResponse.OrderID, orderResponse.Status, nil
}

func (uc *OrderUsecase) GetOrderStatus(ctx context.Context, orderID string) (*dto.OrderResponse, error) {
	response, err := uc.natsClient.Request(ctx, nats.SubjectGetOrderStatus, []byte(orderID))
	if err != nil {
		uc.logger.Error().Err(err).Msg("failed to request NATS")
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, dto.ErrTimeout
		}
		return nil, fmt.Errorf("%w: %v", dto.ErrFailedToRequestNATS, err)
	}

	var orderResponse dto.OrderResponse
	if err := json.Unmarshal(response, &orderResponse); err != nil {
		return nil, fmt.Errorf("%w: %v", dto.ErrFailedToUnmarshalOrderResponse, err)
	}

	if orderResponse.Err != "" {
		return &orderResponse, orderResponse.Error()
	}

	return &orderResponse, nil
}
