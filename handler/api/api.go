package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/kabaliserv/filex/core"
	acl "github.com/kabaliserv/filex/handler/api/acl"
	"github.com/kabaliserv/filex/handler/api/auth"
	"github.com/kabaliserv/filex/handler/api/users"
	"github.com/kabaliserv/filex/service/token"
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
	uploadOpt core.UploadOption
	manager   token.Manager
	users     core.UserStore
	sessions  core.SessionStore
	files     core.FileStore
}

func New(
	uploadOpt core.UploadOption,
	manager token.Manager,
	users core.UserStore,
	sessions core.SessionStore,
	files core.FileStore,
) Server {
	return Server{
		uploadOpt: uploadOpt,
		manager:   manager,
		users:     users,
		sessions:  sessions,
		files:     files,
	}
}

func (s Server) Handler() http.Handler {
	r := chi.NewRouter()

	c := cors.New(corsOpts)
	r.Use(c.Handler)

	permission := acl.New(s.sessions, s.users, s.files, s.manager)
	r.Use(permission.Middleware)
	r.Use(middleware.NoCache)

	r.Route("/auth", func(rr chi.Router) {

		rr.Post("/login", auth.HandleLogin(s.users, s.sessions))
		rr.Post("/signup", auth.HandleRegister(s.users))

	})

	r.Route("/users", func(r chi.Router) {

		r.With(permission.AdminUserRequired).
			Get("/", users.HandlerGetAll(s.users))
		r.With(permission.AdminUserRequired).
			Post("/", users.HandlerPost(s.users))

		r.Route("/{userId}", func(r chi.Router) {
			r.Use(permission.UserRequired)

			r.Get("/", users.HandlerGetOne(s.users))
			//r.Patch("/")

		})
	})

	//r.Route("/me", func(rr chi.Router) {
	//
	//	rr.Use(permission.UserRequired)
	//	rr.Get("/", users.HandlerGetMe())
	//	rr.Get("/storage", storage.FindOneByUserId(storageDB))
	//
	//})
	//
	//r.Route("/users", func(rr chi.Router) {
	//
	//	rr.Use(permission.UserRequired)
	//
	//	rr.Route("/{id}", func(rrr chi.Router) {
	//
	//		rrr.Get("/storage", storage.FindOneByUserId(storageDB))
	//
	//	})
	//
	//})

	//r.Route("/files", func(rr chi.Router) {
	//
	//	handler, err := tusd.NewUnroutedHandler(tusd.Config{
	//		BasePath:                "/api/files",
	//		StoreComposer:           s.files.TusdStoreComposer(),
	//		PreUploadCreateCallback: files.PreUploadCreate(s.manager, s.uploadOpt),
	//		NotifyCreatedUploads:    true,
	//		NotifyCompleteUploads:   true,
	//	})
	//
	//	if err != nil {
	//		panic(err)
	//	}
	//
	//	go func() {
	//		for {
	//			select {
	//			case hook := <-handler.CreatedUploads:
	//				contextUploadId := hook.HTTPRequest.Header.Get("filex-context-upload-id")
	//				cacheFile := core.FileCache{
	//					ID:              hook.Upload.ID,
	//					ContextUploadID: contextUploadId,
	//				}
	//				_ = s.files.CreateInCache(&cacheFile)
	//			case hook := <-handler.CompleteUploads:
	//				_ = s.files.CreateFromCache(hook.Upload.ID)
	//			}
	//		}
	//	}()

	//  rr.Post("/upload", auth.HandlerGetMe(s.manager, filesDB))
	//
	//	rr.Use(handler.Middleware)
	//
	//	rr.With(permission.SecureUpload).
	//		Post("/", handler.PostFile)
	//
	//	rr.Route("/{id:[-+a-z0-9]+}", func(rrr chi.Router) {
	//
	//		rrr.With(permission.SecureUpload).
	//			Patch("/", handler.PatchFile)
	//
	//		rrr.Get("/", handler.GetFile)
	//		rrr.Head("/", handler.HeadFile)
	//		rrr.Delete("/", handler.DelFile)
	//
	//	})
	//
	//})

	return r
}
