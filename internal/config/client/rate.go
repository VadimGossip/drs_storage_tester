package client

import (
	"fmt"
	"os"
	"strconv"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	rateGRPCHostEnvName = "RATE_GRPC_SERVER_HOST"
	rateGRPCPortEnvName = "RATE_GRPC_SERVER_PORT"
)

type rateGRPCClientConfig struct {
	host string
	port int
}

func (cfg *rateGRPCClientConfig) setFromEnv() error {
	var err error
	cfg.host = os.Getenv(rateGRPCHostEnvName)
	if len(cfg.host) == 0 {
		return fmt.Errorf("rateGRPCClientConfig host not found")
	}

	portStr := os.Getenv(rateGRPCPortEnvName)
	if len(portStr) == 0 {
		return fmt.Errorf("rateGRPCClientConfig port not found")
	}

	cfg.port, err = strconv.Atoi(portStr)
	if err != nil {
		return errors.Wrap(err, "failed to parse rateGRPCClientConfig port")
	}
	return nil
}

func NewRateGRPCConfig() (*rateGRPCClientConfig, error) {
	cfg := &rateGRPCClientConfig{}
	if err := cfg.setFromEnv(); err != nil {
		return nil, fmt.Errorf("rateGRPCClientConfig set from env err: %s", err)
	}

	logrus.Infof("rateGRPCClientConfig: [%+v]", *cfg)
	return cfg, nil
}

func (cfg *rateGRPCClientConfig) Address() string {
	return fmt.Sprintf("%s:%d", cfg.host, cfg.port)
}
