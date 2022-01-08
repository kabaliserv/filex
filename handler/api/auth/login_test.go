package auth

import (
	"github.com/kabaliserv/filex/core"
	"github.com/kabaliserv/filex/store/db/sql"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLoginWithUsername(t *testing.T) {
	options := core.StoreOption{
		FileStoreLocalPath: "/tmp/files",
		DatabaseDriver:     "sqlite3",
		DatabaseEndpoint:   "file::memory:?cache=shared",
	}

	db := sql.New(options)

	defer func(db core.Store) {
		err := db.CloseConnection()
		if err != nil {
			t.Error(err)
		}
	}(db)

	userDB := db.UserStore()
	sessionDB := db.SessionStore()

	p, err := bcrypt.GenerateFromPassword([]byte("C0mpleX_P@ssw0rd"), PasswordCost)

	err = userDB.Add(&core.User{
		Username:     "test",
		Email:        "test@gmail.com",
		PasswordHash: string(p),
	})

	if err != nil {
		t.Error(err)
	}

	f := HandleLogin(userDB, sessionDB)

	{ // test login with username
		body := `{"username":"test","password":"C0mpleX_P@ssw0rd"}`

		req := httptest.NewRequest("POST", "/fake", strings.NewReader(body))
		w := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")

		f(w, req)
		res := w.Result()
		defer res.Body.Close()

		if httpCode := res.StatusCode; httpCode != http.StatusNoContent {
			t.Errorf("expected http code 204 got %v", httpCode)
		}
	}

}

func TestLoginWithEmail(t *testing.T) {
	options := core.StoreOption{
		FileStoreLocalPath: "/tmp/files",
		DatabaseDriver:     "sqlite3",
		DatabaseEndpoint:   "file::memory:?cache=shared",
	}

	db := sql.New(options)

	defer func(db core.Store) {
		err := db.CloseConnection()
		if err != nil {
			t.Error(err)
		}
	}(db)

	userDB := db.UserStore()
	sessionDB := db.SessionStore()

	p, err := bcrypt.GenerateFromPassword([]byte("C0mpleX_P@ssw0rd"), PasswordCost)

	err = userDB.Add(&core.User{
		Username:     "test",
		Email:        "test@gmail.com",
		PasswordHash: string(p),
	})

	if err != nil {
		t.Error(err)
	}

	f := HandleLogin(userDB, sessionDB)

	{ // test login with email
		body := `{"username":"test@gmail.com","password":"C0mpleX_P@ssw0rd"}`

		req := httptest.NewRequest("POST", "/fake", strings.NewReader(body))
		w := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")

		f(w, req)
		res := w.Result()
		defer res.Body.Close()

		if httpCode := res.StatusCode; httpCode != http.StatusNoContent {
			t.Errorf("expected http code 204 got %v", httpCode)
		}
	}

}
