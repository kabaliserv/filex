package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/kabaliserv/filex/core"
	acl "github.com/kabaliserv/filex/handler/api/acl"
	"github.com/kabaliserv/filex/handler/api/auth"
	tusd "github.com/tus/tusd/pkg/handler"
	"net/http"
)

var corsOpts = cors.Options{
	AllowedOrigins: []string{"*"},
	AllowedMethods: []string{
		http.MethodHead,
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
	},
	AllowedHeaders:     []string{"*"},
	AllowCredentials:   false,
	OptionsPassthrough: true,
}

type Server struct {
	store core.Store
}

func New(
	store core.Store,
) Server {
	return Server{
		store: store,
	}
}

func (s Server) Handler() http.Handler {
	r := chi.NewRouter()

	sessionsDB := s.store.SessionStore()
	usersDB := s.store.UserStore()
	//accessDB := s.store.AccessStore()
	filesDB := s.store.FileStore()

	c := cors.New(corsOpts)
	r.Use(c.Handler)

	permission := acl.New(sessionsDB, usersDB)
	r.Use(permission.Middleware)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", auth.HandleLogin(usersDB, sessionsDB))
		r.Post("/signup", auth.HandleRegister(usersDB))
	})

	r.Route("/files", func(r chi.Router) {
		handler, err := tusd.NewUnroutedHandler(tusd.Config{
			BasePath:      "/api/files",
			StoreComposer: filesDB.GetTusdStoreComposer(),
			//PreUploadCreateCallback:   files.PreUploadCreate(filesDB),
			//PreFinishResponseCallback: files.PreUploadFinish(filesDB),
			NotifyCreatedUploads:  true,
			NotifyCompleteUploads: true,
		})

		if err != nil {
			panic(err)
		}

		go func() {
			for {
				select {
				case hook := <-handler.CreatedUploads:
					_ = filesDB.AddInCache(hook.Upload.ID)
				case hook := <-handler.CompleteUploads:
					_, _ = filesDB.New(hook.Upload.ID)
				}
			}
		}()

		r.Use(handler.Middleware)

		r.Post("/", handler.PostFile)
		r.Route("/{id:[-+a-z0-9]+}", func(r chi.Router) {
			r.Get("/", handler.GetFile)
			r.Head("/", handler.HeadFile)
			r.Patch("/", handler.PatchFile)
			r.Delete("/", handler.DelFile)
		})
	})

	return r
}
