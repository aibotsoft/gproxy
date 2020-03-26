package main

import (
	"fmt"
	"github.com/aibotsoft/gproxy"
	"github.com/aibotsoft/micro/config"
	"github.com/aibotsoft/micro/logger"
	"github.com/aibotsoft/micro/postgres"
	"github.com/subosito/gotenv"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	gotenv.Must(gotenv.Load)
	log := logger.New()
	cfg := config.New()
	log.Info(cfg)

	db, err := postgres.New(cfg)
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
