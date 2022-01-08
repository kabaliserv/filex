package acl

import (
	"context"
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
		ctx := r.Context()
		session := a.sessions.Get(r)

		userId, ok := session.Values["userId"].(string)

		if ok && userId != "" {
			user, err := a.users.Get(userId)
			if err == nil {
				ctx = context.WithValue(ctx, "user", user)
			}
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(f)
}

func (a ACL) UserRequired(next http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		_, ok := ctx.Value("user").(*core.User)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(f)
}

func (a ACL) AdminUserRequired(next http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		user, ok := ctx.Value("user").(*core.User)
		if !ok || !user.Admin {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(f)
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
				rw.WriteHeader(http.StatusUnauthorized)
				return
			}

			clientId := token.Issuer()
			req.Header.Add("clientId", clientId)

			// check if the upload was initiated by the customer
			if req.Method == "PATCH" {
				fileCache, err := a.files.GetInCacheByClientId(clientId)
				if err != nil {
					log.Error(err)
					rw.WriteHeader(http.StatusInternalServerError)
					return
				}

				if fileCache == nil || fileCache.ClientID != clientId {
					rw.WriteHeader(http.StatusUnauthorized)
					return
				}
			}

		}

		next.ServeHTTP(rw, req.WithContext(ctx))
	}
	return http.HandlerFunc(f)
}
