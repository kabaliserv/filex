package metadata

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/kabaliserv/filex/core"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func HandlerGetMetadata(files core.FileStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fileId := chi.URLParam(r, "fileId")
		if fileId == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		file, err := files.FindById(fileId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Errorf("getMetadata(find in database): %s", err)
			return
		}

		if file == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		out, err := json.Marshal(file)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Errorf("getMetadata(marshal to json): %s", err)
			return
		}

		w.Write(out)

	}
}

type resPatchFile struct {
	Name string `json:"name"`
}

func HandlerPatchMetadata(files core.FileStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var data resPatchFile

		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Errorf("patchMetadata(decode json): %s", err)
			return
		}

		log.Println(data)

		if data.Name == "" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		fileId := chi.URLParam(r, "fileId")
		if fileId == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		file, err := files.FindById(fileId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Errorf("patchMetadata(find in database): %s", err)
			return
		}

		if file == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		user := r.Context().Value(core.User{}).(*core.User)

		if file.StorageID != user.Storage.ID {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		file.Name = data.Name

		if err := files.Save(file); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Errorf("patchMetadata(save in database): %s", err)
			return
		}

		out, err := json.Marshal(file)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Errorf("patchMetadata(marshal to json): %s", err)
			return
		}

		w.Write(out)
	}
}
