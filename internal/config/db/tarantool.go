package db

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	tdbHostEnvName     = "TARANTOOL_HOST"
	tdbPortEnvName     = "TARANTOOL_PORT"
	tdbUsernameEnvName = "TARANTOOL_USERNAME"
	tdbPasswordEnvName = "TARANTOOL_PASSWORD"
	tdbTimeoutSec      = "TARANTOOL_TIMEOUT_SEC"
)

type tdbConfig struct {
	host       string
	port       int
	username   string
	password   string
	timeoutSec int
}

func (cfg *tdbConfig) setFromEnv() error {
	var err error
	cfg.host = os.Getenv(tdbHostEnvName)
	if len(cfg.host) == 0 {
		return fmt.Errorf("tdbConfig host not found")
	}

	portStr := os.Getenv(tdbPortEnvName)
	if len(portStr) == 0 {
		return fmt.Errorf("tdbConfig port not found")
	}

	cfg.port, err = strconv.Atoi(portStr)
	if err != nil {
		return errors.Wrap(err, "failed to parse tdbConfig port")
	}

	cfg.username = os.Getenv(tdbUsernameEnvName)
	cfg.password = os.Getenv(tdbPasswordEnvName)

	timeoutSecStr := os.Getenv(tdbTimeoutSec)
	if len(timeoutSecStr) == 0 {
		return fmt.Errorf("tdbConfig timeout sec not found")
	}

	cfg.timeoutSec, err = strconv.Atoi(timeoutSecStr)
	if err != nil {
		return errors.Wrap(err, "failed to parse tdbConfig timeout sec")
	}

	return nil
}

func NewTarantoolConfig() (*tdbConfig, error) {
	cfg := &tdbConfig{}
	if err := cfg.setFromEnv(); err != nil {
		return nil, fmt.Errorf("tdbConfig set from env err: %s", err)
	}
	logrus.Infof("tdbConfig: [%+v]", *cfg)
	return cfg, nil
}

func (cfg *tdbConfig) Address() string {
	return fmt.Sprintf("%s:%d", cfg.host, cfg.port)
}

func (cfg *tdbConfig) Username() string {
	return cfg.username
}

func (cfg *tdbConfig) Password() string {
	return cfg.password
}

func (cfg *tdbConfig) Timeout() time.Duration {
	return time.Duration(cfg.timeoutSec) * time.Second
}
