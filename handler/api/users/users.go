package users

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/kabaliserv/filex/core"
	"github.com/kabaliserv/filex/handler/api/auth"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type MeResponse map[string]interface{}

func HandleGetMe() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(core.User{}).(core.User)

		out, err := json.Marshal(user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Errorf("users__get-me(marshalJSON): %s", err)
			return
		}

		if _, err := w.Write(out); err != nil {
			log.Error(err)
		}
	}
}

func HandleGetAll(users core.UserStore) http.HandlerFunc {
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

func HandleGet(users core.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userId")
		if userID == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		_, err := uuid.Parse(userID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		user, err := users.FindByUUID(userID)
		if err != nil {
			if errors.Is(err, core.ErrNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			log.Error(err)
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

type responseNewUser struct {
	Login       string `json:"login"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Active      bool   `json:"active"`
	Admin       bool   `json:"admin"`
	Quota       int64  `json:"quota"`
	ActiveQuota bool   `json:"activeQuota"`
}

func HandlePost(users core.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("content-type") != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var user core.User
		var resUser responseNewUser

		if err := json.NewDecoder(r.Body).Decode(&resUser); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Error(err)
			return
		}

		user.Login = resUser.Login
		user.Email = resUser.Email
		user.Active = resUser.Active
		user.Admin = resUser.Admin
		user.Storage.Quota = resUser.Quota
		user.Storage.QuotaEnable = resUser.ActiveQuota

		if err := user.Validate(); err != nil || resUser.Password == "" || len(resUser.Password) < 3 {
			w.WriteHeader(http.StatusBadRequest)
			if err != nil {
				w.Write([]byte(err.Error()))
				return
			}
			w.Write([]byte("Invalid email length"))
			return
		}

		passwordHash, err := bcrypt.GenerateFromPassword([]byte(resUser.Password), auth.PasswordCost)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		user.PasswordHash = string(passwordHash)

		// Code de verification de l'utilisateur à faire ...

		if user.Storage.QuotaEnable && user.Storage.Quota == 0 {
			user.Storage.Quota = 1073741824 // 1GB
		}

		res, err := http.Get("https://picsum.photos/200")
		if err == nil {
			user.Avatar = res.Request.URL.String()
		}

		if err := users.Create(&user); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Errorf("postNewUser(create user in database): %s", err)
			return
		}

		out, err := json.Marshal(user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Errorf("postNewUser(marshal user to json): %s", err)
		}

		w.Write(out)
	}
}

func HandlePatch(users core.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

type resChangePassword struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func HandleChangePassword(users core.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("content type expected application/json"))
			return
		}

		userId := chi.URLParam(r, "userId")

		var contentReq resChangePassword

		if err := json.NewDecoder(r.Body).Decode(&contentReq); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal error"))
			log.Errorf("changePassword(decode json): %s", err)
			return
		}

		if contentReq.OldPassword == "" || contentReq.NewPassword == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("bad request"))
			return
		}

		user, err := users.FindByUUID(userId)
		if err != nil {
			if errors.Is(err, core.ErrNotFound) {
				http.NotFound(w, r)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal error"))
			log.Errorf("changePassword(find user in database): %s", err)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(contentReq.OldPassword)); err != nil {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("forbidden"))
			return
		}

		pass, err := bcrypt.GenerateFromPassword([]byte(contentReq.NewPassword), auth.PasswordCost)
		if err != nil {
			log.Errorf("changePassword(hash password): %s", err)
			return
		}

		user.PasswordHash = string(pass)

		if err := users.Save(&user); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal error"))
			log.Errorf("changePassword(update user in database): %s", err)
			return
		}
	}
}
