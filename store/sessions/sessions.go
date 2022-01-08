package sessions

import (
	"github.com/kabaliserv/filex/core"
	"github.com/wader/gormstore/v2"
	"gorm.io/gorm"
	"net/http"
)

const sessionName = "session"

type sessionStore struct {
	*gormstore.Store
}

func NewSessionStore(db *gorm.DB, options core.StoreOption) core.SessionStore {
	return &sessionStore{
		Store: gormstore.New(db, []byte(options.SessionSecret)),
	}
}

func (st *sessionStore) Get(r *http.Request) *core.Session {
	so, _ := st.Store.Get(r, sessionName)
	s := core.Session{Session: so}
	s.Save = func(rr *http.Request, w http.ResponseWriter) error {
		return st.Store.Save(rr, w, so)
	}
	return &s
}
