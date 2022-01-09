package auth

import (
	"encoding/json"
	"github.com/kabaliserv/filex/core"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type Credential struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type AccessToken struct {
	AccessToken string `json:"access_token"`
}

func HandleLogin(userStore core.UserStore, sessionStore core.SessionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session := sessionStore.Get(r)
		isAuth, authOk := session.Values["auth"].(bool)
		userId, userOk := session.Values["userId"].(string)

		if authOk && userOk && isAuth && userId != "" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		var c Credential

		if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if c.Login == "" || c.Password == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		var user *core.User
		var err error

		if ValidEmail(c.Login) {
			user, err = userStore.FindByEmail(c.Login)
		} else if ValidUserName(c.Login) {
			user, err = userStore.FindByLogin(c.Login)
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if user == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(c.Password)); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		session.Values["auth"] = true
		session.Values["userId"] = user.ID

		if err := session.Save(r, w); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)

		return
	}
}
