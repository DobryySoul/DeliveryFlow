package dto

import (
	"errors"
	"time"
)

type OrderStatus string

const (
	StatusCreated   OrderStatus = "created"
	StatusReserved  OrderStatus = "reserved"
	StatusPaid      OrderStatus = "paid"
	StatusAssigned  OrderStatus = "assigned"
	StatusDelivered OrderStatus = "delivered"
	StatusCancelled OrderStatus = "cancelled"
)

type OrderResponse struct {
	UpdatedAt time.Time   `json:"updated_at"`
	Status    OrderStatus `json:"status"`
	Err       string      `json:"error"`
	OrderID   string      `json:"order_id"`
}

func (o *OrderResponse) Error() error {
	return errors.New(o.Err)
}

func (o *OrderResponse) IsNotFound() bool {
	if o.Err == "" {
		return false
	}

	return o.Err == ErrOrderNotFound.Error()
}

func (o *OrderResponse) IsFailedToRequestNATS() bool {
	if o.Err == "" {
		return false
	}

	return o.Err == ErrFailedToRequestNATS.Error()
}

var (
	ErrFailedToUnmarshalOrderResponse = errors.New("failed to unmarshal order response")
	ErrOrderNotFound                  = errors.New("order not found")
	ErrFailedToRequestNATS            = errors.New("failed to request NATS")
	ErrTimeout                        = errors.New("timeout")
)
