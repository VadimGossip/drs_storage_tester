package config

import (
	"github.com/VadimGossip/drs_storage_tester/pkg/util"
	"github.com/spf13/viper"

	"github.com/VadimGossip/drs_storage_tester/internal/domain"
)

func parseConfigFile(configDir string) error {
	viper.AddConfigPath(configDir)
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func unmarshal(cfg *domain.Config) error {
	if err := viper.UnmarshalKey("keydb", &cfg.TargetDb); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("oracle", &cfg.DataSourceDb); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("task", &cfg.Task); err != nil {
		return err
	}
	cfg.Task.Summary.DbDuration = &domain.DurationSummary{EMA: util.NewEMA(0.01), Histogram: make(map[float64]int)}
	cfg.Task.Summary.TotalDuration = &domain.DurationSummary{EMA: util.NewEMA(0.01), Histogram: make(map[float64]int)}
	return nil
}

func Init(configDir string) (*domain.Config, error) {
	if err := parseConfigFile(configDir); err != nil {
		return nil, err
	}
	cfg := &domain.Config{}
	if err := unmarshal(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
