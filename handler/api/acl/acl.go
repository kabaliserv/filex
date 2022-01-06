package acl

import (
	"context"
	"github.com/kabaliserv/filex/core"
	"net/http"
)

type ACL struct {
	sessions core.SessionStore
	users    core.UserStore
}

func New(
	sessions core.SessionStore,
	users core.UserStore,
) ACL {
	return ACL{
		sessions: sessions,
		users:    users,
	}
}

func (a ACL) Middleware(next http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session := a.sessions.Get(r)

		userId, ok := session.Values["userId"].(string)

		if ok && userId != "" {
			user, err := a.users.GetUserById(userId)
			if err == nil {
				ctx = context.WithValue(ctx, "user", user)
			}
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(f)
}
