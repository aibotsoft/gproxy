package gproxy

import (
	"github.com/aibotsoft/micro/config"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"strconv"
)

type Server struct {
	cfg   *config.Config
	log   *zap.SugaredLogger
	store *Store
	gs    *grpc.Server
	UnimplementedProxyServer
}

func NewServer(cfg *config.Config, log *zap.SugaredLogger, db *pgxpool.Pool) *Server {
	return &Server{
		cfg:   cfg,
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
	lis, err := net.Listen("tcp", strconv.Itoa(s.cfg.ProxyService.GRPCPort))
	if err != nil {
		return err
	}
	RegisterProxyServer(s.gs, s)
	s.log.Info("gRPC Proxy Server listens on port ", strconv.Itoa(s.cfg.ProxyService.GRPCPort))
	return s.gs.Serve(lis)
}
