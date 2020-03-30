package gproxy_test

import (
	"github.com/aibotsoft/gproxy"
	"github.com/aibotsoft/micro/config"
	"github.com/aibotsoft/micro/logger"
	"github.com/aibotsoft/micro/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/subosito/gotenv"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

//func Test_server_CreateProxy(t *testing.T) {
//	s := newServer(t)
//	defer s.GracefulStop()
//	req := &gproxy.CreateProxyRequest{ProxyItem: &gproxy.ProxyItem{
//		ProxyIp:   "0.0.0.0",
//		ProxyPort: 80,
//		ProxyCountry: &gproxy.ProxyCountry{
//			CountryName: "Unknown",
//			CountryCode: "NA",
//		},
//	},
//	}
//	got, err := s.CreateProxy(nil, req)
//	assert.NoError(t, err)
//	assert.Equal(t, req.GetProxyItem(), got.GetProxyItem())
//
//}
//
//func TestServer_GetNextProxy(t *testing.T) {
//	s := newServer(t)
//	defer s.GracefulStop()
//	req := &gproxy.GetNextProxyRequest{}
//	got, err := s.GetNextProxy(nil, req)
//	if assert.NoError(t, err) {
//		assert.NotEmpty(t, got)
//	}
//}

func TestServer_GetNextProxy(t *testing.T) {

	gotenv.Must(gotenv.Load)
	log := logger.New()
	cfg := config.New()
	log.Info(cfg)
	db := postgres.MustConnect(cfg)

	s := gproxy.NewServer(db)
	got, err := s.GetNextProxy(nil, nil)
	//code:= status.Convert(err).Code()
	assert.Equal(t, codes.NotFound, status.Convert(err).Code())
	assert.Empty(t, got)

}