package config

import (
	"github.com/spf13/viper"
	"time"
)

type (
	Config struct {
		HTTP  HTTPConfig  `mapstructure:"http"`
		Other OtherConfig `mapstructure:"others"`
	}
	HTTPConfig struct {
		Port               string        `mapstructure:"port"`
		ReadTimeout        time.Duration `mapstructure:"readTimeout"`
		WriteTimeout       time.Duration `mapstructure:"writeTimeout"`
		MaxHeaderMegabytes int           `mapstructure:"maxHeaderMegaBytes"`
	}
	OtherConfig struct {
		RequestCounterTTL int    `mapstructure:"requestCounterTTL"`
		BackUpFileName    string `mapstructure:"backUpFileName"`
	}
)

func Init(configsDir, environment string) (*Config, error) {
	viper.AddConfigPath(configsDir)
	viper.SetConfigName(environment)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
