package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/kabaliserv/filex/core"
	acl "github.com/kabaliserv/filex/handler/api/acl"
	"github.com/kabaliserv/filex/handler/api/auth"
	"github.com/kabaliserv/filex/handler/api/files"
	"github.com/kabaliserv/filex/handler/api/storage"
	"github.com/kabaliserv/filex/handler/api/users"
	"github.com/kabaliserv/filex/service/token"
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
	store     core.Store
	uploadOpt core.UploadOption
	manager   token.Manager
}

func New(
	store core.Store,
	uploadOpt core.UploadOption,
	manager token.Manager,
) Server {
	return Server{
		store:     store,
		uploadOpt: uploadOpt,
		manager:   manager,
	}
}

func (s Server) Handler() http.Handler {
	r := chi.NewRouter()

	sessionsDB := s.store.SessionStore()
	usersDB := s.store.UserStore()
	//accessDB := s.store.AccessStore()
	filesDB := s.store.FileStore()
	storageDB := s.store.StorageStore()

	c := cors.New(corsOpts)
	r.Use(c.Handler)

	permission := acl.New(sessionsDB, usersDB, filesDB, s.manager)
	r.Use(permission.Middleware)
	r.Use(middleware.NoCache)

	r.Route("/auth", func(rr chi.Router) {

		rr.Post("/login", auth.HandleLogin(usersDB, sessionsDB))
		rr.Post("/signup", auth.HandleRegister(usersDB, storageDB))
		rr.Post("/upload", auth.HandlerGetMe(s.manager, filesDB))

	})

	r.Route("/me", func(rr chi.Router) {

		rr.Use(permission.UserRequired)
		rr.Get("/", users.HandlerGetMe())
		rr.Get("/storage", storage.FindOneByUserId(storageDB))

	})

	r.Route("/users", func(rr chi.Router) {

		rr.Use(permission.UserRequired)

		rr.Route("/{id}", func(rrr chi.Router) {

			rrr.Get("/storage", storage.FindOneByUserId(storageDB))

		})

	})

	r.Route("/files", func(rr chi.Router) {

		handler, err := tusd.NewUnroutedHandler(tusd.Config{
			BasePath:                "/api/files",
			StoreComposer:           filesDB.GetTusdStoreComposer(),
			PreUploadCreateCallback: files.PreUploadCreate(s.manager, s.uploadOpt),
			NotifyCreatedUploads:    true,
			NotifyCompleteUploads:   true,
		})

		if err != nil {
			panic(err)
		}

		go func() {
			for {
				select {
				case hook := <-handler.CreatedUploads:
					clientId := hook.HTTPRequest.Header.Get("clientId")
					_ = filesDB.AddInCache(hook.Upload.ID, clientId)
				case hook := <-handler.CompleteUploads:
					_, _ = filesDB.New(hook.Upload.ID)
					_ = filesDB.DelInCache(hook.Upload.ID)
				}
			}
		}()

		rr.Use(handler.Middleware)

		rr.With(permission.SecureUpload).
			Post("/", handler.PostFile)

		rr.Route("/{id:[-+a-z0-9]+}", func(rrr chi.Router) {

			rrr.With(permission.SecureUpload).
				Patch("/", handler.PatchFile)

			rrr.Get("/", handler.GetFile)
			rrr.Head("/", handler.HeadFile)
			rrr.Delete("/", handler.DelFile)

		})

	})

	return r
}
