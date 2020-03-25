package main

import (
	"fmt"
	"github.com/aibotsoft/gproxy"
	"github.com/aibotsoft/gproxy/internal/config"
	"github.com/aibotsoft/gproxy/internal/pg"
	"github.com/aibotsoft/micro/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log := logger.New()
	cfg := config.New()
	log.Info(cfg)

	db, err := pg.New()
	if err != nil {
		log.Fatal(err)
	}
	errc := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	// Run gPRC proxy server
	s := gproxy.NewServer(db)
	go func() {
		errc <- s.Serve()
	}()

	log.Info("exit: ", <-errc)
	s.GracefulStop()
}
