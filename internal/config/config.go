package config

import (
	"github.com/spf13/viper"
)

type GitHubConfig struct {
	Owner             string `mapstructure:"owner"`
	Repo              string `mapstructure:"repo"`
	Path              string `mapstructure:"path"`
	AppPrivateKey     string `mapstructure:"appPrivateKey"`
	AppInstallationId int64  `mapstructure:"appInstallationId"`
}

type LocalConfig struct {
	FilesPath string `mapstructure:"filesPath"`
}

type Config struct {
	GitHub GitHubConfig `mapstructure:"github"`
	Local  LocalConfig  `mapstructure:"local"`
}

func Load() (Config, error) {
	viper.SetConfigFile("config.yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		return Config{}, err
	}

	cfg := Config{}
	err = viper.Unmarshal(&cfg)
	return cfg, err
}
