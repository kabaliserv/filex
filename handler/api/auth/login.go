package auth

import (
	"encoding/json"
	"github.com/kabaliserv/filex/core"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type Credential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AccessToken struct {
	AccessToken string `json:"access_token"`
}

func HandleLogin(userStore core.UserStore, sessionStore core.SessionStore) http.HandlerFunc {
	f := func(w http.ResponseWriter, r *http.Request) {
		session := sessionStore.Get(r)

		if auth, ok := session.Values["auth"].(bool); ok && auth {
			w.WriteHeader(204)
			return
		}

		if r.Header.Get("Content-Type") != "application/json" {
			// ERROR ...
			return
		}

		var c Credential

		if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
			// ERROR ...
			return
		}

		var u *core.User
		var err error

		if ValidUserName(c.Username) {
			u, err = userStore.GetUserByName(c.Username)
		}

		if ValidEmail(c.Username) {
			u, err = userStore.GetUserByEmail(c.Username)
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			// ERROR ...
			return
		}

		if u == nil {
			w.WriteHeader(http.StatusUnauthorized)
			// ERROR ...
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(c.Password)); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		session.Values["auth"] = true

		if err := session.Save(r, w); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			// ERROR ...
			return
		}

		w.WriteHeader(http.StatusNoContent)
		return
	}

	return f
}
