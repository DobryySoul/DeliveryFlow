package nats

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/DobryySoul/DeliveryFlow/api-gateway/internal/config"
	"github.com/nats-io/nats.go"
)

var (
	errInvalidConnection = errors.New("nats connection is not initialized")
	errInvalidSubject    = errors.New("invalid nats subject")
)

type Client struct {
	nc *nats.Conn
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Connect(ctx context.Context, cfg *config.NatsConfig) error {
	opts := []nats.Option{
		nats.Timeout(cfg.ConnectTimeout),
		nats.ReconnectWait(cfg.ReconnectTimeout),
		nats.MaxReconnects(-1),
		nats.ReconnectHandler(func(c *nats.Conn) {
			log.Println("Reconnected to", c.ConnectedUrl())
		}),
		nats.DisconnectHandler(func(c *nats.Conn) {
			log.Println("Disconnected from NATS")
		}),
		nats.ClosedHandler(func(c *nats.Conn) {
			log.Println("NATS connection is closed.")
		}),
	}

	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		nc, err := nats.Connect(cfg.URL, opts...)
		if err == nil {
			c.nc = nc

			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(cfg.RetryWait):
		}
	}
}

func (c *Client) Request(ctx context.Context, subject string, data []byte) ([]byte, error) {
	if c.nc == nil {
		return nil, errInvalidConnection
	}

	if !IsValidSubject(subject) {
		return nil, errInvalidSubject
	}

	msg, err := c.nc.RequestWithContext(ctx, subject, data)
	if err != nil {
		return nil, fmt.Errorf("failed to request NATS: %w", err)
	}

	return msg.Data, nil
}

func (c *Client) IsConnected() bool {
	return c.nc != nil && c.nc.IsConnected()
}

func (c *Client) Close() error {
	if c.nc == nil {
		return nil
	}

	c.nc.Close()

	return nil
}
