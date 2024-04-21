package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	GH_APP_PRIVATE_KEY     string `mapstructure:"GH_APP_PRIVATE_KEY"`
	GH_APP_INSTALLATION_ID int64  `mapstructure:"GH_APP_INSTALLATION_ID"`
}

func Load() (Config, error) {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		return Config{}, err
	}

	cfg := Config{}
	err = viper.Unmarshal(&cfg)
	return cfg, err
}
