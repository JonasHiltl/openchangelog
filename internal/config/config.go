package config

import (
	"log/slog"

	"github.com/spf13/viper"
)

type CacheTyp string
type AnalyticsProvider string
type SearchType string

const (
	Memory     CacheTyp          = "memory"
	S3         CacheTyp          = "s3"
	Disk       CacheTyp          = "disk"
	Tinybird   AnalyticsProvider = "tinybird"
	SearchDisk SearchType        = "disk"
	SearchMem  SearchType        = "memory"
)

type GithubConfig struct {
	Owner string      `mapstructure:"owner"`
	Repo  string      `mapstructure:"repo"`
	Path  string      `mapstructure:"path"`
	Auth  *GithubAuth `mapstructure:"auth"`
}

type CacheConfig struct {
	Type CacheTyp         `mapstructure:"type"`
	Disk *DiskCacheConfig `mapstructure:"disk"`
	S3   *S3CacheConfig   `mapstructure:"s3"`
}

type DiskCacheConfig struct {
	Location string `mapstructure:"location"`
	MaxSize  uint64 `mapstructure:"maxSize"`
}

type S3CacheConfig struct {
	Bucket string `mapstructure:"bucket"`
}

type GithubAuth struct {
	AppID             int64  `mapstructure:"appId"`
	AppSecret         string `mapstructure:"appSecret"`
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
	Logo          *LogoConfig `mapstructure:"logo"`
	Title         string      `mapstructure:"title"`
	Subtitle      string      `mapstructure:"subtitle"`
	ColorScheme   string      `mapstructure:"colorScheme"`
	HidePoweredBy bool        `mapstructure:"hidePoweredBy"`
	Auth          *AuthConfig `mapstructure:"auth"`
}

type AuthConfig struct {
	Enabled      bool   `mapstructure:"enabled"`
	PasswordHash string `mapstructure:"passwordHash"`
}

type AnalyticsConfig struct {
	Provider AnalyticsProvider `mapstructure:"provider"`
	Tinybird *TinybirdConfig   `mapstructure:"tinybird"`
}

type TinybirdConfig struct {
	AccessToken string `mapstructure:"accessToken"`
}

type AdminConfig struct {
	PasswordHash string `mapstructure:"passwordHash"`
}

type LogConfig struct {
	Level LogLevel `mapstructure:"level"`
	Style string   `mapstructure:"style"`
}

type LogLevel string

func (l LogLevel) ToSlog() slog.Level {
	var sl slog.Level
	sl.UnmarshalText([]byte(l))
	return sl
}

type SearchConfig struct {
	Type SearchType `mapstructure:"type"`
	Disk struct {
		Path string `mapstructure:"path"`
	} `mapstructure:"disk"`
}

type Config struct {
	Addr      string           `mapstructure:"addr"`
	SqliteURL string           `mapstructure:"sqliteUrl"`
	Github    *GithubConfig    `mapstructure:"github"`
	Local     *LocalConfig     `mapstructure:"local"`
	Page      *PageConfig      `mapstructure:"page"`
	Cache     *CacheConfig     `mapstructure:"cache"`
	Analytics *AnalyticsConfig `mapstructure:"analytics"`
	Admin     *AdminConfig     `mapstructure:"admin"`
	Log       *LogConfig       `mapstructure:"log"`
	Search    *SearchConfig    `mapstructure:"search"`
}

func (c Config) HasGithubAuth() bool {
	return c.Github != nil && c.Github.Auth != nil
}

// Returns true if Openchangelog was started in db mode (using sqlite to store changelog configs, multi tenancy)
func (c Config) IsDBMode() bool {
	return c.SqliteURL != ""
}

// Returns true if Openchangelog was started in config mode (using config file to store changelog configs)
func (c Config) IsConfigMode() bool {
	return !c.IsDBMode()
}

// Loads the config file from configPath if specified.
//
// Otherwise searches the standard list of search paths. Returns an error if
// no configuration files could be found.
func Load(configFile string) (Config, error) {
	viper.SetConfigType("yaml")
	viper.SetConfigName("openchangelog") // extension needs to be in the name, otherwise openchangelog binary might be read

	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.AddConfigPath("/etc/")
		viper.AddConfigPath(".")
	}
	err := viper.ReadInConfig()
	if err != nil {
		return Config{}, err
	}

	viper.SetDefault("addr", "localhost:6001")

	cfg := Config{}
	err = viper.Unmarshal(&cfg)
	return cfg, err
}
