package main

import (
	"context"
	"github.com/aibotsoft/gproxy"
	"github.com/aibotsoft/micro/logger"
	"google.golang.org/grpc"
	"time"
)

const (
	port = ":50051"
)

func main() {
	log := logger.New()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, port, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := gproxy.NewProxyClient(conn)
	log.Info(client)
	req := &gproxy.CreateProxyRequest{
		ProxyItem: &gproxy.ProxyItem{},
	}

	res, err := client.CreateProxy(context.Background(), req)
	if err != nil {
		log.Error(err)
	}
	log.Info(res)
}
