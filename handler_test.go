package gproxy_test

import (
	"github.com/aibotsoft/gproxy"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_server_CreateProxy(t *testing.T) {
	s := newServer(t)
	defer s.GracefulStop()
	req := &gproxy.CreateProxyRequest{ProxyItem: &gproxy.ProxyItem{
		ProxyIp:   "0.0.0.0",
		ProxyPort: 80,
		ProxyCountry: &gproxy.ProxyCountry{
			CountryName: "Unknown",
			CountryCode: "NA",
		},
	},
	}
	got, err := s.CreateProxy(nil, req)
	assert.NoError(t, err)
	assert.Equal(t, req.GetProxyItem(), got.GetProxyItem())

}

func TestServer_GetNextProxy(t *testing.T) {
	s := newServer(t)
	defer s.GracefulStop()
	req := &gproxy.GetNextProxyRequest{}
	got, err := s.GetNextProxy(nil, req)
	if assert.NoError(t, err) {
		assert.NotEmpty(t, got)
	}

}
