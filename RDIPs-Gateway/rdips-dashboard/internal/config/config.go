package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	APIBaseURL  string
	APIKey      string
	Port        string
	HTTPTimeout time.Duration
}

func Load() (Config, error) {
	v := viper.New()
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(".")
	v.AutomaticEnv()

	v.SetDefault("API_BASE_URL", "http://sunflower_api:8080")
	v.SetDefault("API_KEY", "a706913c38b555a889218175")
	v.SetDefault("PORT", "3000")
	v.SetDefault("HTTP_TIMEOUT_SECONDS", 20)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return Config{}, err
		}
	}

	return Config{
		APIBaseURL:  strings.TrimRight(v.GetString("API_BASE_URL"), "/"),
		APIKey:      v.GetString("API_KEY"),
		Port:        v.GetString("PORT"),
		HTTPTimeout: time.Duration(v.GetInt("HTTP_TIMEOUT_SECONDS")) * time.Second,
	}, nil
}
