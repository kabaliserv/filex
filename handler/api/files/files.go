package files

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/kabaliserv/filex/storage"
	tusd "github.com/tus/tusd/pkg/handler"
	"net/http"
)

func Handle(storage storage.Storage) http.Handler {
	r := chi.NewRouter()

	handler, err := tusd.NewUnroutedHandler(tusd.Config{
		BasePath:      "/api/files",
		StoreComposer: storage.GetStoreComposer(),
		PreUploadCreateCallback: func(hook tusd.HookEvent) error {

			return errors.New("totoot")
		},
	})

	if err != nil {
		panic(err)
	}

	r.Use(handler.Middleware)

	r.Post("/", handler.PostFile)
	r.Patch("/{id:[-+a-z0-9]+}", handler.PatchFile)

	rr := r.With(acl)
	rr.Head("/{id:[-+a-z0-9]+}", handler.HeadFile)
	rr.Get("/{id:[-+a-z0-9]+}", handler.GetFile)
	rr.Delete("/{id:[-+a-z0-9]+}", handler.DelFile)

	return r
}
