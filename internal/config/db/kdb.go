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
	kdbHostEnvName         = "KDB_HOST"
	kdbPortEnvName         = "KDB_PORT"
	kdbUsernameEnvName     = "KDB_USERNAME"
	kdbPasswordEnvName     = "KDB_PASSWORD"
	kdbDBNameEnvName       = "KDB_DB"
	kdbReadTimeoutEnvName  = "KDB_READ_TIMEOUT_SEC"
	kdbWriteTimeoutEnvName = "KDB_WRITE_TIMEOUT_SEC"
)

type kdbConfig struct {
	host         string
	port         int
	username     string
	password     string
	db           int
	readTimeout  time.Duration
	writeTimeout time.Duration
}

func (cfg *kdbConfig) setFromEnv() error {
	var err error
	cfg.host = os.Getenv(kdbHostEnvName)
	if len(cfg.host) == 0 {
		return fmt.Errorf("kdbConfig host not found")
	}

	portStr := os.Getenv(kdbPortEnvName)
	if len(portStr) == 0 {
		return fmt.Errorf("kdbConfig port not found")
	}
	cfg.username = os.Getenv(kdbUsernameEnvName)
	cfg.password = os.Getenv(kdbPasswordEnvName)

	cfg.port, err = strconv.Atoi(portStr)
	if err != nil {
		return errors.Wrap(err, "failed to parse kdbConfig port")
	}

	dbStr := os.Getenv(kdbDBNameEnvName)
	if len(dbStr) == 0 {
		return fmt.Errorf("kdbConfig db not found")
	}

	cfg.db, err = strconv.Atoi(dbStr)
	if err != nil {
		return errors.Wrap(err, "failed to parse kdbConfig db")
	}

	readTimeoutStr := os.Getenv(kdbReadTimeoutEnvName)
	if len(readTimeoutStr) == 0 {
		return fmt.Errorf("kdbConfig read timeout not found")
	}

	readTimeoutSec, err := strconv.ParseInt(readTimeoutStr, 10, 64)
	if err != nil {
		return errors.Wrap(err, "failed to parse kdbConfig read timeout")
	}
	cfg.readTimeout = time.Duration(readTimeoutSec) * time.Second

	writeTimeoutStr := os.Getenv(kdbWriteTimeoutEnvName)
	if len(writeTimeoutStr) == 0 {
		return fmt.Errorf("kdb write timeout not found")
	}

	writeTimeoutSec, err := strconv.ParseInt(writeTimeoutStr, 10, 64)
	if err != nil {
		return errors.Wrap(err, "failed to parse kdbConfig write timeout")
	}
	cfg.writeTimeout = time.Duration(writeTimeoutSec) * time.Second

	return nil
}

func NewKdbConfig() (*kdbConfig, error) {
	cfg := &kdbConfig{}
	if err := cfg.setFromEnv(); err != nil {
		return nil, fmt.Errorf("kdbConfig set from env err: %s", err)
	}
	logrus.Infof("kdbConfig: [%+v]", *cfg)
	return cfg, nil
}

func (cfg *kdbConfig) Address() string {
	return fmt.Sprintf("%s:%d", cfg.host, cfg.port)
}

func (cfg *kdbConfig) Username() string {
	return cfg.username
}

func (cfg *kdbConfig) Password() string {
	return cfg.password
}

func (cfg *kdbConfig) DB() int {
	return cfg.db
}

func (cfg *kdbConfig) ReadTimeoutSec() time.Duration {
	return cfg.readTimeout
}

func (cfg *kdbConfig) WriteTimeoutSec() time.Duration {
	return cfg.writeTimeout
}
