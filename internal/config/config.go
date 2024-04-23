package config

import (
	"github.com/spf13/viper"
)

type GithubConfig struct {
	Owner string      `mapstructure:"owner"`
	Repo  string      `mapstructure:"repo"`
	Path  string      `mapstructure:"path"`
	Auth  *GithubAuth `mapstructure:"auth"`
}

type GithubAuth struct {
	AppPrivateKey     string `mapstructure:"appPrivateKey"`
	AppInstallationId int64  `mapstructure:"appInstallationId"`
	AccessToken       string `mapstructure:"accessToken"`
}

type LocalConfig struct {
	FilesPath string `mapstructure:"filesPath"`
}

type Config struct {
	Github *GithubConfig `mapstructure:"github"`
	Local  *LocalConfig  `mapstructure:"local"`
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
