package session

import (
	"net/http"

	"github.com/kabaliserv/filex/core"
	"github.com/kabaliserv/filex/store/db/sql"

	"github.com/wader/gormstore/v2"
)

const (
	sessionName = "session"
)

type Store struct {
	*gormstore.Store
}

func New(store core.Store) *Store {
	sqlDB, _ := store.(*sql.Store)
	sessionStore := gormstore.New(sqlDB.GetDB(), []byte("azertFRDSFV563yuiopqsdfVD36514ghjklmwxRcvbnd6f8GRSGs13rg-51h-8hhFs"))

	return &Store{sessionStore}
}

func (st *Store) Get(r *http.Request) *core.Session {
	so, _ := st.Store.Get(r, sessionName)
	s := core.Session{Session: so}
	s.Save = func(r *http.Request, w http.ResponseWriter) error {
		return st.Store.Save(r, w, so)
	}
	return &s
}
