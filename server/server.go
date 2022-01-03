package server

import (
	"context"
	"golang.org/x/sync/errgroup"
	"net/http"
	"time"
)

type Server struct {
	Addr    string
	Host    string
	Handler http.Handler
}

const timeoutGracefulShutdown = 5 * time.Second

func (s Server) ListenAndServe(ctx context.Context) error {
	err := s.listenAndServe(ctx)
	if err == http.ErrServerClosed {
		err = nil
	}
	return err
}

func (s Server) listenAndServe(ctx context.Context) error {
	var g errgroup.Group
	s1 := &http.Server{
		Addr:    s.Addr,
		Handler: s.Handler,
	}
	g.Go(func() error {
		<-ctx.Done()

		ctxShutdown, cancelFunc := context.WithTimeout(context.Background(), timeoutGracefulShutdown)
		defer cancelFunc()

		return s1.Shutdown(ctxShutdown)
	})
	g.Go(s1.ListenAndServe)
	return g.Wait()
}
