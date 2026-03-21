package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	OutputJSON bool
	LogLevel   string
}

func Load(flagJSON bool) (Config, error) {
	v := viper.New()
	v.SetEnvPrefix("AGDEV")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()

	v.SetDefault("output.json", false)
	v.SetDefault("log.level", "info")

	cfg := Config{
		OutputJSON: flagJSON || v.GetBool("output.json"),
		LogLevel:   v.GetString("log.level"),
	}

	return cfg, nil
}
