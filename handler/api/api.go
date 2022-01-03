package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/kabaliserv/filex/cmd/kbs-filex/config"
	"github.com/kabaliserv/filex/storage"
	tusd "github.com/tus/tusd/pkg/handler"
	"net/http"
)

type Server struct {
	config  config.Config
	storage storage.Storage
}

func New(config config.Config, storage storage.Storage) Server {
	return Server{
		config:  config,
		storage: storage,
	}
}

func (s Server) Handler() http.Handler {
	r := chi.NewRouter()

	r.Route("/files", func(r chi.Router) {
		handler, err := tusd.NewUnroutedHandler(tusd.Config{
			BasePath:      "/api/files",
			StoreComposer: s.storage.GetStoreComposer(),
		})

		if err != nil {
			panic(err)
		}

		r.Post("/", http.HandlerFunc(handler.PostFile))
		r.Head("/{id:[-+a-z0-9]+}", http.HandlerFunc(handler.HeadFile))
		r.Patch("/{id:[-+a-z0-9]+}", http.HandlerFunc(handler.PatchFile))
		r.Get("/{id:[-+a-z0-9]+}", http.HandlerFunc(handler.GetFile))
		r.Delete("/{id:[-+a-z0-9]+}", http.HandlerFunc(handler.DelFile))

	})

	return r
}
