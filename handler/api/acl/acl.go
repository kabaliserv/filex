package acl

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/kabaliserv/filex/core"
	"github.com/kabaliserv/filex/service/token"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type ACL struct {
	sessions core.SessionStore
	users    core.UserStore
	files    core.FileStore
	tokens   token.Manager
}

func New(
	sessions core.SessionStore,
	users core.UserStore,
	files core.FileStore,
	tokens token.Manager,
) ACL {
	return ACL{
		sessions: sessions,
		users:    users,
		files:    files,
		tokens:   tokens,
	}
}

func (a ACL) Middleware(next http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		session := a.sessions.Get(r)

		userId, ok := session.Values["userId"].(uint)

		if !ok || userId == 0 {
			next.ServeHTTP(w, r)
			return
		}

		user, _ := a.users.FindByID(userId)

		ctx := context.WithValue(r.Context(), core.User{}, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(f)
}

func (a ACL) UserRequired(next http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		_, ok := ctx.Value(core.User{}).(core.User)
		if !ok {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(f)
}

func (a ACL) AdminUserRequired(next http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {

		user, ok := r.Context().Value(core.User{}).(core.User)
		if !ok || !user.Admin {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(f)
}

func (a ACL) RequireSelfUserOrAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		user, ok := r.Context().Value(core.User{}).(core.User)
		if !ok {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		userId := chi.URLParam(r, "userId")
		if userId == "" || (userId != user.UUID && !user.Admin) {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (a ACL) SecureUpload(next http.Handler) http.Handler {
	f := func(rw http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		if req.Method == "PATCH" || req.Method == "POST" {
			token, err := a.tokens.GetTokenFromHeaders(req.Header)
			if err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				log.Error(err)
				return
			} else if token == nil {
				rw.WriteHeader(http.StatusForbidden)
				return
			}

			FileID := token.Issuer()
			req.Header.Add("filex-File-ID", FileID)

			// check if the upload was initiated by the customer
			if req.Method == "PATCH" {
				file, err := a.files.FindByUUID(FileID)
				if err != nil {
					if errors.Is(err, core.ErrNotFound) {
						rw.WriteHeader(http.StatusUnauthorized)
						return
					}

					log.Error(err)
					rw.WriteHeader(http.StatusInternalServerError)
					return
				}

				if file.IsFinal {
					rw.WriteHeader(http.StatusForbidden)
					return
				}
			}

		}

		next.ServeHTTP(rw, req.WithContext(ctx))
	}
	return http.HandlerFunc(f)
}
