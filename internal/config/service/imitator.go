package service

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
)

const (
	imitatorRequestTypeEnvName = "IMITATOR_REQUEST_TYPE"
	imitatorRPS                = "IMITATOR_RPS"
	imitatorPPS                = "IMITATOR_PPS"
	imitatorTOTAL              = "IMITATOR_TOTAL"
)

const (
	allSupplierRequestType = "ALL_SUP"
	singleRequestType      = "SINGLE"
)

type imitatorConfig struct {
	requestType string
	rps         int
	pps         int
	total       int
}

func (cfg *imitatorConfig) setFromEnv() error {
	var err error
	cfg.requestType = os.Getenv(imitatorRequestTypeEnvName)
	if len(cfg.requestType) == 0 {
		return fmt.Errorf("imitatorConfig requestType not found")
	}

	if cfg.requestType != allSupplierRequestType && cfg.requestType != singleRequestType {
		return fmt.Errorf("imitatorConfig not supported requestType, %s", cfg.requestType)
	}

	rpsStr := os.Getenv(imitatorRPS)
	if len(rpsStr) == 0 {
		return fmt.Errorf("imitatorConfig rps not found")
	}

	cfg.rps, err = strconv.Atoi(rpsStr)
	if err != nil {
		return errors.Wrap(err, "failed to parse imitatorConfig rps")
	}

	ppsStr := os.Getenv(imitatorPPS)
	if len(ppsStr) == 0 {
		return fmt.Errorf("imitatorConfig pps not found")
	}

	cfg.pps, err = strconv.Atoi(ppsStr)
	if err != nil {
		return errors.Wrap(err, "failed to parse imitatorConfig pps")
	}

	totalStr := os.Getenv(imitatorTOTAL)
	if len(totalStr) == 0 {
		return fmt.Errorf("imitatorConfig total not found")
	}

	cfg.total, err = strconv.Atoi(totalStr)
	if err != nil {
		return errors.Wrap(err, "failed to parse imitatorConfig total")
	}

	return nil
}

func NewImitatorConfig() (*imitatorConfig, error) {
	cfg := &imitatorConfig{}
	if err := cfg.setFromEnv(); err != nil {
		return nil, fmt.Errorf("imitatorConfig set from env err: %s", err)
	}
	logrus.Infof("imitatorConfig: [%+v]", *cfg)

	return cfg, nil
}

func (cfg *imitatorConfig) RequestType() string {
	return cfg.requestType
}

func (cfg *imitatorConfig) AllSupplierRequestType() string {
	return allSupplierRequestType
}

func (cfg *imitatorConfig) SingleRequestType() string {
	return singleRequestType
}

func (cfg *imitatorConfig) RequestPerSecond() int {
	return cfg.rps
}

func (cfg *imitatorConfig) PackPerSecond() int {
	return cfg.pps
}

func (cfg *imitatorConfig) TotalRequests() int {
	return cfg.total
}
