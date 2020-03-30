package gproxy_test

import (
	"github.com/aibotsoft/gproxy"
	"github.com/aibotsoft/micro/config"
	"github.com/aibotsoft/micro/logger"
	"github.com/aibotsoft/micro/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/subosito/gotenv"
	"testing"
	"time"
)

func TestNewServer(t *testing.T) {
	gotenv.Must(gotenv.Load)
	var err error
	cfg := config.New()
	log := logger.New()
	db := postgres.MustConnect(cfg)
	defer db.Close()
	server := gproxy.NewServer(cfg, log, db)
	go func() {
		err = server.Serve()
	}()
	time.Sleep(time.Millisecond * 100)
	assert.NoError(t, err)
	server.GracefulStop()
}
