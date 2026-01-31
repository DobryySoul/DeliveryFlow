package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	defaultHTTPAddr = ":8080"
	defaultNATSURL  = "nats_rpc:4222"
)

type Config struct {
	HTTPConfig *HTTPConfig `yaml:"http_config" env:"HTTP_CONFIG" env-default:":8080"`
	NATSConfig *NatsConfig `yaml:"nats" env-prefix:"NATS_"`
}

type HTTPConfig struct {
	ReadTimeout  time.Duration `yaml:"read_timeout" env:"READ_TIMEOUT" env-default:"10s"`
	WriteTimeout time.Duration `yaml:"write_timeout" env:"WRITE_TIMEOUT" env-default:"10s"`
	IdleTimeout  time.Duration `yaml:"idle_timeout" env:"IDLE_TIMEOUT" env-default:"10s"`
	HTTPAddr     string        `yaml:"http_addr" env:"HTTP_ADDR" env-default:":8080"`
	Concurrency  int           `yaml:"concurrency" env:"CONCURRENCY" env-default:"1000"`
	Prefork      bool          `yaml:"prefork" env:"PREFORK" env-default:"false"`
}
type NatsConfig struct {
	ConnectTimeout   time.Duration `yaml:"connect_timeout" env:"CONNECT_TIMEOUT" env-default:"5s"`
	ReconnectTimeout time.Duration `yaml:"reconnect_timeout" env:"RECONNECT_TIMEOUT" env-default:"2s"`
	RetryWait        time.Duration `yaml:"retry_wait" env:"RETRY_WAIT" env-default:"1s"`
	URL              string        `yaml:"url" env:"NATS_URL" env-default:"nats_rpc:4222"`
}

func Load() (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadConfig("./internal/config/config.yaml", &cfg); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	return &cfg, nil
}
