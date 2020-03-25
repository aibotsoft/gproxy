package gproxy_test

import (
	"github.com/aibotsoft/gproxy"
	"github.com/aibotsoft/gproxy/internal/pg"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func newServer(t *testing.T) *gproxy.Server {
	t.Helper()
	db, err := pg.New()
	assert.NoError(t, err)
	server := gproxy.NewServer(db)
	return server
}

func TestNewServer(t *testing.T) {
	var err error
	db, err := pg.New()
	assert.NoError(t, err)
	server := gproxy.NewServer(db)
	go func() {
		err = server.Serve()
	}()
	time.Sleep(time.Millisecond * 100)
	assert.NoError(t, err)
	server.GracefulStop()
}
