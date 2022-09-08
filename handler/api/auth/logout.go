package auth

import (
	"github.com/kabaliserv/filex/core"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func HandleLogout(sessions core.SessionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, ok := r.Context().Value(core.User{}).(core.User); !ok {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		session := sessions.Get(r)

		delete(session.Values, "userId")

		if err := session.Save(r, w); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Errorf("logout(save-session): %s", err)
			return
		}

		w.WriteHeader(http.StatusNoContent)

	}
}
