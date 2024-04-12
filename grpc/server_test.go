package grpc_test

import (
	"context"
	"testing"
	"time"

	grpcx "github.com/jonasohland/ext/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
)

func TestServerShutdown(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	srv := grpc.NewServer()
	healthServer := health.NewServer()
	healthgrpc.RegisterHealthServer(srv, healthServer)

	csrv, err := grpcx.NewContextServer(ctx, srv, "tcp", "127.0.0.1:8029")
	if err != nil {
		panic(err)
	}

	<-ctx.Done()
	csrv.Wait(time.Second * 2)
}
