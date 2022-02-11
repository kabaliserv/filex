package files

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/kabaliserv/filex/core"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func FileCtx(files core.FileStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := r.Context().Value(core.User{}).(core.User)
			if !ok {
				return
			}

			fileId := chi.URLParam(r, "fileId")
			if fileId == "" {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			_, err := uuid.Parse(fileId)
			if err != nil {
				http.NotFound(w, r)
				return
			}

			file, err := files.FindByUUID(fileId)
			if err != nil {
				if errors.Is(err, core.ErrNotFound) {
					http.NotFound(w, r)
					return
				}
				w.WriteHeader(http.StatusInternalServerError)
				log.Errorf("fileCtx(find in database): %s", err)
				return
			}

			if file.StorageID != user.Storage.ID {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), core.File{}, file)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetAll(files core.FileStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(core.User{}).(core.User)
		if !ok {
			return
		}

		var where core.File

		if tusFileId := r.URL.Query().Get("uploadId"); tusFileId != "" {
			where = core.File{FileId: tusFileId, StorageID: user.Storage.ID}

		} else {
			where = core.File{StorageID: user.Storage.ID}
		}

		f, err := files.Find(where)
		if err != nil && !errors.Is(err, core.ErrNotFound) {
			w.WriteHeader(http.StatusInternalServerError)
			log.Errorf("getAllFile(find in database): %s", err)
			return
		}

		out, err := json.Marshal(f)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Errorf("getAllFile(marshal to json): %s", err)
			return
		}

		w.Write(out)
	}
}

func GetOne() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		file, ok := r.Context().Value(core.File{}).(core.File)
		if !ok {
			return
		}

		out, err := json.Marshal(file)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Errorf("getOneFile(marshal to json): %s", err)
			return
		}

		w.Write(out)
	}
}

func Delete(files core.FileStore, storages core.StorageStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(core.User{}).(core.User)
		if !ok {
			return
		}
		file, ok := r.Context().Value(core.File{}).(core.File)
		if !ok {
			return
		}

		user.Storage.Size -= file.Size

		if err := storages.Save(&user.Storage); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Errorf("deleteFile(update Storage): %s", err)
			return
		}

		if err := files.Delete(&file); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Errorf("deleteFile(detele file): %s", err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
