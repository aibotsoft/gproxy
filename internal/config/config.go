package config

import (
	"github.com/kelseyhightower/envconfig"
	"os"
)

const envFileName = ".env"
const devEnv = "dev"

func init() {
	if os.Getenv("SERVICE_ENV") == devEnv {
		MustLoadEnv()
	}
}

// New returns the settings from the environment.
func New() *Config {
	cfg := &Config{}
	err := envconfig.Process("", cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}

func IsDev() bool {
	return os.Getenv("SERVICE_ENV") == devEnv
}
