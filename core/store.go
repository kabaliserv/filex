package core

import (
	"github.com/gorilla/sessions"
	"net/http"
)

type (
	Store interface {
		UserStore() UserStore
		AccessStore() AccessStore
		SessionStore() SessionStore
		Close() error
	}

	UserStore interface {
		GetUserById(id string) (*User, error)
		GetUserByName(name string) (*User, error)
		GetUserByEmail(email string) (*User, error)
		InsertUser(data User) (*User, error)
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
