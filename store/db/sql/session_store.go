package sql

import (
	"github.com/kabaliserv/filex/core"
	"github.com/wader/gormstore/v2"
	"net/http"
)

const (
	sessionName = "session"
)

type SessionStore struct {
	*gormstore.Store
}

func (st *SessionStore) Get(r *http.Request) *core.Session {
	so, _ := st.Store.Get(r, sessionName)
	s := core.Session{Session: so}
	s.Save = func(r *http.Request, w http.ResponseWriter) error {
		return st.Store.Save(r, w, so)
	}
	return &s
}
