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

type LogoConfig struct {
	Src    string `mapstructure:"src"`
	Width  string `mapstructure:"width"`
	Height string `mapstructure:"height"`
	Link   string `mapstructure:"link"`
	Alt    string `mapstructure:"alt"`
}

type PageConfig struct {
	Logo     *LogoConfig `mapstructure:"logo"`
	Title    string      `mapstructure:"title"`
	Subtitle string      `mapstructure:"subtitle"`
}

type Config struct {
	Port        int           `mapstructure:"port"`
	DatabaseURL string        `mapstructure:"databaseUrl"`
	Github      *GithubConfig `mapstructure:"github"`
	Local       *LocalConfig  `mapstructure:"local"`
	Page        *PageConfig   `mapstructure:"page"`
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
