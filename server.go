package gproxy

import (
	"github.com/aibotsoft/micro/logger"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
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

func NewServer(db *pgxpool.Pool) *Server {
	log := logger.New()
	return &Server{
		log:   log,
		store: New(log, db),
		gs:    grpc.NewServer(),
	}
}

func (s *Server) GracefulStop() {
	s.log.Debug("begin proxy server gracefulStop")
	s.gs.GracefulStop()
	s.store.db.Close()
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
