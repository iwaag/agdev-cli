package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	APIBaseURL     string
	AuthToken      string
	RequestTimeout time.Duration
	OutputJSON     bool
	LogLevel       string
}

func Load(flagJSON bool) (Config, error) {
	v := viper.New()
	v.SetEnvPrefix("AGDEV")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()

	v.SetDefault("api.base_url", "")
	v.SetDefault("auth.token", "")
	v.SetDefault("request.timeout", "30s")
	v.SetDefault("output.json", false)
	v.SetDefault("log.level", "info")

	timeout, err := time.ParseDuration(v.GetString("request.timeout"))
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		APIBaseURL:     v.GetString("api.base_url"),
		AuthToken:      v.GetString("auth.token"),
		RequestTimeout: timeout,
		OutputJSON:     flagJSON || v.GetBool("output.json"),
		LogLevel:       v.GetString("log.level"),
	}

	return cfg, nil
}
