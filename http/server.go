package http

import (
	"context"
	"net"
	"net/http"
	"sync"
	"time"

	syncx "github.com/jonasohland/ext/sync"
	slx "github.com/jonasohland/slog-ext/pkg/slog-ext"
)

type Server interface {
	Wait(timeout time.Duration) error
}

type contextServer struct {
	listener net.Listener
	server   http.Server
	wg       sync.WaitGroup
}

func (s *contextServer) Wait(timeout time.Duration) error {
	return syncx.WaitTimeout(&s.wg, timeout)
}

func (s *contextServer) run(ctx context.Context) {
	defer s.wg.Done()
	if err := s.server.Serve(s.listener); err != nil {
		if err != http.ErrServerClosed {
			slx.FromContext(ctx).Error("serve", "error", err)
		}
	}
}

func (s *contextServer) waitShutdown(ctx context.Context) {
	defer s.wg.Done()
	<-ctx.Done()
	if err := s.server.Shutdown(context.Background()); err != nil {
		slx.FromContext(ctx).Error("shutdown", "error", err)
	}
}

func NewContextServer(ctx context.Context, handler http.Handler, network string, address string) (Server, error) {
	srv := &contextServer{}

	srv.server.Addr = address
	srv.server.BaseContext = func(l net.Listener) context.Context { return ctx }
	srv.server.Handler = handler

	listener, err := net.Listen(network, address)
	if err != nil {
		return nil, err
	}

	srv.listener = listener

	slx.FromContext(ctx).Info("listening", "address", listener.Addr().String())

	srv.wg.Add(2)
	go srv.run(ctx)
	go srv.waitShutdown(ctx)
	return srv, nil
}
