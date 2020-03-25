package gproxy

import (
	"context"
	"github.com/aibotsoft/micro/logger"
	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"time"
)

const (
	port = ":50051"
)

type Server struct {
	log   *zap.SugaredLogger
	store *Store
	gs    *grpc.Server
	UnimplementedProxyServer
}

func NewServer(db *pgx.Conn) *Server {
	log := logger.New()
	return &Server{
		log:   log,
		store: New(log, db),
		gs:    grpc.NewServer(),
	}
}

func (s *Server) GracefulStop() {
	s.log.Debug("begin proxy server gracefulStop")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s.gs.GracefulStop()
	err := s.store.db.Close(ctx)
	if err != nil {
		s.log.Error(err)
	}
	s.log.Debug("end proxy server gracefulStop")
}

func (s *Server) Serve() error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}
	RegisterProxyServer(s.gs, s)
	s.log.Info("gRPC Proxy Server listens on port ", port)
	return s.gs.Serve(lis)
}
