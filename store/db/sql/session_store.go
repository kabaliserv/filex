package sql

import (
	"github.com/kabaliserv/filex/core"
	"github.com/wader/gormstore/v2"
	"gorm.io/gorm"
	"net/http"
)

const (
	sessionName = "session"
)

type sessionStore struct {
	*gormstore.Store
}

func newSessionStore(db *gorm.DB, secret string) *sessionStore {
	return &sessionStore{
		Store: gormstore.New(db, []byte(secret)),
	}
}

func (st *sessionStore) Get(r *http.Request) *core.Session {
	so, _ := st.Store.Get(r, sessionName)
	s := core.Session{Session: so}
	s.Save = func(r *http.Request, w http.ResponseWriter) error {
		return st.Store.Save(r, w, so)
	}
	return &s
}
