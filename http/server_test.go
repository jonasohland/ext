package http_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	htx "github.com/jonasohland/ext/http"
)

func TestServerShutdown(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	s, err := htx.NewContextServer(ctx, http.DefaultServeMux, "tcp", "127.0.0.1:45893")
	if err != nil {
		panic(err)
	}

	if err := s.Wait(time.Second * 7); err != nil {
		panic(err)
	}

}
