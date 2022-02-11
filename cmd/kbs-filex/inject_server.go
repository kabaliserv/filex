package main

import (
	"github.com/kabaliserv/filex/cmd/kbs-filex/config"
	"github.com/kabaliserv/filex/core"
	"github.com/kabaliserv/filex/handler/api"
	"github.com/kabaliserv/filex/handler/health"
	"github.com/kabaliserv/filex/handler/web"
	"github.com/kabaliserv/filex/server"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/wire"
	"github.com/unrolled/secure"
)

type (
	healthzHandler http.Handler
)

var serverSet = wire.NewSet(
	api.New,
	web.New,
	provideRouter,
	provideServer,
	provideHealthz,
	provideServerOptions,
	provideUploadOptions,
)

func provideRouter(api api.Server, web web.Server, healthz healthzHandler) *chi.Mux {
	r := chi.NewRouter()
	r.Mount("/api", api.Handler())
	r.Mount("/healthz", healthz)
	r.Mount("/", web.Handler())
	return r
}

func provideHealthz() healthzHandler {
	v := health.New()
	return healthzHandler(v)
}

func provideUploadOptions(config config.Config) core.UploadOption {
	return core.UploadOption{
		GuestAllow:         config.Guest.AllowUpload,
		GuestMaxUploadSize: config.Guest.MaxUploadSize,
		GuestMaxDuration:   config.Guest.MaxFileDuration,
	}
}

func provideServer(handler *chi.Mux, config config.Config) *server.Server {
	return &server.Server{
		Addr:    config.Server.Port,
		Host:    config.Server.Host,
		Cert:    config.Server.Cert,
		Key:     config.Server.Key,
		Handler: handler,
	}
}

func provideServerOptions(config config.Config) secure.Options {
	return secure.Options{
		AllowedHosts:          config.HTTP.AllowedHosts,
		HostsProxyHeaders:     config.HTTP.HostsProxyHeaders,
		SSLRedirect:           config.HTTP.SSLRedirect,
		SSLTemporaryRedirect:  config.HTTP.SSLTemporaryRedirect,
		SSLHost:               config.HTTP.SSLHost,
		SSLProxyHeaders:       config.HTTP.SSLProxyHeaders,
		STSSeconds:            config.HTTP.STSSeconds,
		STSIncludeSubdomains:  config.HTTP.STSIncludeSubdomains,
		STSPreload:            config.HTTP.STSPreload,
		ForceSTSHeader:        config.HTTP.ForceSTSHeader,
		FrameDeny:             config.HTTP.FrameDeny,
		ContentTypeNosniff:    config.HTTP.ContentTypeNosniff,
		BrowserXssFilter:      config.HTTP.BrowserXSSFilter,
		ContentSecurityPolicy: config.HTTP.ContentSecurityPolicy,
		ReferrerPolicy:        config.HTTP.ReferrerPolicy,
	}
}
