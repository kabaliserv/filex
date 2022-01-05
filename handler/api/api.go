package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/kabaliserv/filex/cmd/kbs-filex/config"
	"github.com/kabaliserv/filex/core"
	"github.com/kabaliserv/filex/handler/api/auth"
	"github.com/kabaliserv/filex/handler/api/files"
	"github.com/kabaliserv/filex/storage"
	"net/http"
)

type Server struct {
	config  config.Config
	storage storage.Storage
	store   core.Store
}

func New(
	config config.Config,
	storage storage.Storage,
	store core.Store,
) Server {
	return Server{
		config:  config,
		storage: storage,
		store:   store,
	}
}

func (s Server) Handler() http.Handler {
	r := chi.NewRouter()

	r.Mount("/files", files.Handle(s.storage))
	r.Mount("/auth", auth.NewHandler(s.store))

	return r
}
