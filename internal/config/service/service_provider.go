package service

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

const (
	testDbEnvName = "TEST_DB"
)
const (
	tarantoolTestDB = "TARANTOOL"
	kdbTestDb       = "KDB"
	cacheTestDb     = "CACHE"
)

type serviceProviderConfig struct {
	testDB string
}

func (cfg *serviceProviderConfig) setFromEnv() error {
	cfg.testDB = os.Getenv(testDbEnvName)
	if len(cfg.testDB) == 0 {
		return fmt.Errorf("serviceProviderConfig testDB not found")
	}
	if cfg.testDB != tarantoolTestDB && cfg.testDB != kdbTestDb && cfg.testDB != cacheTestDb {
		return fmt.Errorf("serviceProviderConfig not supported testDB, %s", cfg.testDB)
	}

	return nil
}

func NewServiceProviderConfig() (*serviceProviderConfig, error) {
	cfg := &serviceProviderConfig{}
	if err := cfg.setFromEnv(); err != nil {
		return nil, fmt.Errorf("serviceProviderConfig set from env err: %s", err)
	}
	logrus.Infof("serviceProviderConfig: [%+v]", *cfg)

	return cfg, nil
}

func (cfg *serviceProviderConfig) TestDB() string {
	return cfg.testDB
}

func (cfg *serviceProviderConfig) TarantoolTestDB() string {
	return tarantoolTestDB
}

func (cfg *serviceProviderConfig) KdbTestDB() string {
	return kdbTestDb
}

func (cfg *serviceProviderConfig) CacheTestDB() string {
	return cacheTestDb
}
