package core

import (
	"github.com/gorilla/sessions"
	"net/http"
)

type (
	Store interface {
		AccessStore() AccessStore
		SessionStore() SessionStore
		FileStore() FileStore
		StorageStore() StorageStore
		CloseConnection() error
	}

	AccessStore interface {
		GetAccess(token string) (interface{}, error)
		GetUserAccess(token string) (*UserAccess, error)
		NewUserAccess() (*UserAccess, error)
	}

	SessionStore interface {
		Get(r *http.Request) *Session
	}

	Session struct {
		*sessions.Session
		Save func(r *http.Request, w http.ResponseWriter) error
	}
)
