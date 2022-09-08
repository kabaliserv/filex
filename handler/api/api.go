package api

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/kabaliserv/filex/core"
	acl "github.com/kabaliserv/filex/handler/api/acl"
	"github.com/kabaliserv/filex/handler/api/auth"
	"github.com/kabaliserv/filex/handler/api/download"
	"github.com/kabaliserv/filex/handler/api/files"
	"github.com/kabaliserv/filex/handler/api/upload"
	"github.com/kabaliserv/filex/handler/api/users"
	"github.com/kabaliserv/filex/service/token"
	log "github.com/sirupsen/logrus"
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
	files          core.FileStore
	sessions       core.SessionStore
	storages       core.StorageStore
	tokens         token.Manager
	uploadOpt      core.UploadOption
	uploads        core.UploadStore
	users          core.UserStore
	uploadAccesses core.AccessUploadStore
	fileStorage    core.FileStoreComposer
}

func New(
	files core.FileStore,
	sessions core.SessionStore,
	storages core.StorageStore,
	tokens token.Manager,
	uploadOpt core.UploadOption,
	uploads core.UploadStore,
	users core.UserStore,
	uploadAccesses core.AccessUploadStore,
	fileStorage core.FileStoreComposer,
) Server {
	return Server{
		files:          files,
		sessions:       sessions,
		storages:       storages,
		tokens:         tokens,
		uploadOpt:      uploadOpt,
		uploads:        uploads,
		users:          users,
		uploadAccesses: uploadAccesses,
		fileStorage:    fileStorage,
	}
}

func (s Server) Handler() http.Handler {
	r := chi.NewRouter()

	c := cors.New(corsOpts)
	r.Use(c.Handler)

	permission := acl.New(s.sessions, s.users, s.files, s.tokens)
	uploads := upload.NewUploadRouter(s.uploadOpt, s.files, s.users, s.tokens, s.uploads, s.uploadAccesses, s.storages, s.fileStorage)

	r.Use(middleware.NoCache)
	r.Use(middleware.Recoverer)
	r.Use(permission.Middleware)

	r.Route("/access", func(r chi.Router) {
		//r.Get("/upload")
		//r.Get("/download")
		r.Route("/requests", func(r chi.Router) {
			//r.Post("/upload",)
			//r.Post("/download")
		})
	})

	r.Route("/auth", func(r chi.Router) {

		r.Post("/login", auth.HandleLogin(s.users, s.sessions))
		r.Post("/logout", auth.HandleLogout(s.sessions))
		r.Post("/signup", auth.HandleRegister(s.users))
		r.Get("/check", auth.HandlerCheckAuth(s.sessions))

	})

	r.Route("/files", func(r chi.Router) {
		r.Use(permission.UserRequired)

		r.Get("/", files.GetAll(s.files))

		r.Route("/{fileId}", func(r chi.Router) {
			r.Use(files.FileCtx(s.files))

			r.Get("/", files.GetOne())
			r.Delete("/", files.Delete(s.files, s.storages))

			r.Get("/download", download.GetByFileCtx(s.files))

		})
	})

	r.Route("/users", func(r chi.Router) {
		r.With(permission.AdminUserRequired).
			Get("/", users.HandleGetAll(s.users))
		r.With(permission.AdminUserRequired).
			Post("/", users.HandlePost(s.users))

		r.Route("/{userId}", func(r chi.Router) {
			r.Use(permission.RequireSelfUserOrAdmin)

			r.Get("/", users.HandleGet(s.users))
			r.Post("/change-password", users.HandleChangePassword(s.users))

			//r.Patch("/")

		})
	})

	r.Mount("/upload", uploads.Handler())

	r.Get("/me", users.HandleGetMe())

	r.Get("/serverOptions", s.serverOptions)

	return r
}

func (s Server) serverOptions(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"signup": true,
		"guest": map[string]interface{}{
			"upload":  true,
			"maxSize": 1073741824,
		},
	}

	out, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}
