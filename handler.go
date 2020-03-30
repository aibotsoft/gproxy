package gproxy

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateProxyStat(ctx context.Context, req *CreateProxyStatRequest) (*CreateProxyStatResponse, error) {
	proxyStat := req.GetProxyStat()
	err := s.store.CreateProxyStat(ctx, proxyStat)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "CreateProxyStat: %v", err)
	}
	return &CreateProxyStatResponse{ProxyStat: proxyStat}, nil
}

func (s *Server) CreateProxy(ctx context.Context, req *CreateProxyRequest) (*CreateProxyResponse, error) {
	proxyItem := req.GetProxyItem()
	err := s.store.GetOrCreateProxyItem(ctx, proxyItem)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "GetOrCreateProxyItem: %v", err)
	}
	return &CreateProxyResponse{ProxyItem: proxyItem}, nil
}

// GetNextProxy возвращает прокси которое нужно проверить.
// Возвращаются те которые еще не проверялись, либо отсортированные по времени проверки.
func (s *Server) GetNextProxy(ctx context.Context, req *GetNextProxyRequest) (*GetNextProxyResponse, error) {
	proxyItem, err := s.store.GetNextProxyItem(ctx)
	switch err {
	case nil:
		return &GetNextProxyResponse{ProxyItem: proxyItem}, nil
	case ErrNoRows:
		return nil, status.Error(codes.NotFound, "no new proxy for check to return")
	default:
		return nil, status.Errorf(codes.Internal, "GetNextProxy error %v", err)
	}
}
