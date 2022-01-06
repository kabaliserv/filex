package auth

import (
	"github.com/go-chi/chi/v5"
	"github.com/kabaliserv/filex/core"
	"net/http"
)

func Handler(store core.Store) http.Handler {
	r := chi.NewRouter()

	userStore := store.UserStore()
	sessionStore := store.SessionStore()

	r.Post("/login", HandleLogin(userStore, sessionStore))
	r.Post("/signup", HandleRegister(userStore))

	return r
}
