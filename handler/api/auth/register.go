package auth

import (
	"encoding/json"
	"github.com/kabaliserv/filex/core"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

const (
	PasswordCost = 10
)

type CreateUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func HandleRegister(users core.UserStore, storages core.StorageStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusBadGateway)
			return
		}

		var c CreateUser

		if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}

		if !ValidUserName(c.Username) || !ValidEmail(c.Email) || c.Password == "" {
			w.WriteHeader(http.StatusBadGateway)
			return
		}

		fromPassword, err := bcrypt.GenerateFromPassword([]byte(c.Password), PasswordCost)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		user := core.User{
			Username:     c.Username,
			PasswordHash: string(fromPassword),
			Email:        c.Email,
		}

		if err := users.Add(&user); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error(err)
			return
		}

		storage := core.Storage{UserID: user.ID}

		if err := storages.Add(&storage); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error(err)
			return
		}

		w.WriteHeader(http.StatusNoContent)

		return
	}
}
