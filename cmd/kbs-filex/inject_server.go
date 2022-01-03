package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/google/wire"
	"github.com/kabaliserv/filex/cmd/kbs-filex/config"
	"github.com/kabaliserv/filex/handler/api"
	"github.com/kabaliserv/filex/handler/web"
	"github.com/kabaliserv/filex/server"
)

var serverSet = wire.NewSet(
	api.New,
	web.New,
	provideRouter,
	provideServer,
)

func provideRouter(api api.Server, web web.Server) *chi.Mux {
	r := chi.NewRouter()
	r.Mount("/api", api.Handler())
	r.Mount("/", web.Handler())
	return r
}

func provideServer(handler *chi.Mux, config config.Config) *server.Server {
	return &server.Server{
		Addr:    config.Server.Port,
		Host:    config.Server.Host,
		Handler: handler,
	}
}
