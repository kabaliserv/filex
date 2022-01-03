package web

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kabaliserv/filex/web"
	"net/http"
)

type Server struct {
}

func New() Server {
	return Server{}
}

func (s Server) Handler() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.NoCache)
	r.Use(middleware.StripSlashes)

	h := http.FileServer(http.FS(web.New()))

	r.Handle("/favicon.ico", h)
	r.Handle("/assets/*", h)
	r.NotFound(HandleIndex())

	return r
}
