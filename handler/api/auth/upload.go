package auth

import (
	"github.com/kabaliserv/filex/core"
	"github.com/kabaliserv/filex/service/token"
	"github.com/lestrrat-go/jwx/jwt"
	gonanoid "github.com/matoous/go-nanoid"
	"net/http"
	"time"
)

func HandlerGetMe(manager token.Manager, files core.FileStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		token := manager.NewToken()

		token.Set(jwt.SubjectKey, `upload file`)
		token.Set(jwt.ExpirationKey, time.Now().Add(time.Second*5))

		if user, ok := ctx.Value("user").(*core.User); ok && user != nil {
			token.Set("user_id", user.ID)
		}

		// generate client id
		clientId := gonanoid.MustID(10)
		token.Set(jwt.IssuerKey, clientId)

		payload, err := manager.Sign(token)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(payload)
	}
}
