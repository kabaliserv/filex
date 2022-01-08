package users

import (
	"encoding/json"
	"github.com/kabaliserv/filex/core"
	"github.com/prometheus/common/log"
	"net/http"
)

type MeResponse map[string]interface{}

func HandlerGetMe() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var res = make(MeResponse)

		res["auth"] = false

		ctx := r.Context()

		user, ok := ctx.Value("user").(*core.User)
		if ok {
			res["auth"] = true
			res["userId"] = user.ID
		}

		out, err := json.Marshal(res)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(out); err != nil {
			log.Error(err)
		}
	}
}
