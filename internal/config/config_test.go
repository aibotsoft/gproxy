package config_test

import (
	"github.com/aibotsoft/gproxy/internal/config"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestLoadEnv(t *testing.T) {
	if !config.IsDev() {
		return
	}
	err := config.LoadEnv()
	assert.NoError(t, err)
	assert.Equal(t, "true", os.Getenv("TEST_LOAD_ENV"))
}

func TestNew(t *testing.T) {
	cfg := config.New()
	assert.Equal(t, true, cfg.Service.TestLoadEnv)
}
