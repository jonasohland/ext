package grpc

import (
	"context"
	"net"
	"sync"
	"time"

	httpx "github.com/jonasohland/ext/http"
	syncx "github.com/jonasohland/ext/sync"
	slogext "github.com/jonasohland/slog-ext/pkg/slog-ext"
	"google.golang.org/grpc"
)

type contextServer struct {
	wg sync.WaitGroup

	srv      *grpc.Server
	listener net.Listener
}

// Wait implements http.Server.
func (c *contextServer) Wait(timeout time.Duration) error {
	return syncx.WaitTimeout(&c.wg, timeout)
}

func (s *contextServer) run(ctx context.Context) {
	s.srv.Serve(s.listener)
	defer s.wg.Done()
	slogext.FromContext(ctx).Debug("done serving")
}

func (s *contextServer) waitShutdown(ctx context.Context) {
	defer s.wg.Done()
	<-ctx.Done()
	slogext.FromContext(ctx).Debug("shutting down")
	s.srv.GracefulStop()
	slogext.FromContext(ctx).Debug("shutdown ok")
}

func NewContextServer(ctx context.Context, srv *grpc.Server, network string, address string) (httpx.Server, error) {
	listener, err := net.Listen(network, address)
	if err != nil {
		return nil, err
	}

	c := &contextServer{
		srv:      srv,
		listener: listener,
	}

	slogext.FromContext(ctx).Info("listening", "network", network, "address", listener.Addr())

	c.wg.Add(2)
	go c.run(ctx)
	go c.waitShutdown(ctx)
	return c, nil
}
