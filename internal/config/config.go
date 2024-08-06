package config

import (
	"github.com/spf13/viper"

	"github.com/VadimGossip/tj-drs-storage/internal/domain"
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
