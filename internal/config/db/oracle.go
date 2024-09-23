package db

import (
	"fmt"
	"os"
	"strconv"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	oracleHostEnvName     = "ORACLE_HOST"
	oraclePortEnvName     = "ORACLE_PORT"
	oracleUsernameEnvName = "ORACLE_USERNAME"
	oraclePasswordEnvName = "ORACLE_PASSWORD"
	oracleServiceEnvName  = "ORACLE_SERVICE"
)

type oracleConfig struct {
	host     string
	port     int
	username string
	password string
	service  string
}

func (cfg *oracleConfig) setFromEnv() error {
	var err error
	cfg.host = os.Getenv(oracleHostEnvName)
	if len(cfg.host) == 0 {
		return fmt.Errorf("oracleConfig host not found")
	}

	portStr := os.Getenv(oraclePortEnvName)
	if len(portStr) == 0 {
		return fmt.Errorf("oracleConfig port not found")
	}

	cfg.port, err = strconv.Atoi(portStr)
	if err != nil {
		return errors.Wrap(err, "failed to parse oracleConfig port")
	}

	cfg.username = os.Getenv(oracleUsernameEnvName)
	if len(cfg.username) == 0 {
		return fmt.Errorf("oracleConfig username not found")
	}

	cfg.password = os.Getenv(oraclePasswordEnvName)
	if len(cfg.password) == 0 {
		return fmt.Errorf("oracleConfig password not found")
	}

	cfg.service = os.Getenv(oracleServiceEnvName)
	if len(cfg.service) == 0 {
		return fmt.Errorf("oracleConfig sevice not found")
	}

	return nil
}

func NewOracleConfig() (*oracleConfig, error) {
	cfg := &oracleConfig{}
	if err := cfg.setFromEnv(); err != nil {
		return nil, fmt.Errorf("oracleConfig set from env err: %s", err)
	}
	logrus.Infof("oracleConfig: [%+v]", *cfg)

	return cfg, nil
}

func (cfg *oracleConfig) DSN() string {
	return fmt.Sprintf(`user=%s password=%s connectString=%s:%d/%s`, cfg.username, cfg.password, cfg.host, cfg.port, cfg.service)
}
