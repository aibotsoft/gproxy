package gproxy_test

import (
	"github.com/aibotsoft/gproxy"
	"github.com/aibotsoft/micro/config"
	"github.com/aibotsoft/micro/logger"
	"github.com/aibotsoft/micro/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/subosito/gotenv"
	"testing"
)

func TestStore_GetNextProxyItem(t *testing.T) {
	gotenv.Must(gotenv.Load)
	log := logger.New()
	cfg := config.New()
	log.Info(cfg)

	db := postgres.MustConnect(cfg)
	assert.NotEmpty(t, db)
	store :=gproxy.New(log, db)
	err := store.GetNextProxyItem(&gproxy.ProxyItem{})
	t.Log(err)
	//s := &Store{
	//	log:   ni,
	//	db:    tt.fields.db,
	//	cache: tt.fields.cache,
	//}
	//if err := s.GetNextProxyItem(tt.args.p); (err != nil) != tt.wantErr {
	//	t.Errorf("GetNextProxyItem() error = %v, wantErr %v", err, tt.wantErr)
	//}
}