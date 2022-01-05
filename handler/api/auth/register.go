package auth

import (
	"encoding/json"
	"fmt"
	"github.com/kabaliserv/filex/core"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

const (
	PasswordCost = 15
)

type CreateUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func HandleRegister(store core.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(400)
			// ERROR ...
			return
		}

		var c CreateUser

		if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
			w.WriteHeader(500)
			// ERROR ...
			return
		}

		if !ValidUserName(c.Username) || !ValidEmail(c.Email) || c.Password == "" {
			w.WriteHeader(403)
			// Error ...
			return
		}

		p, err := bcrypt.GenerateFromPassword([]byte(c.Password), PasswordCost)

		if err != nil {
			w.WriteHeader(500)
			// ERROR ...
			return
		}

		u := core.User{
			Username:     c.Username,
			PasswordHash: string(p),
			Email:        c.Email,
		}

		_, err = store.InsertUser(u)

		if err != nil {
			w.WriteHeader(500)
			fmt.Println("insert ERROR")
			// ERROR ...
			return
		}

		w.WriteHeader(http.StatusNoContent)

		return
	}
}
