package users

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/kabaliserv/filex/core"
	"github.com/prometheus/common/log"
	"net/http"
)

type MeResponse map[string]interface{}

func HandlerGetMe() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var res = make(MeResponse)

		res["auth"] = false

		ctx := r.Context()

		user, ok := ctx.Value(core.User{}).(*core.User)
		if ok {
			res["auth"] = true
			res["userId"] = user.ID
		}

		out, err := json.Marshal(res)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(out); err != nil {
			log.Error(err)
		}
	}
}

func HandlerGetAll(users core.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		allUsers, err := users.Find(core.User{})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error(err)
			return
		}

		if len(allUsers) == 0 {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("content-type", "application/json")
			w.Write([]byte("[]"))
			return
		}

		out, err := json.Marshal(allUsers)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error(err)
			return
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err = w.Write(out)
		if err != nil {
			log.Error(err)
		}
	}
}

func HandlerGetOne(users core.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIdString := chi.URLParam(r, "userId")
		if userIdString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userId, err := uuid.Parse(userIdString)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		user, err := users.FindByUUID(userId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error(err)
			return
		} else if user == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		out, err := json.Marshal(user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error(err)
			return
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err = w.Write(out)
		if err != nil {
			log.Error(err)
		}
	}
}

func HandlerPost(users core.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("content-type") != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var user core.User

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Error(err)
			return
		}

		// Code de verification de l'utilisateur à faire ...

		user.ID = 0
		user.Storage.ID = 0
		user.Storage.UserID = 0

		if user.Storage.EnableQuota && user.Storage.Quota == 0 {
			user.Storage.Quota = 1073741824 // 1GB
		}

		if err := users.Create(&user); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
