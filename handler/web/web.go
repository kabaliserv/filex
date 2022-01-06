package web

import (
	"net/http"

	"github.com/kabaliserv/filex/web"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/unrolled/secure"
)

type Server struct {
	Options secure.Options
}

func New(option secure.Options) Server {
	return Server{
		Options: option,
	}
}

func (s Server) Handler() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.NoCache)
	r.Use(middleware.StripSlashes)

	sec := secure.New(s.Options)
	r.Use(sec.Handler)

	h := http.FileServer(http.FS(web.New()))

	r.Handle("/favicon.ico", h)
	r.Handle("/assets/*", h)
	r.NotFound(HandleIndex())

	return r
}
