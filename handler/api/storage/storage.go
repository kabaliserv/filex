package storage

import (
	"encoding/json"
	"github.com/kabaliserv/filex/core"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func FindOneByUserId(storages core.StorageStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value("user").(*core.User)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error("User are not inject in context")
			return
		}

		storage, err := storages.GetByUserId(user.ID)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error(err)
			return
		} else if storage == nil {
			http.NotFound(w, r)
			return
		}

		out, err := json.Marshal(storage)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error(err)
			return
		}

		w.WriteHeader(http.StatusOK)

		_, err = w.Write(out)
		if err != nil {
			log.Error(err)
			return
		}
	}
}
