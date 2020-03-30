package gproxy

import (
	"context"
	"github.com/aibotsoft/micro/config"
	"github.com/aibotsoft/micro/logger"
	"github.com/aibotsoft/micro/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/subosito/gotenv"
	"testing"
)

func TestStore_GetNextProxyItem_Multi(t *testing.T) {
	gotenv.Must(gotenv.Load)
	log := logger.New()
	cfg := config.New()
	db := postgres.MustConnect(cfg)
	assert.NotEmpty(t, db)
	store :=New(log, db)
	for i := 0; i < 100; i++ {
		got, err := store.GetNextProxyItem(context.Background())
		assert.NoError(t, err)
		assert.NotEmpty(t, got)
		t.Log(got)
	}
}
