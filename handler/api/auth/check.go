package auth

import (
	"encoding/json"
	"github.com/kabaliserv/filex/core"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func HandlerCheckAuth(sessions core.SessionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := make(map[string]interface{})

		res["auth"] = false

		session := sessions.Get(r)

		user, ok := r.Context().Value(core.User{}).(core.User)

		if ok && user.Active {
			res["auth"] = true
		} else {
			delete(session.Values, "userId")
		}

		if err := session.Save(r, w); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Errorf("auth check(save-session): %s", err)
			return
		}

		out, err := json.Marshal(res)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Errorf("auth check(parse-to-json): %s", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
	}
}
